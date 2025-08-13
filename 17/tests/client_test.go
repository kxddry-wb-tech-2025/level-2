package tests

import (
	"bytes"
	"context"
	"net"
	"strings"
	"telnet/internal/config"
	"telnet/internal/run"
	"testing"
	"time"
)

// startMockServer поднимает TCP echo-сервер для тестов
func startMockServer(t *testing.T) (addr string, stop func()) {
	ln, err := net.Listen("tcp", "127.0.0.1:0") // свободный порт
	if err != nil {
		t.Fatalf("failed to start mock server: %v", err)
	}

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			go func(c net.Conn) {
				defer c.Close()
				buf := make([]byte, 1024)
				for {
					n, err := c.Read(buf)
					if err != nil {
						return
					}
					_, _ = c.Write(buf[:n]) // echo обратно
				}
			}(conn)
		}
	}()

	return ln.Addr().String(), func() { ln.Close() }
}

func TestTelnetClient_Exchange(t *testing.T) {
	addr, stopServer := startMockServer(t)
	defer stopServer()

	input := bytes.NewBufferString("hello\nworld\n")
	output := &bytes.Buffer{}

	opt := &config.Options{
		Timeout: 2 * time.Second,
		Address: addr,
	}
	ctx, cancel := context.WithCancel(context.Background())

	err := run.Run(ctx, cancel, opt, input, output)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got := output.String()
	if !strings.Contains(got, "hello") || !strings.Contains(got, "world") {
		t.Errorf("unexpected output: %q", got)
	}
}

func TestTelnetClient_ServerCloses(t *testing.T) {
	// сервер закроет соединение после первого сообщения
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()

	go func() {
		conn, _ := ln.Accept()
		conn.Close() // сразу закрываем
	}()

	input := bytes.NewBufferString("ping\n")
	output := &bytes.Buffer{}
	opt := &config.Options{
		Timeout: 2 * time.Second,
		Address: ln.Addr().String(),
	}
	ctx, cancel := context.WithCancel(context.Background())
	err = run.Run(ctx, cancel, opt, input, output)
	if err == nil {
		t.Error("expected error due to closed connection, got nil")
	}
}

func TestTelnetClient_Timeout(t *testing.T) {
	// Запускаем "медленный" сервер, который принимает соединение, но не отвечает
	ln, err := net.Listen("tcp", "127.0.0.1:0") // :0 — выбрать свободный порт
	if err != nil {
		t.Fatalf("failed to start fake server: %v", err)
	}
	defer ln.Close()

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			// Держим соединение открытым, но ничего не пишем — имитация зависшего сервера
			time.Sleep(10 * time.Second)
			conn.Close()
		}
	}()

	input := &bytes.Buffer{}
	output := &bytes.Buffer{}

	start := time.Now()
	opt := &config.Options{
		Timeout: 2 * time.Second, // ожидаем, что по этому таймауту упадёт
		Address: ln.Addr().String(),
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = run.Run(ctx, cancel, opt, input, output)
	elapsed := time.Since(start)

	if err == nil {
		t.Error("expected timeout error, got nil")
	}
	if elapsed < 2*time.Second || elapsed > 3*time.Second {
		t.Errorf("expected timeout around 2s, got %v", elapsed)
	}
}

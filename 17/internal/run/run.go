package run

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"sync"
	"telnet/internal/config"
)

func Run(ctx context.Context, cancel context.CancelFunc, opt *config.Options, in io.Reader, out io.Writer) error {
	conn, err := net.DialTimeout("tcp", opt.Address, opt.Timeout)
	if err != nil {
		return err
	}
	defer conn.Close()

	var wg sync.WaitGroup
	wg.Add(2)

	errCh := make(chan error, 1)
	// Чтение из in → запись в сокет
	go func() {
		defer wg.Done()
		if _, err := io.Copy(conn, in); err != nil && !errors.Is(err, io.EOF) {
			fmt.Fprintln(os.Stderr, "write error:", err)
			errCh <- err
		}
		// Закрываем write-часть, чтобы сервер понял, что данных больше не будет
		if tcpConn, ok := conn.(*net.TCPConn); ok {
			_ = tcpConn.CloseWrite()
		}
	}()

	// Чтение из сокета → запись в out
	go func() {
		defer wg.Done()
		if _, err := io.Copy(out, conn); err != nil && !errors.Is(err, io.EOF) {
			fmt.Fprintln(os.Stderr, "read error:", err)
			errCh <- err
		}
		// Если сервер закрыл соединение — завершаем контекст
		cancel()
	}()

	// Ждём либо завершения контекста, либо таймаута
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-ctx.Done():
	case <-done:
	}

	select {
	case err, ok := <-errCh:
		if ok && err != nil {
			return err
		}
	default:
	}

	return nil
}

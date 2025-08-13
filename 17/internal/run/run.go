package run

import (
	"context"
	"fmt"
	"io"
	"net"
	"sync"
	"telnet/internal/config"
	"time"
)

// Run connects to the TCP server and relays data between stdin and stdout.
func Run(ctx context.Context, cancel context.CancelFunc, opt *config.Options, in io.Reader, out io.Writer) error {
	defer cancel()

	conn, err := net.DialTimeout("tcp", opt.Address, opt.Timeout)
	if err != nil {
		return fmt.Errorf("failed to connect to %s: %w", opt.Address, err)
	}
	defer func() { _ = conn.Close() }()

	ctxTimeout, cancelTimeout := context.WithTimeout(ctx, opt.Timeout)
	defer cancelTimeout()

	go func() {
		<-ctxTimeout.Done()
		_ = conn.SetDeadline(time.Now())
	}()

	var (
		wg    sync.WaitGroup
		errCh = make(chan error, 2)
	)

	// server -> stdout
	wg.Add(1)
	go func() {
		defer wg.Done()
		buf := make([]byte, 32*1024)
		for {
			n, err := conn.Read(buf)
			if n > 0 {
				if _, werr := out.Write(buf[:n]); werr != nil {
					errCh <- werr
					return
				}
			}
			if err != nil {
				errCh <- err
				return
			}
		}
	}()

	// stdin -> server
	wg.Add(1)
	go func() {
		defer wg.Done()
		if _, err := io.Copy(conn, in); err != nil && err != io.EOF {
			errCh <- err
		}
		if tcp, ok := conn.(*net.TCPConn); ok {
			_ = tcp.CloseWrite()
		} else {
			_ = conn.Close()
		}
	}()

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		close(errCh)
		var firstErr error
		for e := range errCh {
			if e != nil && firstErr == nil {
				firstErr = e
			}
		}
		return firstErr
	case <-ctxTimeout.Done():
		_ = conn.SetDeadline(time.Now())
		_ = conn.Close()
		wg.Wait()
		return ctxTimeout.Err()
	}
}

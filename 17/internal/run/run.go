package run

import (
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"telnet/internal/config"
)

// Run starts the client
func Run(opt *config.Options) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	conn, err := net.DialTimeout("tcp", opt.Address, opt.Timeout)
	if err != nil {
		return err
	}
	defer func() { _ = conn.Close() }()

	go func() {
		if _, err := io.Copy(conn, os.Stdin); err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
			cancel()
		}
	}()

	go func() {
		if _, err := io.Copy(os.Stdout, conn); err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
			cancel()
		}
	}()

	<-ctx.Done()
	return nil
}

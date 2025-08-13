package run

import (
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"telnet/internal/config"
)

func Run(opt *config.Options) error {
	conn, err := net.DialTimeout("tcp", opt.Address, opt.Timeout)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer func() { cancel(); _ = conn.Close() }()

	go func() {
		if _, err := io.Copy(os.Stdout, conn); err != nil {
			_, _ = fmt.Fprintln(os.Stderr, "read error:", err)
			cancel()
		}
	}()

	go func() {
		if _, err := io.Copy(conn, os.Stdin); err != nil {
			_, _ = fmt.Fprintln(os.Stderr, "read error:", err)
			cancel()
		}
	}()

	<-ctx.Done()
	return nil
}

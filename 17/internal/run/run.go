package run

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"syscall"
	"telnet/internal/config"
	"telnet/internal/telnet"
)

func Run(opt *config.Options) error {
	conn, err := net.DialTimeout("tcp", opt.Address, opt.Timeout)
	if err != nil {
		return err
	}
	defer conn.Close()
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()
	cli := telnet.New(conn)
	if err = cli.Run(ctx); err != nil {
		if errors.Is(err, io.EOF) {
			return nil
		}
		fmt.Fprintln(os.Stderr, "telnet run error:", err)
	}
	return nil
}

package telnet

import (
	"context"
	"io"
	"net"
	"os"
)

const (
	SE = 240
	SB = 249 + iota
	WILL
	WONT
	DO
	DONT
	IAC
	ECHO = 1
)

type Client struct {
	conn net.Conn
	out  io.Writer
	err  io.Writer
	in   io.Reader
	echo bool
}

func (c *Client) Run(ctx context.Context) error {
	go func() {
		if err := c.fromConn(ctx); err != nil && ctx.Err() == nil {
			_, _ = c.err.Write([]byte("fromConn error: " + err.Error() + "\n"))
		}
	}()

	return c.toConn(ctx)
}

func New(conn net.Conn) *Client {
	return &Client{
		conn: conn,
		out:  os.Stdout,
		err:  os.Stderr,
		in:   os.Stdin,
		echo: true,
	}
}

func (c *Client) toConn(ctx context.Context) error {
	buf := make([]byte, 1024)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		_, err := c.in.Read(buf)
		if err != nil {
			return err
		}

		if c.echo {
			if _, err := c.out.Write(buf); err != nil {
				return err
			}
		}

		if buf[0] == IAC {
			if _, err := c.conn.Write([]byte{IAC, IAC}); err != nil {
				return err
			}
		} else {
			if _, err := c.conn.Write(buf); err != nil {
				return err
			}
		}
	}
}

func (c *Client) fromConn(ctx context.Context) error {
	buf := make([]byte, 1)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		_, err := c.conn.Read(buf)
		if err != nil {
			return err
		}

		if buf[0] == IAC {
			cmd := make([]byte, 2)
			if _, err := io.ReadFull(c.conn, cmd); err != nil {
				return err
			}

			switch cmd[0] {
			case WILL:
				if cmd[1] == ECHO {
					c.echo = false
					c.conn.Write([]byte{IAC, DO, ECHO})
				}
			case WONT:
				if cmd[1] == ECHO {
					c.echo = true
					c.conn.Write([]byte{IAC, DONT, ECHO})
				}
			default:
			}
		} else {
			_, err := c.out.Write(buf)
			if err != nil {
				return err
			}
		}
	}
}

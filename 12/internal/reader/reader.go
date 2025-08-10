package reader

import (
	"io"
	"os"
)

// Open determines which file to open.
func Open(args []string) (io.ReadCloser, error) {
	if len(args) == 0 {
		return os.Stdin, nil
	}
	f, err := os.Open(args[0])
	if err != nil {
		return nil, err
	}
	return f, nil
}

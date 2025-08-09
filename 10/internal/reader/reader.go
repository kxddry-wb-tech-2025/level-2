package reader

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

// OpenInput returns an opened file
func OpenInput(args []string) (io.ReadCloser, error) {
	if len(args) == 0 {
		// os.Stdin is already a ReadCloser
		return os.Stdin, nil
	}
	f, err := os.Open(args[0])
	if err != nil {
		return nil, fmt.Errorf("open file: %w", err)
	}
	return f, nil
}

// ReadAllLines reads all lines from r without imposing Scanner's token limit.
func ReadAllLines(r io.Reader) ([]string, error) {
	rd := bufio.NewReaderSize(r, 64*1024)
	var lines []string
	for {
		line, err := rd.ReadString('\n')
		if err == io.EOF {
			if len(line) > 0 {
				// last line without newline
				if line[len(line)-1] == '\r' {
					line = line[:len(line)-1]
				}
				lines = append(lines, line)
			}
			break
		}
		if err != nil {
			return nil, err
		}
		// normalize: strip trailing '\n' and optional '\r'
		if n := len(line); n > 0 {
			if line[n-1] == '\n' {
				line = line[:n-1]
				n--
			}
			if n > 0 && line[n-1] == '\r' {
				line = line[:n-1]
			}
		}
		lines = append(lines, line)
	}
	return lines, nil
}

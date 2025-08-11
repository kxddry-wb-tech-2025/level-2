package cmd

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"cut/internal/config"
)

// Process reads lines from r, selects fields according to cfg and writes results to w.
func Process(r io.Reader, w io.Writer, cfg config.Options) error {
	scanner := bufio.NewScanner(r)
	outw := bufio.NewWriter(w)
	defer func() { _ = outw.Flush() }()

	for scanner.Scan() {
		line := scanner.Text()
		words := strings.Split(line, string(cfg.Delimiter))
		if cfg.SepOnly && len(words) == 1 {
			continue
		}

		var out []string
		if cfg.ShowAll {
			out = words
		} else {
			for i := range words {
				if _, ok := cfg.Fields[i+1]; ok {
					out = append(out, words[i])
				}
			}
		}
		_, _ = fmt.Fprintln(outw, strings.Join(out, "\t"))
	}

	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}

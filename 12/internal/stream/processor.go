package stream

import (
	"bufio"
	"fmt"
	"grep/internal/config"
	"io"
	"regexp"
	"strings"
)

// Processor implements the run.StreamProcesser interface
type Processor struct {
	opt *config.Flags
}

// NewProcessor creates a Processor.
func NewProcessor(opt *config.Flags) *Processor {
	return &Processor{opt: opt}
}

type prevLine struct {
	num  int
	text string
}

const maxToken = 10 * 1024 * 1024

// ProcessStream processes stream according to the options
func (p *Processor) ProcessStream(r io.Reader, pattern string, re *regexp.Regexp) (int, error) {
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 4096), maxToken)

	prev := make([]prevLine, 0, p.opt.Before+1)
	idx := 0
	trailing := 0
	matchCount := 0
	last := 0

	_print := func(num int, text string) {
		prefix := ""
		if p.opt.PrintNumbers {
			prefix = fmt.Sprintf("%d\t", num)
		}
		fmt.Println(prefix + text)
		last = num
	}

	for scanner.Scan() {
		line := scanner.Text()
		idx++

		match := p.isMatch(line, pattern, re)
		if p.opt.Invert {
			match = !match
		}

		if match {
			matchCount++
			if !p.opt.OnlyCount {
				for _, pl := range prev {
					if pl.num > last {
						_print(pl.num, pl.text)
					}
				}

				_print(idx, line)
			}

			trailing = max(trailing, p.opt.After)
		} else {
			if trailing > 0 {
				if !p.opt.OnlyCount {
					_print(idx, line)
				}
				trailing--
			}
		}

		if p.opt.Before > 0 {
			prev = append(prev, prevLine{idx, line})
			if len(prev) > p.opt.Before {
				prev = prev[1:]
			}
		} else {
			if len(prev) > 0 {
				prev = prev[:0]
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return matchCount, err
	}
	return matchCount, nil

}

func (p *Processor) isMatch(line, pattern string, re *regexp.Regexp) bool {
	if p.opt.FixedString {
		if p.opt.IgnoreCase {
			return strings.Contains(strings.ToLower(line), strings.ToLower(pattern))
		}
		return strings.Contains(line, pattern)
	}
	if re == nil {
		return false
	}
	return re.FindStringIndex(line) != nil
}

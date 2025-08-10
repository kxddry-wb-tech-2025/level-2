package run

import (
	"fmt"
	"grep/internal/config"
	"grep/internal/reader"
	"io"
	"os"
	"regexp"
)

// StreamProcesser can process stream.
type StreamProcesser interface {
	ProcessStream(r io.Reader, pattern string, re *regexp.Regexp) (int, error)
}

func compileRegexp(pattern string, ignoreCase bool) (*regexp.Regexp, error) {
	if ignoreCase {
		pattern = "(?i)" + pattern
	}
	return regexp.Compile(pattern)
}

// Run runs the CLI tool.
func Run(args []string, opt config.Flags, sp StreamProcesser) (err error) {
	pattern := args[0]
	r, err := reader.Open(args[1:])
	defer func() {
		err = r.Close()
	}()

	if err != nil {
		return err
	}

	var re *regexp.Regexp
	if !opt.FixedString {
		re, err = compileRegexp(args[0], opt.IgnoreCase)
		if err != nil {
			println("regex compile error:", err.Error())
			os.Exit(2)
		}
	}

	count, err := sp.ProcessStream(r, pattern, re)
	if err != nil {
		println("process error:", err.Error())
		os.Exit(3)
	}

	if opt.OnlyCount {
		fmt.Println(count)
	}

	foundAny := count > 0

	if foundAny {
		os.Exit(0)
	}
	os.Exit(1)
	return nil
}

package builtins

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Cd changes directory
func Cd(in io.Reader, out io.Writer, args ...string) bool {
	var path string
	if len(args) > 0 {
		path = args[0]
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err.Error())
			return false
		}
		err = os.Chdir(home)
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err.Error())
			return false
		}
		return true
	}

	if strings.HasPrefix(path, "~") {
		home, err := os.UserHomeDir()
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err.Error())
			return false
		}
		path = filepath.Join(home, path[1:])
	}
	if err := os.Chdir(path); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, err.Error())
		return false
	}

	return true
}

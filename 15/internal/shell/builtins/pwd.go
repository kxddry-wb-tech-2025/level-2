package builtins

import (
	"fmt"
	"io"
	"os"
)

// Pwd shows the current working directory
func Pwd(stdout io.Writer) bool {
	pwd, err := os.Getwd()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "failed to get current working directory")
		return false
	}
	_, _ = fmt.Fprintln(stdout, pwd)
	return true
}

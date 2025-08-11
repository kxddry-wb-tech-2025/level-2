package builtins

import (
	"fmt"
	"io"
	"os"
	"os/exec"
)

// Ps shows the currently running processes
func Ps(in io.Reader, out io.Writer, args ...string) bool {
	cmd := exec.Command("ps", args...)
	cmd.Stdin = in
	cmd.Stdout = out
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "ps: %v\n", err)
		return false
	}
	return true
}

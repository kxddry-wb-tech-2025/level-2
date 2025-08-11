package builtins

import (
	"fmt"
	"os"
	"strconv"
	"syscall"
)

// Kill sends a syscall signal to process to kill it
func Kill(args ...string) bool {
	if len(args) == 0 {
		_, _ = fmt.Fprintln(os.Stderr, "kill: missing PID argument")
		return false
	}

	pid, err := strconv.Atoi(args[0])
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "kill: invalid PID '%s'\n", args[1])
		return false
	}

	process, err := os.FindProcess(pid)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "kill: process %d not found: %v\n", pid, err)
		return false
	}

	if err := process.Signal(syscall.SIGTERM); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "kill: failed to kill process %d: %v\n", pid, err)
		return false
	}

	return true
}

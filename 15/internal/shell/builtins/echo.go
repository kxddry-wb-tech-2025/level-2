package builtins

import (
	"fmt"
	"io"
	"strings"
)

// Echo displays something user asked for on the screen
func Echo(in io.Reader, out io.Writer, args ...string) bool {
	if len(args) == 0 {
		_, _ = fmt.Fprintln(out)
		return true
	}
	_, _ = fmt.Fprintf(out, "%s\n", strings.Join(args, " "))
	return true
}

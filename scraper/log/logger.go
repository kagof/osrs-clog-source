package log

import (
	"fmt"
	"os"
)

var Writer = os.Stderr // by default, log to stderr

func Printf(format string, a ...any) {
	_, _ = fmt.Fprintf(Writer, format, a...)
}

func Println(msg string) {
	_, _ = fmt.Fprintln(Writer, msg)
}

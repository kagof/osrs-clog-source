package log

import (
	"fmt"
	"io"
	"os"
)

var Writer io.Writer = os.Stderr // by default, log to stderr

func Disable() {
	Writer = io.Discard
}

func Printf(format string, a ...any) {
	_, _ = fmt.Fprintf(Writer, format, a...)
}

func Println(msg string) {
	_, _ = fmt.Fprintln(Writer, msg)
}

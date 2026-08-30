package log

import (
	"fmt"
	"io"
	"os"
)

var Writer io.Writer = os.Stderr // by default, log to stderr
var debugEnabled = false

func EnableDebug() {
	debugEnabled = true
}

func Disable() {
	Writer = io.Discard
}

func Printf(format string, a ...any) {
	_, _ = fmt.Fprintf(Writer, format, a...)
}

func Println(msg string) {
	_, _ = fmt.Fprintln(Writer, msg)
}

func Debugf(format string, a ...any) {
	if debugEnabled {
		_, _ = fmt.Fprintf(Writer, "DBG: "+format, a...)
	}
}

func Debugln(msg string) {
	if debugEnabled {
		_, _ = fmt.Fprintln(Writer, "DBG: "+msg)
	}
}

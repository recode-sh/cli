package system

import (
	"fmt"
	"os"

	"github.com/recode-sh/cli/internal/constants"
)

type Logger struct{}

func NewLogger() Logger {
	return Logger{}
}

func (Logger) Info(format string, v ...interface{}) {
	fmt.Fprintf(os.Stderr, constants.Cyan(format)+"\n", v...)
}

func (Logger) Warning(format string, v ...interface{}) {
	fmt.Fprintf(os.Stderr, constants.Yellow(format)+"\n", v...)
}

func (Logger) Error(format string, v ...interface{}) {
	fmt.Fprintf(os.Stderr, constants.Red(format)+"\n", v...)
}

func (Logger) Log(format string, v ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", v...)
}

func (Logger) LogNoNewline(format string, v ...interface{}) {
	fmt.Fprintf(os.Stderr, format, v...)
}

func (l Logger) Write(p []byte) (n int, err error) {
	l.Log(string(p))
	return len(p), nil
}

package system

import (
	"fmt"
	"io"
)

type Displayer struct{}

func NewDisplayer() Displayer {
	return Displayer{}
}

func (Displayer) Display(w io.Writer, format string, args ...interface{}) {
	fmt.Fprintf(w, format, args...)
}

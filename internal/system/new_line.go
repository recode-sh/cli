package system

import "runtime"

var NewLineChar = "\n"

func init() {
	if runtime.GOOS == "windows" {
		NewLineChar = "\r\n"
	}
}

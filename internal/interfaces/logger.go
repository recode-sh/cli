package interfaces

type Logger interface {
	Info(format string, v ...interface{})
	Warning(format string, v ...interface{})
	Error(format string, v ...interface{})
	Log(format string, v ...interface{})
	LogNoNewline(format string, v ...interface{})
	Write(p []byte) (n int, err error)
}

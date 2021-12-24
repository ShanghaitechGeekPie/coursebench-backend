package errors

type ErrorLevel int

const (
	Silent ErrorLevel = iota + 1
	Info
	Error
	Fatal
)

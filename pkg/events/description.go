package events

import "go.uber.org/zap/zapcore"

type description struct {
	level   zapcore.Level
	name    string
	message string
	errno   string
	status  int
}

func (d *description) Level() zapcore.Level {
	return d.level
}

func (d *description) Name() string {
	return d.name
}

func (d *description) Message() string {
	return d.message
}

func (d *description) Errno() string {
	return d.errno
}
func (d *description) HttpStatus() int {
	return d.status
}

var errorDescriptionList = make([]*description, 0)

var errnoGeneratorInstance = newErrnoSequenceGenerator()

func createDescription(name, errorMessage string, level zapcore.Level, status int) (errorDescription *description) {
	errorDescription = &description{
		level:   level,
		name:    name,
		message: errorMessage,
		errno:   errnoGeneratorInstance.Value(),
		status:  status,
	}
	errnoGeneratorInstance = errnoGeneratorInstance.Next()
	errorDescriptionList = append(errorDescriptionList, errorDescription)
	return
}

package errors

import (
	"fmt"
	"github.com/pkg/errors"
	"time"
)

type description struct {
	Option
	name    string
	message string
	errno   string
	status  int
}

func (d *description) Name() string {
	return d.name
}

func (d *description) Error() string {
	return d.message
}

func (d *description) Errno() string {
	return d.errno
}

func (d *description) StatusCode() int {
	return d.status
}

var errorDescriptionList = make([]*description, 0)

var errnoGeneratorInstance = newErrnoSequenceGenerator()

func createDescription(name, errorMessage string, option Option, status ...int) (errorDescription *description) {
	statusCode := 500
	for _, statusCode = range status {
		if statusCode < 100 || statusCode > 600 {
			panic(fmt.Sprintf("status code %d is not in the range of 100-600", statusCode))
		}
	}
	errorDescription = &description{
		name:    name,
		Option:  option,
		message: errorMessage,
		errno:   errnoGeneratorInstance.Value(),
		status:  statusCode,
	}
	errnoGeneratorInstance = errnoGeneratorInstance.Next()
	errorDescriptionList = append(errorDescriptionList, errorDescription)
	return
}

// New 返回一个通过ErrorDescription创建的Error
// New 也会记录堆栈信息
// description 要显示的错误
// options 返回错误的参数，例如是否在日志中输出
// 此方法多用于抛出逻辑错误, 例如用户名或密码错误, 因为在这个过程中, 没有错误变量产生
func New(description *description) error {
	return &userErrorImpl{
		description: description,
		err:         errors.New(description.message),
		time:        time.Now(),
	}
}

// Wrap 传入一个err, 返回一个通过ErrorDescription包装的Error
// Wrap 也会记录堆栈信息
// err 原始error
// description 要显示的错误
// options 返回错误的参数，例如是否在日志中输出
// 此方法多用于抛出有错误变量产生的错误
func Wrap(err error, description *description) error {
	if err == nil {
		return nil
	}

	if _, ok := err.(*userErrorImpl); ok {
		return err
	}

	return &userErrorImpl{
		description: description,
		err:         errors.WithStack(err),
		time:        time.Now(),
	}
}

type UserError interface {
	Errno() string
	Error() string
	Name() string
	Cause() error
	Time() time.Time
	Stacktrace() string
	StatusCode() int
	Option
}

type userErrorImpl struct {
	*description
	err  error
	time time.Time
}

func (err *userErrorImpl) Cause() error {
	return err.err
}

func (err *userErrorImpl) Stacktrace() string {
	return fmt.Sprintf("%+v", err.err)
}

func (err *userErrorImpl) Time() time.Time {
	return err.time
}

func Is(err error, target error) bool {
	if err, ok := err.(*userErrorImpl); ok {
		return errors.Is(err.description, target)
	}
	if target, ok := target.(*userErrorImpl); ok {
		return errors.Is(err, target.description)
	}
	return errors.Is(err, target)
}

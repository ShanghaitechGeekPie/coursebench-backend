package events

import (
	"context"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"time"
)

// New 返回一个通过ErrorDescription创建的Error
// New 也会记录堆栈信息
// description 要显示的错误
// options 返回错误的参数，例如是否在日志中输出
// 此方法多用于抛出逻辑错误, 例如用户名或密码错误, 因为在这个过程中, 没有错误变量产生
func New(description *description) *AttributedEvent {
	return &AttributedEvent{
		description: description,
		cause:       errors.New(description.name),
		time:        time.Now(),
		fields: []zap.Field{
			zap.String("cause", description.name),
		},
	}
}

// Wrap 传入一个err, 返回一个通过ErrorDescription包装的Error
// Wrap 也会记录堆栈信息
// err 原始error
// description 要显示的错误
// options 返回错误的参数，例如是否在日志中输出
// 此方法多用于抛出有错误变量产生的错误
func Wrap(err error, description *description) *AttributedEvent {
	if err == nil {
		return nil
	}

	if err, ok := err.(*AttributedError); ok {
		return err.Event
	}

	return &AttributedEvent{
		description: description,
		cause:       errors.WithStack(err),
		time:        time.Now(),
		fields: []zap.Field{
			zap.String("cause", err.Error()),
		},
	}
}

type AttributedEvent struct {
	*description
	cause  error
	time   time.Time
	fields []zap.Field
	logged bool
}

func (event *AttributedEvent) Add(fields ...zap.Field) *AttributedEvent {
	oldFields := event.fields
	if oldFields == nil {
		oldFields = make([]zap.Field, 0, len(fields))
	}
	newFields := append(oldFields, fields...)
	return &AttributedEvent{
		description: event.description,
		cause:       event.cause,
		time:        event.time,
		fields:      newFields,
	}
}

func (event *AttributedEvent) Cause() error {
	return event.cause
}

func (event *AttributedEvent) Description() *description {
	return event.description
}

func (event *AttributedEvent) Fields() []zap.Field {
	return event.fields
}

func (event *AttributedEvent) Time() time.Time {
	return event.time
}

func (event *AttributedEvent) IsError() bool {
	if event.description == nil {
		return false
	}
	return event.description.level > zap.WarnLevel
}

func (event *AttributedEvent) ToError() error {
	if event == nil {
		return nil
	}
	return &AttributedError{
		Event: event,
	}
}

func (event *AttributedEvent) Log(ctx context.Context) *AttributedEvent {
	return logger.WithContext(ctx).Log(event)
}

type AttributedError struct {
	Event *AttributedEvent
}

func (err *AttributedError) Error() string {
	return err.Event.Cause().Error()
}

func HasSpecificDescription(event *AttributedEvent, description *description) bool {
	return event.description == description
}

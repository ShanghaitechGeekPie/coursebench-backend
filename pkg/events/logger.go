package events

import (
	"context"
	"coursebench-backend/internal/config"
	"encoding/json"
	"fmt"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

type CtxLogger struct {
	logger *zap.Logger
	ctx    context.Context
}

var logger *CtxLogger = nil

func initLogger(InDevelopment bool) {
	logger = new(CtxLogger)
	var err error
	if InDevelopment {
		logger.logger, err = zap.NewDevelopment()
	} else {
		logger.logger, err = zap.NewProduction()
	}
	if err != nil {
		panic(err)
	}
}

func (l *CtxLogger) WithContext(ctx context.Context) *CtxLogger {
	newLogger := *l
	newLogger.ctx = ctx
	return &newLogger
}

func (l *CtxLogger) Log(event *AttributedEvent) *AttributedEvent {
	if event == nil {
		return nil
	}
	if event.logged {
		return event
	}
	event.logged = true

	if event.description != CacheMiss {
		level := event.Description().Level()
		if ce := l.logger.Check(level, event.Name()); ce != nil {
			ce.Write(event.Fields()...)
		}
	}
	if l.ctx != nil {
		span := trace.SpanFromContext(l.ctx)
		if span.IsRecording() {
			encoder := zapcore.NewJSONEncoder(zapcore.EncoderConfig{
				LevelKey:       "log.severity",
				TimeKey:        "log.timestamp",
				MessageKey:     "exception.message",
				StacktraceKey:  "exception.stacktrace",
				EncodeLevel:    zapcore.CapitalLevelEncoder,
				EncodeTime:     zapcore.ISO8601TimeEncoder,
				EncodeDuration: zapcore.StringDurationEncoder,
				EncodeCaller:   zapcore.ShortCallerEncoder,
			})

			encoded, err := encoder.EncodeEntry(zapcore.Entry{
				Level:   event.Level(),
				Time:    event.Time(),
				Message: event.Name(),
				Stack:   fmt.Sprintf("%+v", event.Cause()),
			}, event.Fields())
			if err != nil {
				l.logger.Error("Unable to encode entry", zap.Error(err))
			} else {
				defer encoded.Free()
				kvs := make(map[string]interface{})
				kvs["log.message"] = event.Name() + " " + string(encoded.Bytes())
				err = json.Unmarshal(encoded.Bytes(), &kvs)
				if err != nil {
					panic(err)
				}
				attrs := make([]attribute.KeyValue, 0, len(kvs))
				for key, value := range kvs {
					bytes, _ := json.Marshal(value)
					attrs = append(attrs, attribute.String(key, string(bytes)))
				}
				span.AddEvent(event.Name(), trace.WithTimestamp(time.Now()), trace.WithStackTrace(true), trace.WithAttributes(attrs...))
				if event.IsError() {
					span.SetStatus(codes.Error, event.Description().Message())
				}
			}
		}
	}
	return event
}

func init() {
	initLogger(config.GlobalConf.InDevelopment)
}

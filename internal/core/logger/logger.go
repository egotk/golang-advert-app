package corelogger

import "context"

type loggerContextKey struct{}

var key = loggerContextKey{}

type Logger interface {
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	DPanic(msg string, fields ...Field)
	Panic(msg string, fields ...Field)
	Fatal(msg string, fields ...Field)

	With(field ...Field) Logger
	Close()
}

func ToContext(ctx context.Context, log Logger) context.Context {
	return context.WithValue(ctx, key, log)
}

func FromContext(ctx context.Context) Logger {
	log, ok := ctx.Value(key).(Logger)
	if !ok {
		panic("no logger in context")
	}

	return log
}

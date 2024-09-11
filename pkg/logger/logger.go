package logger

import (
	"context"
	"errors"
)

const (
	// ErrorKey is the key of error used in logging messages.
	ErrorKey = "err"
)

var (
	// ErrNotRegistered means needing to call Register() once in the beginning.
	ErrNotRegistered = errors.New("not registered")
)

// Fields type, used to pass to `WithFields`.
type Fields map[string]interface{}

type LoggerBase interface {
	// Flush flushes any buffered log entries.
	// Applications should take care to call Sync before exiting.
	Flush() error

	// Debug logs a message at DebugLevel. The message includes any fields passed
	// at the log site, as well as any fields accumulated on the logger.
	Debug(message string, options ...LoggerOptions)
	// Info logs a message at InfoLevel. The message includes any fields passed
	// at the log site, as well as any fields accumulated on the logger.
	Info(message string, options ...LoggerOptions)
	// Warn logs a message at WarnLevel. The message includes any fields passed
	// at the log site, as well as any fields accumulated on the logger.
	Warn(message string, options ...LoggerOptions)
	// Error logs a message at ErrorLevel. The message includes any fields passed
	// at the log site, as well as any fields accumulated on the logger.
	Error(message string, options ...LoggerOptions)
	// Panic logs a message at PanicLevel. The message includes any fields passed
	// at the log site, as well as any fields accumulated on the logger.
	//
	// The logger then panics, even if logging at PanicLevel is disabled.
	Panic(message string, options ...LoggerOptions)
	// Fatal logs a message at FatalLevel. The message includes any fields passed
	// at the log site, as well as any fields accumulated on the logger.
	//
	// The logger then calls os.Exit(1), even if logging at FatalLevel is disabled.
	Fatal(message string, options ...LoggerOptions)
}

type LoggerWithCtx interface {
	LoggerBase
}

type Logger interface {
	LoggerBase

	// Ctx returns a new logger with the context.
	Ctx(ctx context.Context) LoggerWithCtx
}

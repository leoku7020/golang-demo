package logger

import "context"

func newEmpty() Logger {
	return &empty{}
}

func newEmptyWithCtx() LoggerWithCtx {
	return &empty{}
}

type empty struct{}

func (e *empty) Ctx(ctx context.Context) LoggerWithCtx {
	return &empty{}
}

func (e *empty) Flush() error {
	return nil
}

func (e *empty) Debug(message string, options ...LoggerOptions) {}

func (e *empty) Info(message string, options ...LoggerOptions) {}

func (e *empty) Warn(message string, options ...LoggerOptions) {}

func (e *empty) Error(message string, options ...LoggerOptions) {}

func (e *empty) Panic(message string, options ...LoggerOptions) {}

func (e *empty) Fatal(message string, options ...LoggerOptions) {}

package logger

import (
	"context"
	"fmt"
	"os"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func convertMapToFields[V any](m map[string]V) []zapcore.Field {
	var fields []zapcore.Field
	for k, v := range m {
		fields = append(fields, zap.Any(k, v))
	}

	return fields
}

func injectTags(lg *zap.Logger, tags map[string]string) *zap.Logger {
	return lg.With(convertMapToFields(tags)...)
}

func newZapLogger(cfg Config) Logger {
	c := zap.NewProductionEncoderConfig()
	c.MessageKey = "message"
	enc := zapcore.NewJSONEncoder(c)

	if cfg.Development {
		c = zap.NewDevelopmentEncoderConfig()
		c.EncodeLevel = zapcore.CapitalColorLevelEncoder
		enc = zapcore.NewConsoleEncoder(c)
		cfg.Level = LevelDebug // override existing level with LevelDebug
	}

	cores := []zapcore.Core{zapcore.NewCore(
		enc,
		getLoggerWriter(),
		zapcore.Level(cfg.Level),
	)}

	// duplicates log entries into Airbrake client
	if cfg.AirbrakeCli != nil {
		zapCore, err := NewZapBrakeCore(zapcore.ErrorLevel, cfg.AirbrakeCli)
		if err != nil {
			panic(err)
		}

		cores = append(cores, zapCore)
	}

	// TODO: actually, AddCallerSkip(2) makes logger with Ctx() display well.
	// But, logger without Ctx() will show the previous stack of caller.
	lg := zap.New(zapcore.NewTee(cores...), zap.AddCaller(), zap.AddCallerSkip(2))

	return &zapLogger{
		lg:  injectTags(lg, cfg.CommonTags),
		ctx: context.TODO(),
	}
}

type zapLogger struct {
	lg  *zap.Logger
	ctx context.Context
}

func (z *zapLogger) Flush() error {
	return z.lg.Sync()
}

func (z *zapLogger) write(lvl zapcore.Level, msg string, o *loggerOptions) {
	entry := z.lg.Check(lvl, msg)
	if entry == nil {
		// do nothing
		return
	}

	m := map[string]interface{}{}
	if entry.Caller.Defined {
		// prevent naming collision
		m["runtime.file"] = entry.Caller.File
		m["runtime.line"] = entry.Caller.Line
		m["runtime.function"] = entry.Caller.Function
	}

	if fs, ok := fields(z.ctx); ok {
		for k, v := range fs {
			m[k] = v
		}
	}

	// overwrite fields in context with the one in option fields
	for k, v := range o.fields {
		m[k] = v
	}

	// overwrite fields in option fields with the one in option pairs
	for _, p := range o.pairs {
		m[p.Key] = p.Val
	}

	// overwrite fields with error
	if o.err != nil {
		m[ErrorKey] = o.err
	}

	entry.Write(convertMapToFields(m)...)
}

func (z *zapLogger) Debug(message string, options ...LoggerOptions) {
	z.write(zapcore.DebugLevel, message, applyLoggerOptions(options...))
}

func (z *zapLogger) Info(message string, options ...LoggerOptions) {
	z.write(zapcore.InfoLevel, message, applyLoggerOptions(options...))
}

func (z *zapLogger) Warn(message string, options ...LoggerOptions) {
	z.write(zapcore.WarnLevel, message, applyLoggerOptions(options...))
}

func (z *zapLogger) Error(message string, options ...LoggerOptions) {
	z.write(zapcore.ErrorLevel, message, applyLoggerOptions(options...))
}

func (z *zapLogger) Panic(message string, options ...LoggerOptions) {
	z.write(zapcore.PanicLevel, message, applyLoggerOptions(options...))
}

func (z *zapLogger) Fatal(message string, options ...LoggerOptions) {
	z.write(zapcore.FatalLevel, message, applyLoggerOptions(options...))
}

func (z *zapLogger) Ctx(ctx context.Context) LoggerWithCtx {
	return &zapLogger{
		lg:  z.lg,
		ctx: ctx,
	}
}

func getLoggerWriter() zapcore.WriteSyncer {
	if os.Getenv("LOGGER_FILE_PREFIX") != "" {
		logFile := fmt.Sprintf("%s-%s", os.Getenv("LOGGER_FILE_PREFIX"), "%Y-%m-%d.log")
		rotator, err := rotatelogs.New(
			logFile,
			rotatelogs.WithMaxAge(90*24*time.Hour),
			rotatelogs.WithRotationTime(time.Hour*24))
		if err != nil {
			fmt.Println("failed to log, err:", err)
		}
		return zapcore.AddSync(rotator)
	}

	return zapcore.AddSync(os.Stdout)
}

package logger

import (
	"context"
	"errors"
	"sync"

	"demo/pkg/envkit"
)

var (
	// registration
	regLogger   Logger
	regExporter = ExporterNone
	regDev      bool
	regLevel    Level

	// registerOnce limits registering once
	registerOnce = sync.Once{}
)

// Registered returns if the logger was registered or not in the beginning
func Registered() bool {
	return regLogger != nil
}

func debugEnabled() bool {
	if regDev || regLevel == LevelDebug {
		return true
	}

	return false
}

func Register(cfg Config) {
	registerOnce.Do(func() {
		if _, ok := _ExporterMap[cfg.Exporter]; !ok {
			panic(errors.New("no such exporter"))
		}

		regExporter = cfg.Exporter
		regDev = cfg.Development
		regLevel = cfg.Level

		// set default value
		if cfg.CommonTags == nil {
			cfg.CommonTags = map[string]string{}
		}
		// override the common tags with envkit if registered it before
		if envkit.Registered() {
			cfg.CommonTags["envkit.env"] = envkit.EnvNamespace()
			cfg.CommonTags["envkit.service"] = envkit.ServiceName()
			cfg.CommonTags["envkit.project"] = envkit.ProjectName()
			cfg.CommonTags["envkit.pod"] = envkit.PodName()
		}

		switch regExporter {
		case ExporterZap:
			regLogger = newZapLogger(cfg)
		default:
			regLogger = newEmpty()
		}
	})
}

func Flush() error {
	if !Registered() {
		return ErrNotRegistered
	}

	return regLogger.Flush()
}

func Debug(message string, options ...LoggerOptions) {
	if !Registered() {
		return
	}

	regLogger.Debug(message, options...)
}

func Info(message string, options ...LoggerOptions) {
	if !Registered() {
		return
	}

	regLogger.Info(message, options...)
}

func Warn(message string, options ...LoggerOptions) {
	if !Registered() {
		return
	}

	regLogger.Warn(message, options...)
}

func Error(message string, options ...LoggerOptions) {
	if !Registered() {
		return
	}

	regLogger.Error(message, options...)
}

func Panic(message string, options ...LoggerOptions) {
	if !Registered() {
		return
	}

	regLogger.Panic(message, options...)
}

func Fatal(message string, options ...LoggerOptions) {
	if !Registered() {
		return
	}

	regLogger.Fatal(message, options...)
}

func Ctx(ctx context.Context) LoggerWithCtx {
	if !Registered() {
		return newEmptyWithCtx()
	}

	return regLogger.Ctx(ctx)
}

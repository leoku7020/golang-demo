package logger

import (
	newrelic "github.com/newrelic/go-agent/v3/newrelic"
)

// NewNRLogger implements newrelic.Logger with pkg/logger
func NewNRLogger() newrelic.Logger {
	return &nrLggr{}
}

// ConfigNRLogger configures the newrelic.Application to send log messsages to the
// provided zap logger.
func ConfigNRLogger() newrelic.ConfigOption {
	return newrelic.ConfigLogger(NewNRLogger())
}

type nrLggr struct{}

func (lg *nrLggr) Error(msg string, context map[string]interface{}) {
	Error(msg, WithFields(Fields(context)))
}

func (lg *nrLggr) Warn(msg string, context map[string]interface{}) {
	Warn(msg, WithFields(Fields(context)))
}

func (lg *nrLggr) Info(msg string, context map[string]interface{}) {
	Info(msg, WithFields(Fields(context)))
}

func (lg *nrLggr) Debug(msg string, context map[string]interface{}) {
	Debug(msg, WithFields(Fields(context)))
}

func (lg *nrLggr) DebugEnabled() bool {
	return debugEnabled()
}

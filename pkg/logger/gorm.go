package logger

import (
	"context"
	"fmt"
	"time"

	gormLogger "gorm.io/gorm/logger"
)

func NewGormLogger() gormLogger.Interface {
	if !Registered() {
		return nil
	}

	return &gormLggr{}
}

type gormLggr struct{}

func (lg *gormLggr) LogMode(l gormLogger.LogLevel) gormLogger.Interface {
	return lg
}

func (lg *gormLggr) Info(ctx context.Context, msg string, data ...interface{}) {
	Ctx(ctx).Info(fmt.Sprintf(msg, data...))
}

func (lg *gormLggr) Warn(ctx context.Context, msg string, data ...interface{}) {
	Ctx(ctx).Warn(fmt.Sprintf(msg, data...))
}

func (lg *gormLggr) Error(ctx context.Context, msg string, data ...interface{}) {
	Ctx(ctx).Error(fmt.Sprintf(msg, data...))
}

func (lg *gormLggr) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	explain, affectedRows := fc()
	Ctx(ctx).Debug(explain,
		WithError(err),
		WithFields(Fields{
			"begin":        begin,
			"affectedRows": affectedRows,
		}),
	)
}

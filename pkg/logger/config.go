//go:generate go-enum -f=$GOFILE --nocase

package logger

import (
	"github.com/airbrake/gobrake/v5"
)

// Exporter is an enumeration of exporters.
/*
ENUM(
None // Not existed
Zap // Coordinate logs with Zap (https://github.com/uber-go/zap)
)
*/
type Exporter int32

// Config defines details used in Register()
type Config struct {
	// Exporter defines the way to coordinate logs.
	Exporter Exporter
	// Development `true` is used in development env or test case only
	Development bool
	// Level means Log Level used in different envs.
	Level Level
	// AirbrakeCli means Airbrake client used to broadcast logs to Airbrake
	AirbrakeCli *gobrake.Notifier
	// CommonTags add common tags (key-value pairs) for each log
	CommonTags map[string]string
}

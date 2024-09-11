package initkit

import (
	"demo/pkg/logger"
	"os"
	"strconv"
)

func InitLogger() func() error {
	development := os.Getenv("LOGGER_DEVELOPMENT")
	level, _ := strconv.ParseInt(os.Getenv("LOGGER_LEVEL"), 10, 64)
	logger.Register(logger.Config{
		Exporter:    logger.ExporterZap,
		Development: development == "true",
		Level:       logger.Level(level),
	})

	logger.Info("Logger done", logger.WithFields(logger.Fields{
		"level":       development,
		"development": level,
	}))

	return func() error {
		return logger.Flush()
	}
}

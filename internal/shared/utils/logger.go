package utils

import (
	"os"

	"go_boilerplate/internal/shared/config"

	"github.com/sirupsen/logrus"
)

// InitLogger initializes the logger with the given configuration
func InitLogger(cfg *config.Config) *logrus.Logger {
	logger := logrus.New()

	// Set log level
	switch cfg.Logger.Level {
	case "debug":
		logger.SetLevel(logrus.DebugLevel)
	case "info":
		logger.SetLevel(logrus.InfoLevel)
	case "warn":
		logger.SetLevel(logrus.WarnLevel)
	case "error":
		logger.SetLevel(logrus.ErrorLevel)
	default:
		logger.SetLevel(logrus.InfoLevel)
	}

	// Set log format
	if cfg.Logger.Format == "json" {
		logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
		})
	} else {
		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
		})
	}

	// Set output to stdout
	logger.SetOutput(os.Stdout)

	return logger
}

// WithFields creates a logger entry with fields
func WithFields(logger *logrus.Logger, fields logrus.Fields) *logrus.Entry {
	return logger.WithFields(fields)
}

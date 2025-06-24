package logger

import (
	"github.com/sirupsen/logrus"
)

// NewLogger initializes a new logger instance using Logrus
func NewLogger() *logrus.Logger {
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetLevel(logrus.InfoLevel)
	return log
}

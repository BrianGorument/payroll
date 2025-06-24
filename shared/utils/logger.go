package utils

import (
	"os"

	"github.com/sirupsen/logrus"
)

// NewLogger membuat instance baru dari logrus.Logger
func NewLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)
	return logger
}

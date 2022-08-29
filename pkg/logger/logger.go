package logger

import (
	"fmt"
	"github.com/sirupsen/logrus"
)

func NewWithConfig(logLevel string) (*logrus.Logger, error) {
	logger := New()
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		return nil, fmt.Errorf("failed to parse logger level %w", err)
	}
	logger.SetLevel(level)
	return logger, nil
}

func New() *logrus.Logger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	return logger
}

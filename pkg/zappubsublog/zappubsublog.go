package zappubsublog

import (
	"go.uber.org/zap"
)

type Logger struct {
	sugaredLogger *zap.SugaredLogger
}

func NewLogger() *Logger {
	logger, _ := zap.NewProduction() // zap.NewDevelopment()

	return &Logger{
		sugaredLogger: logger.Sugar(),
	}
}

func (logger *Logger) Infof(msg string, args ...any) {
	logger.sugaredLogger.Infof(msg, args...)
}

func (logger *Logger) Warn(msg string, args ...any) {
	logger.sugaredLogger.Warnf(msg, args...)
}

func (logger *Logger) Errorf(msg string, args ...any) {
	logger.sugaredLogger.Errorf(msg, args...)
}

func (logger *Logger) Debug(msg string, args ...any) {
	logger.sugaredLogger.Debugf(msg, args...)
}

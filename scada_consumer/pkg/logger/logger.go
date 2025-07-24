package logger

import (
	"go.uber.org/zap"
)

var log *zap.Logger

func InitLogger() {
	logger, _ := zap.NewProduction()
	log = logger
}

func Sync() {
	log.Sync()
}

func Info(msg string, fields ...zap.Field) {
	log.Info(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	log.Error(msg, fields...)
}

func Fatalf(msg string, fields ...zap.Field) {
	log.Fatal(msg, fields...)
}

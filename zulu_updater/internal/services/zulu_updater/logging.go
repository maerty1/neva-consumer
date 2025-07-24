package zulu_updater

import (
	"context"
	"fmt"
	"log"
	"runtime"
)

type LogLevel string
type LogCategory string

const (
	LogLevelInfo    LogLevel = "INFO"
	LogLevelWarning LogLevel = "WARNING"
	LogLevelError   LogLevel = "ERROR"
	LogLevelSuccess LogLevel = "SUCCESS"
)

func (s service) logUpdaterError(ctx context.Context, err error) {
	s.logUpdater(ctx, LogLevelError, err.Error())
}

func (s service) logUpdaterSuccess(ctx context.Context, message string) {
	if message == "" {
		message = "Успешный запрос"
	}
	s.logUpdater(ctx, LogLevelSuccess, message)
}

func (s service) logUpdaterWarning(ctx context.Context, err error) {
	s.logUpdater(ctx, LogLevelWarning, err.Error())
}

func (s service) logUpdaterInfo(ctx context.Context, message string) {
	s.logUpdater(ctx, LogLevelInfo, message)
}

func (s service) logUpdater(ctx context.Context, level LogLevel, message string) {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "unknown"
		line = 0
	}
	log.Println(fmt.Sprintf("\n	Level: %s;\n	Сообщение: %s;\n	Файл: %s;\n	Строка: %d", level, message, file, line))
}

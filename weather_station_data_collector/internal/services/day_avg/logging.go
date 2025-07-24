package day_avg

import (
	"context"
	"fmt"
	"log"
)

type LogLevel string

const (
	LogLevelInfo    LogLevel = "INFO"
	LogLevelWarning LogLevel = "WARNING"
	LogLevelError   LogLevel = "ERROR"
	LogLevelSuccess LogLevel = "SUCCESS"
)

func (s service) logInterrogateError(ctx context.Context, err error) {
	s.logInterrogate(ctx, LogLevelError, err.Error())
}

func (s service) logInterrogateSuccess(ctx context.Context) {
	s.logInterrogate(ctx, LogLevelSuccess, "Среднее значение успешно рассчитано и записано в бд")
}

func (s service) logInterrogateWarning(ctx context.Context, err error) {
	s.logInterrogate(ctx, LogLevelWarning, err.Error())
}

func (s service) logInterrogateInfo(ctx context.Context, message string) {
	s.logInterrogate(ctx, LogLevelInfo, message)
}

func (s service) logInterrogate(ctx context.Context, level LogLevel, message string) {
	log.Println(fmt.Sprintf("\n	Level: %s;\nСообщение: %s", level, message))
}

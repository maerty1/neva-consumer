package interrogator

import (
	"context"
	"fmt"
	"log"
)

type LogLevel string
type LogCategory string

const (
	LogLevelInfo    LogLevel = "INFO"
	LogLevelWarning LogLevel = "WARNING"
	LogLevelError   LogLevel = "ERROR"
	LogLevelSuccess LogLevel = "SUCCESS"
)

func (s *service) logInterrogateError(ctx context.Context, time string, err error) {
	s.logInterrogate(ctx, LogLevelError, time, err.Error())
}

func (s *service) logInterrogateSuccess(ctx context.Context, time string) {
	s.logInterrogate(ctx, LogLevelSuccess, time, "Успешный опрос")
}

func (s *service) logInterrogateWarning(ctx context.Context, time string, err error) {
	s.logInterrogate(ctx, LogLevelWarning, time, err.Error())
}

func (s *service) logInterrogateInfo(ctx context.Context, time string, message string) {
	s.logInterrogate(ctx, LogLevelInfo, time, message)
}

func (s *service) logInterrogate(ctx context.Context, level LogLevel, time string, message string) {
	log.Println(fmt.Sprintf("\n	Level: %s;\n	Опрашиваемое время: %s;\n	Сообщение: %s", level, time, message))
}

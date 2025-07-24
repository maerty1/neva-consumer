package poller

import (
	"context"
	"fmt"
)

type LogLevel string
type LogCategory string

const (
	LogLevelInfo    LogLevel = "INFO"
	LogLevelWarning LogLevel = "WARNING"
	LogLevelError   LogLevel = "ERROR"
)

func (s *service) logSyncError(ctx context.Context, accountID int, measurePointID int, err error) error {
	return s.logSync(ctx, LogLevelError, accountID, measurePointID, err.Error())
}

func (s *service) logSyncSuccess(ctx context.Context, accountID int, measurePointID int) error {
	return s.logSync(ctx, LogLevelInfo, accountID, measurePointID, "Успешный опрос")
}

func (s *service) logSync(ctx context.Context, level LogLevel, accountID int, measurePointID int, message string) error {
	logErr := s.measurePointsRepository.InsertSyncLog(ctx, accountID, measurePointID, string(level), message)
	if logErr != nil {
		fmt.Printf("Ошибка записи лога синхронизации для аккаунта %d: %v\n", accountID, logErr)
	}
	return nil
}

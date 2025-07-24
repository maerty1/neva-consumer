package poller

import (
	"context"
	"fmt"
	"lers_integration_service/internal/api_clients/lers"
	"lers_integration_service/internal/models"
	"log"
	"time"
)

const syncMeasurePointTimeout = time.Minute * 1

func (s *service) RunPoller(ctx context.Context) error {
	accountsToSync, err := s.measurePointsRepository.FindAccountsToSync(ctx)
	if err != nil {
		return err
	}

	if len(accountsToSync) == 0 {
		return fmt.Errorf("нет аккаунтов для синхронизации")
	}

	for _, account := range accountsToSync {
		fmt.Printf("Синхронизация аккаунта %d\n", account.ID)
		if err := s.pollMeasurePointsForAccount(ctx, account); err != nil {
			fmt.Printf("Ошибка синхронизации аккаунта %d: %v\n", account.ID, err)
		}
	}

	return nil
}

func (s *service) pollMeasurePointsForAccount(ctx context.Context, account models.AccountToSync) error {
	isRunning, err := s.lersApiClient.ArePollingsCurrentlyRunning(account.Token, account.ServerHost)
	if err != nil {
		return s.logSyncError(ctx, account.ID, 0, fmt.Errorf("не удалось получить текущую очередь опроса"))
	}

	if isRunning {
		s.waitForQueue(ctx, account, 0)
	}

	measurePoints, err := s.lersApiClient.GetMeasurePoints(account.Token, account.ServerHost)
	if err != nil {
		return s.logSyncError(ctx, account.ID, 0, fmt.Errorf("не удалось получить точки измерения: %w", err))
	}

	if len(measurePoints) == 0 {
		return s.logSyncError(ctx, account.ID, 0, fmt.Errorf("точки измерения не найдены"))
	}

	for _, point := range measurePoints {
		if err := s.syncMeasurePoint(ctx, account.ID, point, account); err != nil {
			continue
		}

	}

	return s.logSyncSuccess(ctx, account.ID, 0)
}

func (s *service) syncMeasurePoint(ctx context.Context, accountID int, point lers.MeasurePoint, account models.AccountToSync) error {
	if err := s.measurePointsRepository.InsertMeasurePoint(ctx, accountID, point.ID, point.DeviceID, point.Title, point.FullTitle, point.Address, point.SystemType); err != nil {
		return s.logSyncError(ctx, accountID, point.ID, fmt.Errorf("не удалось вставить точку измерения %d: %w", point.ID, err))
	}

	startDate, err := s.getStartDate(ctx, point.ID)
	if err != nil {
		return s.logSyncError(ctx, accountID, point.ID, err)
	}

	endDate := time.Now().Format(time.RFC3339)
	fmt.Printf("Опрос данных для точки измерения %d с %s по %s\n", point.ID, startDate, endDate)

	resp, err := s.lersApiClient.PollMeasurePoints(account.Token, account.ServerHost, []int{point.ID}, startDate, endDate, syncMeasurePointTimeout)
	if err != nil {
		return s.logSyncError(ctx, accountID, point.ID, fmt.Errorf("не удалось сделать опрос для точки измерения %d: %w", point.ID, err))
	}

	err = s.measurePointsRepository.InsertMeasurePointPollLog(ctx, resp.PollSessionID, point.ID, account.ID, resp.Status)
	if err != nil {
		return s.logSyncError(ctx, accountID, point.ID, fmt.Errorf("не удалось вставить лог опроса %d: %w", point.ID, err))
	}
	return nil
}

func (s *service) getStartDate(ctx context.Context, measurePointID int) (string, error) {
	lastDatetime, err := s.measurePointsRepository.GetLastMeasurePointDatetime(ctx, measurePointID)
	if err != nil {
		return "", fmt.Errorf("не удалось получить последний datetime для точки измерения %d: %w", measurePointID, err)
	}

	if lastDatetime == "" {
		return "2000-01-01T00:00:00Z", nil // Используем очень раннюю дату, если записей не существует.
	}

	lastTime, err := time.Parse(time.RFC3339, lastDatetime)
	if err != nil {
		return "", fmt.Errorf("не удалось распарсить последний datetime для точки измерения %d: %w", measurePointID, err)
	}

	return lastTime.Add(time.Hour).Format(time.RFC3339), nil
}

func (s *service) waitForQueue(ctx context.Context, account models.AccountToSync, pointID int) error {
	log.Println("в очереди остались задачи на опрос, ожидаем 60 секунд...")
	isRunning, err := s.lersApiClient.ArePollingsCurrentlyRunning(account.Token, account.ServerHost)
	if err != nil {
		return s.logSyncError(ctx, account.ID, pointID, fmt.Errorf("не удалось получить текущую очередь опроса"))
	}

	if isRunning {
		time.Sleep(time.Second * 60)
		s.waitForQueue(ctx, account, pointID)
	}
	return nil
}

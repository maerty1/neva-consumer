package retryer

import (
	"context"
	"fmt"
	"lers_integration_service/internal/models"
	"lers_integration_service/internal/repositories/measure_points"
	"log"
	"time"
)

const getPollSessionsTimeout = time.Second * 30
const pollMeasurePointsTimeout = time.Minute * 10

// RunRetryer запускает процесс повторного выполнения опросов для аккаунтов, требующих повторного опроса.
func (s *service) RunRetryer(ctx context.Context) error {
	accountsToRetry, err := s.measurePointsRepository.FindAccountsToSync(ctx)
	if err != nil {
		return err
	}

	if len(accountsToRetry) == 0 {
		return fmt.Errorf("нет аккаунтов для повторного опроса")
	}

	for _, account := range accountsToRetry {
		fmt.Printf("Повторный опрос для аккаунта %d\n", account.ID)
		if err := s.processAccountRetry(ctx, account); err != nil {
			fmt.Printf("Ошибка при повторном опросе аккаунта %d: %v\n", account.ID, err)
		}
	}
	fmt.Println("Процесс завершен")
	return nil
}

// processAccountRetry выполняет повторные опросы для конкретного аккаунта.
func (s *service) processAccountRetry(ctx context.Context, account models.AccountToSync) error {
	// Проверить, есть ли активные опросы. Если есть, ожидаем, если нет, продолжаем.
	isRunning, err := s.lersApiClient.ArePollingsCurrentlyRunning(account.Token, account.ServerHost)
	if err != nil {
		return s.logSyncError(ctx, account.ID, 0, fmt.Errorf("не удалось получить текущую очередь опроса: %v", err))
	}

	if isRunning {
		s.waitForPollingCompletion(ctx, account, 0)
	}

	endDate := time.Now().AddDate(0, 0, 1).Format("2006-01-02")
	startDate := time.Now().AddDate(0, 0, -7).Format("2006-01-02")
	// Получаем сессии опросов за последнюю неделю.
	pollSessions, err := s.lersApiClient.GetPollSessions(account.Token, account.ServerHost, startDate, endDate, pollMeasurePointsTimeout)
	if err != nil {
		return s.logSyncError(ctx, account.ID, 0, fmt.Errorf("не удалось получить сессии опроса: %v", err))
	}

	pollingsToRetry, err := s.measurePointsRepository.FindPollSessionsToRetry(ctx, account.ID)
	if err != nil {
		return s.logSyncError(ctx, account.ID, 0, fmt.Errorf("не удалось получить опросы для повторного выполнения: %v", err))
	}

	err = s.updateFailedPollSessions(ctx, pollSessions, pollingsToRetry)
	if err != nil {
		fmt.Println("Ошибка при обновлении статуса неудачных опросов")
		return s.logSyncError(ctx, account.ID, 0, fmt.Errorf("не удалось проверить и обновить poll сессии: %v", err))
	}

	retryPollSession, err := s.measurePointsRepository.FindRetryPollSessions(ctx, account.ID)
	if err != nil {
		return s.logSyncError(ctx, account.ID, 0, fmt.Errorf("не удалось получить poll сессии для повторного выполнения: %v", err))
	}
	err = s.updateFailedRetrySessions(ctx, pollSessions, retryPollSession)
	if err != nil {
		return s.logSyncError(ctx, account.ID, 0, fmt.Errorf("не удалось проверить и обновить poll сессии для повторного выполнения: %v", err))
	}

	err = s.performRetries(ctx, account)
	if err != nil {
		return s.logSyncError(ctx, account.ID, 0, fmt.Errorf("не удалось провести повторный опрос для аккаунта %v: %v", account.ID, err))
	}

	return nil
}

// waitForPollingCompletion ожидает завершения текущих опросов.
func (s *service) waitForPollingCompletion(ctx context.Context, account models.AccountToSync, pointID int) error {
	log.Println("В очереди остались задачи на опрос, ожидаем 60 секунд...")
	isRunning, err := s.lersApiClient.ArePollingsCurrentlyRunning(account.Token, account.ServerHost)
	if err != nil {
		return s.logSyncError(ctx, account.ID, pointID, fmt.Errorf("не удалось получить текущую очередь опроса: %v", err))
	}

	if isRunning {
		time.Sleep(time.Second * 60)
		s.waitForPollingCompletion(ctx, account, pointID)
	}
	return nil
}

// updateFailedPollSessions обновляет статус неудачных сессий опросов.
func (s *service) updateFailedPollSessions(ctx context.Context, pollSessions map[int]string, pollingsToRetry []measure_points.PollSessionsToRetry) error {
	for _, pollToRetry := range pollingsToRetry {
		status, ok := pollSessions[pollToRetry.PollID]
		if !ok {
			err := s.measurePointsRepository.UpdatePollStatus(ctx, pollToRetry.PollID, status)
			if err != nil {
				return err
			}
			continue
		}
		if status == "None" {
			err := s.measurePointsRepository.UpdatePollStatus(ctx, pollToRetry.PollID, "SUCCESS")
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// updateFailedRetrySessions обновляет статус неудачных повторных сессий опросов.
func (s *service) updateFailedRetrySessions(ctx context.Context, pollSessions map[int]string, pollingsToRetry []measure_points.RetryPollSessions) error {
	for _, pollToRetry := range pollingsToRetry {
		status, ok := pollSessions[pollToRetry.RetryPollID]
		if !ok {
			err := s.measurePointsRepository.UpdateRetryPollStatus(ctx, pollToRetry, "FAILED")
			if err != nil {
				return err
			}
			continue
		}
		if status == "None" {
			err := s.measurePointsRepository.UpdateRetryPollStatus(ctx, pollToRetry, "SUCCESS")
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// performRetries выполняет повторные опросы для сессий, которые ранее не были успешными.
func (s *service) performRetries(ctx context.Context, account models.AccountToSync) error {
	// Получить список poll сессий без 'SUCCESS', которые еще можно повторить сегодня.
	toRetryPolls, err := s.measurePointsRepository.FindPollSessionsToRetry2(ctx, account.ID)
	if err != nil {
		return s.logSyncError(ctx, account.ID, 0, fmt.Errorf("не удалось получить опросы по FindPollSessionsToRetry2: %v", err))
	}
	if len(toRetryPolls) == 0 {
		fmt.Println("Нет сессий для повторного выполнения")
		return nil
	}

	// Отправляем запрос на получение данных для таких сессий (retry).
	for _, poll := range toRetryPolls {
		startDate, err := s.getPollingStartDate(ctx, poll.MeasurePointID)
		if err != nil {
			return s.logSyncError(ctx, account.ID, poll.MeasurePointID, err)
		}

		endDate := time.Now().Format(time.RFC3339)
		resp, err := s.lersApiClient.PollMeasurePoints(account.Token, account.ServerHost, []int{poll.MeasurePointID}, startDate, endDate, pollMeasurePointsTimeout)
		if err != nil {
			return s.logSyncError(ctx, account.ID, poll.PollID, fmt.Errorf("не удалось сделать опрос для точки измерения %d: %w", poll.PollID, err))
		}
		err = s.measurePointsRepository.InsertMeasurePointPollRetry(ctx, poll.PollID, resp.PollSessionID, resp.Status)
		if err != nil {
			return s.logSyncError(ctx, account.ID, poll.MeasurePointID, fmt.Errorf("не удалось вставить лог retry опроса %d: %w", poll.MeasurePointID, err))
		}
	}
	// Ожидание завершения всех опросов.
	time.Sleep(time.Second * 30)
	// Проверка, чтобы очередь была пустой.
	isRunning, err := s.lersApiClient.ArePollingsCurrentlyRunning(account.Token, account.ServerHost)
	if err != nil {
		return s.logSyncError(ctx, account.ID, 0, fmt.Errorf("не удалось получить текущую очередь опроса: %v", err))
	}

	if isRunning {
		s.waitForPollingCompletion(ctx, account, 0)
	}
	// Обновление статусов повторных опросов.
	endDate := time.Now().AddDate(0, 0, 1).Format("2006-01-02")
	startDate := time.Now().AddDate(0, 0, -7).Format("2006-01-02")
	pollSessions, err := s.lersApiClient.GetPollSessions(account.Token, account.ServerHost, startDate, endDate, getPollSessionsTimeout)
	if err != nil {
		return s.logSyncError(ctx, account.ID, 0, fmt.Errorf("не удалось получить сессии опроса: %v", err))
	}
	retryPollSession, err := s.measurePointsRepository.FindRetryPollSessions(ctx, account.ID)
	if err != nil {
		return s.logSyncError(ctx, account.ID, 0, fmt.Errorf("не удалось получить poll сессии для повторного выполнения: %v", err))
	}
	err = s.updateFailedRetrySessions(ctx, pollSessions, retryPollSession)
	if err != nil {
		return s.logSyncError(ctx, account.ID, 0, fmt.Errorf("не удалось проверить poll сессии для retry: %v", err))
	}

	// Повторяем процесс ретрая, пока все сессии не будут завершены.
	s.performRetries(ctx, account)
	return nil
}

// getPollingStartDate возвращает стартовую дату для опроса точки измерения.
func (s *service) getPollingStartDate(ctx context.Context, measurePointID int) (string, error) {
	lastDatetime, err := s.measurePointsRepository.GetLastMeasurePointDatetime(ctx, measurePointID)
	if err != nil {
		return "", fmt.Errorf("не удалось получить последний datetime для точки измерения %d: %w", measurePointID, err)
	}

	if lastDatetime == "" {
		// Используем очень раннюю дату, если записей не существует.
		return "2000-01-01T00:00:00Z", nil
	}

	lastTime, err := time.Parse(time.RFC3339, lastDatetime)
	if err != nil {
		return "", fmt.Errorf("не удалось распарсить последний datetime для точки измерения %d: %w", measurePointID, err)
	}

	// Добавляем один час к последнему времени, чтобы определить стартовую дату нового опроса.
	return lastTime.Add(time.Hour).Format(time.RFC3339), nil
}

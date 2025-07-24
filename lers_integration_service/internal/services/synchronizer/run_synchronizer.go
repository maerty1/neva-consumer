package synchronizer

import (
	"context"
	"encoding/json"
	"fmt"
	"lers_integration_service/internal/api_clients/lers"
	"lers_integration_service/internal/models"
	"lers_integration_service/internal/repositories/measure_points"
	"time"
)

func (s *service) RunSynchronizer(ctx context.Context) error {
	accountsToSync, err := s.measurePointsRepository.FindAccountsToSync(ctx)
	if err != nil {
		return err
	}

	if len(accountsToSync) == 0 {
		return fmt.Errorf("нет аккаунтов для синхронизации")
	}

	for _, account := range accountsToSync {
		fmt.Printf("Синхронизация аккаунта %d\n", account.ID)
		if err := s.syncMeasurePointsDataForAccount(ctx, account); err != nil {
			fmt.Printf("Ошибка синхронизации аккаунта %d: %v\n", account.ID, err)
		}
	}

	return nil
}

func (s *service) syncMeasurePointsDataForAccount(ctx context.Context, account models.AccountToSync) error {
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
	fmt.Printf("Получение данных для точки измерения %d с %s по %s\n", point.ID, startDate, endDate)

	consumptionResponse, err := s.lersApiClient.GetConsumptionData(accountID, account.Token, account.ServerHost, point.ID, startDate, endDate)
	if err != nil {
		return s.logSyncError(ctx, accountID, point.ID, fmt.Errorf("не удалось получить данные о потреблении для точки измерения %d: %w", point.ID, err))
	}

	s.storeHourConsumptionData(ctx, accountID, point.ID, consumptionResponse.HourConsumption)
	s.storeDayConsumptionData(ctx, accountID, point.ID, consumptionResponse.DayConsumption)
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

func (s *service) storeHourConsumptionData(ctx context.Context, accountID, measurePointID int, consumptionData []lers.ConsumptionData) error {
	// TODO: Сделать оптимизированную вставку
	for _, data := range consumptionData {
		if len(data.Values) == 0 {
			continue
		}

		jsonValues, err := json.Marshal(data.Values)
		if err != nil {
			s.logSyncError(ctx, accountID, measurePointID, fmt.Errorf("ошибка маршалинга values в JSON для точки измерения %d: %w", measurePointID, err))
			continue
		}

		if err := s.measurePointsRepository.InsertMeasurePointData(ctx, measurePointID, data.DateTime, string(jsonValues)); err != nil {
			s.logSyncError(ctx, accountID, measurePointID, fmt.Errorf("не удалось вставить данные о потреблении для точки измерения %d: %w", measurePointID, err))
			continue
		}
	}

	return nil
}

func (s *service) storeDayConsumptionData(ctx context.Context, accountID, measurePointID int, consumptionData []lers.ConsumptionData) error {
	// TODO: Сделать оптимизированную вставку
	// for _, data := range consumptionData {
	// 	if len(data.Values) == 0 {
	// 		continue
	// 	}

	// 	jsonValues, err := json.Marshal(data.Values)
	// 	if err != nil {
	// 		s.logSyncError(ctx, accountID, measurePointID, fmt.Errorf("ошибка маршалинга values в JSON для точки измерения %d: %w", measurePointID, err))
	// 		continue
	// 	}

	// 	if err := s.measurePointsRepository.InsertMeasurePointDayData(ctx, measurePointID, data.DateTime, string(jsonValues)); err != nil {
	// 		s.logSyncError(ctx, accountID, measurePointID, fmt.Errorf("не удалось вставить данные о потреблении для точки измерения %d: %w", measurePointID, err))
	// 		continue
	// 	}
	// }

	var batchData []measure_points.MeasurePointsDataDay

	for _, data := range consumptionData {
		if len(data.Values) == 0 {
			continue
		}

		jsonValues, err := json.Marshal(data.Values)
		if err != nil {
			s.logSyncError(ctx, accountID, measurePointID, fmt.Errorf("ошибка маршалинга values в JSON для точки измерения %d: %w", measurePointID, err))
			continue
		}

		batchData = append(batchData, measure_points.MeasurePointsDataDay{
			MeasurePointID: measurePointID,
			DateTime:       data.DateTime,
			Values:         string(jsonValues),
		})
	}

	if len(batchData) == 0 {
		return nil
	}

	if err := s.measurePointsRepository.InsertMeasurePointDayDataBatch(ctx, batchData); err != nil {
		s.logSyncError(ctx, accountID, measurePointID, fmt.Errorf("не удалось выполнить пакетную вставку данных о потреблении для точки измерения %d: %w", measurePointID, err))
		return err
	}

	return nil
}

package interrogator

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
	"weather_station_data_collector/internal/models"
)

func (s *service) RunInterrogator(ctx context.Context) {
	interrogationChannel := make(chan models.WeatherStationResponseChannel)
	defer close(interrogationChannel)

	go s.checkFailedInterrogationsTicker(ctx)

	for {
		go func() {
			s.Interrogation(ctx, interrogationChannel)
		}()
		response := <-interrogationChannel

		if response.Err != nil {
			s.logInterrogateError(ctx, time.Now().Truncate(time.Minute).String(), response.Err)
		}

		if response.Code != http.StatusOK || response.Err != nil || response.Body.Dateutc == "" {
			s.logInterrogateWarning(ctx, time.Now().Truncate(time.Minute).String(),
				errors.New(fmt.Sprintf("Код ответа: %d, Ошибка: %s запрашиваемое время будет опрошено позже",
					response.Code, response.Err.Error())))

			err := s.RawWeatherDataRepository.InsertRawWeatherData(ctx, models.WeatherDataRaw{
				Dateutc: time.Now().Format("2006-01-02 15:04:05")})
			if err != nil {
				s.logInterrogateError(ctx, time.Now().Truncate(time.Minute).String(),
					errors.New(fmt.Sprintf("не получилось записать в бд, ошибка: %s", err.Error())))
			}
		}

		if response.Code == http.StatusOK && response.Err == nil {
			go func() {
				s.logInterrogateInfo(ctx, response.Body.Dateutc, "записть в бд")
				err := s.RawWeatherDataRepository.InsertRawWeatherData(ctx, response.Body)
				if err != nil {
					s.logInterrogateError(ctx, response.Body.Dateutc, errors.New(
						fmt.Sprintf("не получилось записать в бд, ошибка: %s", err.Error())))
				} else {
					s.logInterrogateSuccess(ctx, response.Body.Dateutc)
				}
			}()
		}

		time.Sleep(time.Second * s.checkGapDuration)
	}
}

func (s *service) checkFailedInterrogationsTicker(ctx context.Context) {
	startTime := time.Now()
	nextCheckTime := startTime.Add(time.Second * s.checkFailedTimeGapDuration)

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			now := time.Now()
			if now.After(nextCheckTime) {
				failedTimes, err := s.RawWeatherDataRepository.SelectTimeWithNullData(ctx)
				if err != nil {
					s.logInterrogateError(ctx, "", errors.New(
						fmt.Sprintf("ошибка при получении списка времён без данных. Ошибка: %s", err.Error())))
				}
				if len(failedTimes) > 0 {
					go s.handleFailedInterrogations(ctx, failedTimes)
				}
				nextCheckTime = nextCheckTime.Add(time.Second * s.checkFailedTimeGapDuration)
			}
		case <-ctx.Done():
			return
		}
	}
}

func (s *service) handleFailedInterrogations(ctx context.Context, failedInterrogationsTime []time.Time) {
	for _, failedTime := range failedInterrogationsTime {
		s.logInterrogateInfo(ctx, failedTime.String(), "проверка времени, которое не получилось проверить ранее")
		fromTime := failedTime.Format("2006-01-02 15:04")
		toTime := failedTime.Add(time.Minute).Format("2006-01-02 15:04")

		fromTime = strings.ReplaceAll(fromTime, " ", "%20")
		toTime = strings.ReplaceAll(toTime, " ", "%20")

		res, err := s.ApiClient.GetHistoryMinute(ctx, fromTime, toTime)
		if err != nil {
			return
		}

		var closestElement models.WeatherDataRaw
		var closestTimeDiff time.Duration

		closestTimeDiff = time.Hour * 24

		for _, v := range res {
			parsedTime, err := time.Parse("2006-01-02 15:04:05", v.Dateutc)
			if err != nil {

			}
			timeDiff := failedTime.Sub(parsedTime)

			if timeDiff < 0 {
				timeDiff = -timeDiff
			}

			if timeDiff < closestTimeDiff {
				closestTimeDiff = timeDiff
				closestElement = v
			}
		}

		if closestElement.Dateutc != "" {
			go func() {
				err := s.RawWeatherDataRepository.UpdateRawWeatherData(ctx, closestElement, failedTime)
				if err != nil {
					s.logInterrogateError(ctx, closestElement.Dateutc, errors.New(
						fmt.Sprintf("не получилось записать в бд, ошибка: %s", err.Error())))
				}
			}()
		}
	}
}

func (s *service) Interrogation(ctx context.Context, ch chan models.WeatherStationResponseChannel) {
	s.logInterrogateInfo(ctx, time.Now().Truncate(time.Minute).String(), "Проверка по времени")
	res := s.ApiClient.GetLastData(ctx)
	ch <- res
}

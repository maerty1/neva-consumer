package interrogator

import ( 
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
	"weather_station_data_collector/internal/models"
)

func (s *service) RunHandleInterrogator(ctx context.Context, from, to string) {
	timeFromParsed, err := time.Parse("2006-01-02 15:04", from)
	if err != nil {
		s.logInterrogateError(ctx, from, errors.New(
			fmt.Sprintf("неверный формат времени: %s", err.Error())))
	}
	timeToParsed, err := time.Parse("2006-01-02 15:04", to)
	if err != nil {
		s.logInterrogateError(ctx, from, errors.New(
			fmt.Sprintf("неверный формат времени: %s", err.Error())))
	}

	from = timeFromParsed.Format("2006-01-02 15:04")
	to = timeToParsed.Format("2006-01-02 15:04")

	from = strings.ReplaceAll(from, " ", "%20")
	to = strings.ReplaceAll(to, " ", "%20")

	res, err := s.ApiClient.GetHistoryMinute(ctx, from, to)
	if err != nil {
		s.logInterrogateError(ctx, from, errors.New(
			fmt.Sprintf("ошибка при получении времени: %s", err.Error())))
		return
	}

	if len(res) == 0 {
		s.logInterrogateError(ctx, from, errors.New(
			fmt.Sprintf("ошибка при получении времени: %s", errors.New("история пуста"))))
		return
	}

	filteredRes := []models.WeatherDataRaw{res[0]}
	for i := 1; i < len(res); i++ {
		currentTime, err := time.Parse("2006-01-02 15:04:05", res[i].Dateutc)
		if err != nil {
			s.logInterrogateError(ctx, from, errors.New(
				fmt.Sprintf("неверный формат времени: %s", err.Error())))
			continue
		}

		lastTime, err := time.Parse("2006-01-02 15:04:05", filteredRes[len(filteredRes)-1].Dateutc)
		if err != nil {
			s.logInterrogateError(ctx, from, errors.New(
				fmt.Sprintf("неверный формат времени: %s", err.Error())))
			continue
		}

		if currentTime.Sub(lastTime) >= s.checkGapDuration*time.Second {
			filteredRes = append(filteredRes, res[i])
		}
	}

	err = s.RawWeatherDataRepository.InsertWeatherDataBatch(ctx, filteredRes)
	if err != nil {
		s.logInterrogateError(ctx, "", errors.New(
			fmt.Sprintf("не получилось записать в бд, ошибка: %s", err.Error())))
	}

	s.logInterrogateSuccess(ctx, "")
}

package interrogator

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"
	weatherStation "weather_station_data_collector/internal/api_client/weather_station"
	weatherDataRaw "weather_station_data_collector/internal/repositories/weather_data_raw"
)

type Service interface {
	RunInterrogator(ctx context.Context)
	RunHandleInterrogator(ctx context.Context, timeFrom string, timeTo string)
}

var _ Service = (*service)(nil)

type service struct {
	RawWeatherDataRepository weatherDataRaw.Repository
	ApiClient                weatherStation.ApiClient

	checkGapDuration           time.Duration
	checkFailedTimeGapDuration time.Duration
}

func NewService(repository weatherDataRaw.Repository, apiClient weatherStation.ApiClient) Service {
	var res service
	gap := os.Getenv("CHECK_GAP_SECONDS")
	gapInt, err := strconv.Atoi(gap)
	if err != nil {
		log.Fatal("не получилось распарсить CHECK_GAP_SECONDS")
	}
	res.checkGapDuration = time.Duration(gapInt)
	if os.Getenv("TIME_FROM") == "" {
		gapFailedTime := os.Getenv("CHECK_FAILED_TIME_GAP_SECONDS")
		gapFailedTimeInt, err := strconv.Atoi(gapFailedTime)
		if err != nil {
			log.Fatal("не получилось распарсить CHECK_FAILED_TIME_GAP_SECONDS")
		}
		res.checkFailedTimeGapDuration = time.Duration(gapFailedTimeInt)
	}
	res.RawWeatherDataRepository = repository
	res.ApiClient = apiClient
	return &res
}

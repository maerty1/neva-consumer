package zulu_updater

import (
	"context"
	"log"
	"os"
	"strconv"
	"zulu_updater/internal/api_client/weather"
	zuluApi "zulu_updater/internal/api_client/zulu"
	measurePointsDataDay "zulu_updater/repositories/measure_points_data_day"
	zuluRepo "zulu_updater/repositories/zulu"
)

type Service interface {
	RunZuluUpdater(ctx context.Context)
}

var _ Service = (*service)(nil)

type service struct {
	rootRepository   measurePointsDataDay.Repository
	zuluRepository   zuluRepo.Repository
	zuluApiClient    zuluApi.ApiClient
	weatherApiClient weather.ApiClient

	elemId int
}

func NewService(
	rootRepository measurePointsDataDay.Repository,
	zuluApiClient zuluApi.ApiClient,
	weatherApiClient weather.ApiClient,
	zuluRepository zuluRepo.Repository) Service {

	var res service
	var err error

	res.elemId, err = strconv.Atoi(os.Getenv("ELEM_ID"))
	if err != nil {
		log.Fatal("невозможно преобразовать ELEM_ID int")
	}
	res.zuluApiClient = zuluApiClient
	res.rootRepository = rootRepository
	res.weatherApiClient = weatherApiClient
	res.zuluRepository = zuluRepository

	return &res
}

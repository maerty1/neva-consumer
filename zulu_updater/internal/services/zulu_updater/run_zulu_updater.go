package zulu_updater

import (
	"context"
	"errors"
	"fmt"
	"time"
)

const InternalTemperature = "T_in"
const ExternalTemperature = "T1_t"
const ZuluExternalTemperature = "Tnv_t"
const SleepDuration = time.Duration(30)
const UpdateType = "calculated"
const ExtractionType = "calculated_daily"

func (s service) RunZuluUpdater(ctx context.Context) {
	lersTs := time.Now().Truncate(time.Hour * 24)
	fmt.Println("Текущая дата: ", lersTs)
	tIn, err := s.rootRepository.GetDataParameterByDay(ctx, lersTs, InternalTemperature)
	if err != nil {
		s.logUpdaterError(ctx, err)
		return
	}
	s.logUpdaterSuccess(ctx, "Параметры получены")

	weatherData, err := s.weatherApiClient.GetWeatherData()
	if err != nil {
		s.logUpdaterError(ctx, err)
		return
	}
	s.logUpdaterSuccess(ctx, "Данные о погоде получены")

	err = s.zuluApiClient.UpdateAttribute(ctx, ExternalTemperature, tIn, s.elemId)
	if err != nil {
		s.logUpdaterError(ctx, err)
		return
	}
	s.logUpdaterSuccess(ctx, "Атрибут T1_t изменён")

	err = s.zuluApiClient.UpdateAttribute(ctx, ZuluExternalTemperature, weatherData.Today.Temp, s.elemId)
	if err != nil {
		s.logUpdaterError(ctx, err)
		return
	}
	s.logUpdaterSuccess(ctx, "Атрибут Tnv_t изменён")

	taskHandle, err := s.zuluApiClient.RunZuluCalculation(ctx)
	if err != nil {
		s.logUpdaterError(ctx, err)
		return
	}
	s.logUpdaterInfo(ctx, "Ожидание расчёта")
	time.Sleep(time.Second * SleepDuration)

	status, err := s.zuluApiClient.GetZuluTaskStatus(ctx, taskHandle)

	if status != "finished" {
		s.logUpdaterWarning(ctx, errors.New("статус расчёта: "+status))
		return
	}
	s.logUpdaterSuccess(ctx, "Расчёт завершён")

	errCount, errs, err := s.zuluApiClient.GetZuluErrors(ctx, taskHandle)
	if err != nil {
		s.logUpdaterError(ctx, err)
		return
	}
	if errCount > 0 {
		message := fmt.Sprintf("расчёт окончен с ошибками в количестве %d\nОшибки:", errCount)
		for _, v := range errs {
			message = message + v.Name + ": " + v.Text
		}
		s.logUpdaterWarning(ctx, errors.New(message))
	}
	s.logUpdaterSuccess(ctx, "Расчёт успешно завершён")

	zwsType, err := s.zuluRepository.SelectZwsTypeByElemId(ctx, s.elemId)
	if err != nil {
		s.logUpdaterError(ctx, err)
		return
	}

	valNames, err := s.zuluRepository.SelectValNameByZwsType(ctx, zwsType, ExtractionType)
	if err != nil {
		s.logUpdaterError(ctx, err)
		return
	}

	result, err := s.zuluApiClient.ExecuteSqlGetParametersByValNames(ctx, valNames, s.elemId)
	if err != nil {
		s.logUpdaterError(ctx, err)
		return
	}
	s.logUpdaterSuccess(ctx, "Данные расчёта получены")

	err = s.zuluRepository.InsertRecords(ctx, result.Record[0].Field, s.elemId)
	if err != nil {
		s.logUpdaterError(ctx, err)
		return
	}
	s.logUpdaterSuccess(ctx, "Данные расчёта сохранены в бд")

	res, err := s.zuluRepository.InsertObjectRecordsJson(ctx, *result, UpdateType, s.elemId, lersTs)
	if err != nil {
		s.logUpdaterError(ctx, err)
		return
	}
	if res.ErrMsg != "" {
		s.logUpdaterWarning(ctx, errors.New(res.ErrMsg))
		return
	}
	s.logUpdaterSuccess(ctx, "Данные расчёта сохранены в бд")
}

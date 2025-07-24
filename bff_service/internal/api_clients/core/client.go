package core

import "bff_service/internal/config"

type ApiClient interface {
	GetPointsData(reqData []GetPointsDataRequest, timestamp string) ([]GetPointsDataResponse, error)
	GetPointsDataHistory(reqData []GetPointsDataHistoryRequest, nDays int, timestamp string) (GetPointsDataHistoryResponse, error)
	GetElementIDs() ([]int, error)
}

var _ ApiClient = (*apiClient)(nil)

type apiClient struct {
	serviceMapper config.ServiceMapper
}

func NewApiClient(serviceMapper config.ServiceMapper) ApiClient {
	return &apiClient{serviceMapper: serviceMapper}
}

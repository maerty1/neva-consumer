package zulu

import "bff_service/internal/config"

type ApiClient interface {
	GetPoints(zwsTypeIds []int) ([]Point, error)
	GetFullPoint(elemID int, nDays int) (FullElementData, error)
	GetPointCategoryDataGroup(elemID int, categoryID int, timestamp string) (GetPointDataByCategoryGroup, error)
	GetPointCategoryDataKeyvalue(elemID int, categoryID int) (GetPointDataByCategoryKeyvalue, error)
	GetFilteredPoints(elementIds []int, zwsTypeIDs []int, timestamp string) ([]Point, error)
	GetPointsDataCategoryByZwsType(zwsTypeIds []int, categoryID int, timestamp string) (GetPointsDataCategoryResponse, error)
}

var _ ApiClient = (*apiClient)(nil)

type apiClient struct {
	serviceMapper config.ServiceMapper
}

func NewApiClient(serviceMapper config.ServiceMapper) ApiClient {
	return &apiClient{serviceMapper: serviceMapper}
}

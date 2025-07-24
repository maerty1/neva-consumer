package geodata

import (
	"context"

	"zulu_service/internal/db"
	"zulu_service/internal/models/geodata"
)

type Repository interface {
	GetGeoJson(ctx context.Context) ([]byte, error)
	GetGeoJsonV2(ctx context.Context) ([]byte, error)
	GetStates(ctx context.Context) ([]geodata.ObjectState, error)
	GetElementDataByID(ctx context.Context, elementID int) ([]geodata.ElementData, error)
	GetPipelineDepth(ctx context.Context) (map[int]geodata.PipelineDepth, error)

	GetPoints(ctx context.Context, zwsTypeIDs []int) ([]geodata.Point, error)
	GetFilteredPoints(ctx context.Context, elementIDs []int, zwsTypeIDs []int, timestamp string) ([]geodata.Point, error)
	GetMeasurementGroupsEnum(ctx context.Context) (map[int]geodata.MeasurementGroupEnum, error)
	GetFullByElemID(ctx context.Context, elemID int, nDays int) (*geodata.FullElementData, error)
	GetPointCategories(ctx context.Context, elemID int) (geodata.PointWithCategories, error)
	GetPointDataByCategoryGroup(ctx context.Context, elemID int, categoryID int, timestamp string, nDays int) (*geodata.GetPointDataByCategoryGroup, error)
	GetPointDataByCategoryKeyvalue(ctx context.Context, elemID int, categoryID int) (geodata.GetPointDataByCategoryKeyvalue, error)
	GetIconIdByElemId(ctx context.Context, elemID int) (int, error)
	GetIconIdByZwsType(ctx context.Context, elemID int) (int, error)
	GetSchemaIdByElemId(ctx context.Context, elemID int) (int, error)

	GetPointsDataByCategoryGroup(ctx context.Context, elemIDs []int, categoryID int) (map[int]*geodata.GetPointsDataByCategoryGroup, error)
	GetPointsDataByZwsTypes(ctx context.Context, zwsTypeIDs []int, categoryID int) (map[int]*geodata.GetPointsDataByCategoryGroup, error)
	GetPointsDataByZwsTypesV2(ctx context.Context, zwsTypeIDs []int, categoryID int, timestamp string) (map[int]*geodata.GetPointsDataByCategoryGroup, error)
}

var _ Repository = (*repository)(nil)

type repository struct {
	db db.PostgresClient
}

func NewRepository(db db.PostgresClient) *repository {
	return &repository{
		db: db,
	}
}

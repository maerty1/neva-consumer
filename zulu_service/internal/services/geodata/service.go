package geodata

import (
	geodataRepository "zulu_service/internal/repositories/geodata"
)

type Service interface {
	// GetGeoJson(ctx context.Context) ([]byte, error)
	// GetGeoJsonV2(ctx context.Context) ([]byte, error)
	// GetStates(ctx context.Context) ([]geodata.ObjectState, error)
	// GetElementDataByID(ctx context.Context, elementID int) ([]geodata.ElementData, error)
}

var _ Service = (*service)(nil)

type service struct {
	geodataRepository geodataRepository.Repository
}

func NewService(geodataRepository geodataRepository.Repository) Service {
	return &service{
		geodataRepository: geodataRepository,
	}
}

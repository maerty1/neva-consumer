package geodata

import (
	"context"
	"zulu_service/internal/models/geodata"
)

func (r *service) GetElementDataByID(ctx context.Context, elementID int) ([]geodata.ElementData, error) {

	elementData, err := r.geodataRepository.GetElementDataByID(ctx, elementID)
	if err != nil {
		return nil, err
	}

	return elementData, nil

}

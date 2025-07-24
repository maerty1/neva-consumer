package geodata

import (
	"context"
)

func (r *service) GetGeoJsonV2(ctx context.Context) ([]byte, error) {

	geojson, err := r.geodataRepository.GetGeoJsonV2(ctx)
	if err != nil {
		return nil, err
	}

	return geojson, nil

}

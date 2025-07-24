package geodata

import (
	"context"
)

func (r *service) GetGeoJson(ctx context.Context) ([]byte, error) {

	geojson, err := r.geodataRepository.GetGeoJson(ctx)
	if err != nil {
		return nil, err
	}

	return geojson, nil

}

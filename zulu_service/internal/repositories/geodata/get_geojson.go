package geodata

import (
	"context"
	"log"
)

func (r *repository) GetGeoJson(ctx context.Context) ([]byte, error) {
	var geoJSON []byte
	query := `
        SELECT jsonb_build_object(
            'type',     'FeatureCollection',
            'features', jsonb_agg(features.feature)
        )
        FROM (
            SELECT jsonb_build_object(
                'type',       'Feature',
                'id',         elem_id,
                'geometry',   ST_AsGeoJSON(ST_FlipCoordinates(zws_geometry))::jsonb,
                'properties', to_jsonb(inputs) - 'zws_geometry'
            ) AS feature
            FROM (
                SELECT elem_id, parent_id, zws_mode, zws_type, zws_linecolor, zws_geometry 
                FROM zulu.objects_geometry_log
            ) inputs
        ) features;
    `

	err := r.db.DB().QueryRow(ctx, query).Scan(&geoJSON)
	if err != nil {
		log.Printf("Ошибка при получении GeoJSON: %v", err)
		return nil, err
	}

	return geoJSON, nil
}

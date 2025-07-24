package geodata

import (
	"context"
	"log"
)

func (r *repository) GetGeoJsonV2(ctx context.Context) ([]byte, error) {
	var geoJSON []byte
	query := `
        SELECT jsonb_build_object(
            'type',     'FeatureCollection',
            'features', jsonb_agg(features.feature)
        )
        FROM (
            SELECT jsonb_build_object(
                'type',       'Feature',
                'id',         inputs.elem_id,
                'geometry',   ST_AsGeoJSON(ST_FlipCoordinates(inputs.zws_geometry))::jsonb,
                'properties', to_jsonb(inputs) - 'zws_geometry' || jsonb_build_object(
                    'Name', name_record.val,
                    'Adres', adres_record.val
                )
            ) AS feature
            FROM (
                SELECT elem_id, parent_id, zws_mode, zws_type, zws_linecolor, zws_geometry 
                FROM zulu.objects_geometry_log
            ) inputs
            LEFT JOIN LATERAL (
                SELECT val
                FROM zulu.object_records
                WHERE parameter = 'Name'
                  AND elem_id = inputs.elem_id
                  AND is_deleted = false
                  AND td IS NULL
                ORDER BY inserted_ts DESC
                LIMIT 1
            ) name_record ON TRUE
            LEFT JOIN LATERAL (
                SELECT val
                FROM zulu.object_records
                WHERE parameter = 'Adres'
                  AND elem_id = inputs.elem_id
                  AND is_deleted = false
                  AND td IS NULL
                ORDER BY inserted_ts DESC
                LIMIT 1
            ) adres_record ON TRUE
        ) features;
    `

	err := r.db.DB().QueryRow(ctx, query).Scan(&geoJSON)
	if err != nil {
		log.Printf("Ошибка при получении GeoJSON: %v", err)
		return nil, err
	}

	return geoJSON, nil
}

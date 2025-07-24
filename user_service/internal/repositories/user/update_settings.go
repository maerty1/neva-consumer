package user

import (
	"context"
	"log"

	"user_service/internal/config/errors"
	"user_service/internal/models"
)

func (r *repository) UpdateSettings(ctx context.Context, clientID int, settings models.UserSettingsUpdate) error {

	_, err := r.db.DB().Exec(
		ctx,
		`INSERT INTO user_settings (
			user_id, 
			bounding_box_top_left_lat, 
			bounding_box_top_left_lon, 
			bounding_box_bottom_right_lat, 
			bounding_box_bottom_right_lon, 
			default_zoom, 
			default_center_lat, 
			default_center_lon
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		ON CONFLICT (user_id) 
		DO UPDATE SET 
			bounding_box_top_left_lat = EXCLUDED.bounding_box_top_left_lat, 
			bounding_box_top_left_lon = EXCLUDED.bounding_box_top_left_lon, 
			bounding_box_bottom_right_lat = EXCLUDED.bounding_box_bottom_right_lat, 
			bounding_box_bottom_right_lon = EXCLUDED.bounding_box_bottom_right_lon, 
			default_zoom = EXCLUDED.default_zoom, 
			default_center_lat = EXCLUDED.default_center_lat, 
			default_center_lon = EXCLUDED.default_center_lon`,
		clientID,
		settings.BoundingBoxTopLeftLat,
		settings.BoundingBoxTopLeftLon,
		settings.BoundingBoxTopRightLat,
		settings.BoundingBoxTopRightLon,
		settings.DefaultZoom,
		settings.DefaultCenterLat,
		settings.DefaultCenterLon,
	)

	if err != nil {
		log.Println(err.Error())
		return errors.ErrInternalError
	}

	return nil
}

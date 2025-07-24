package user

import (
	"context"

	"user_service/internal/config/errors"
	"user_service/internal/models"

	"github.com/jackc/pgx/v4"
)

func (r *repository) GetSettings(ctx context.Context, clientID int) (models.UserSettingsResponse, error) {
	var userSettings models.UserSettingsResponse

	row := r.db.DB().QueryRow(
		ctx,
		`SELECT bounding_box_top_left_lat, bounding_box_top_left_lon, bounding_box_bottom_right_lat, bounding_box_bottom_right_lon, default_zoom, default_center_lat, default_center_lon
		FROM user_settings
		WHERE user_id = $1
		`,
		clientID,
	)

	err := row.Scan(&userSettings.BoundingBoxTopLeftLat, &userSettings.BoundingBoxTopLeftLon, &userSettings.BoundingBoxTopRightLat, &userSettings.BoundingBoxTopRightLon, &userSettings.DefaultZoom, &userSettings.DefaultCenterLat, &userSettings.DefaultCenterLon)
	if err != nil {
		if err == pgx.ErrNoRows {
			return models.UserSettingsResponse{}, nil
		}
		return models.UserSettingsResponse{}, errors.ErrInternalError
	}

	return userSettings, nil
}

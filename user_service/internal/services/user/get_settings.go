package user

import (
	"context"
	"user_service/internal/models"
)

func (r *service) GetSettings(ctx context.Context, clientID int) (models.UserSettingsResponse, error) {

	userSettingsResponse, err := r.userRepository.GetSettings(ctx, clientID)
	if err != nil {
		return models.UserSettingsResponse{}, err
	}

	return userSettingsResponse, nil

}

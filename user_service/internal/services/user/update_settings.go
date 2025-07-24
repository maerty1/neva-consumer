package user

import (
	"context"
	"user_service/internal/models"
)

func (r *service) UpdateSettings(ctx context.Context, clientID int, settings models.UserSettingsUpdate) error {

	err := r.userRepository.UpdateSettings(ctx, clientID, settings)
	if err != nil {
		return err
	}

	return nil

}

package user

import (
	"context"
	"user_service/internal/models"
	userRepository "user_service/internal/repositories/user"
)

type Service interface {
	CreateUser(ctx context.Context, user models.UserRegisterRequest) error
	AuthenticateUser(ctx context.Context, user models.UserAuthenticateRequest) (*models.UserAuthenticateResponse, error)
	GetSettings(ctx context.Context, clientID int) (models.UserSettingsResponse, error)
	UpdateSettings(ctx context.Context, clientID int, settings models.UserSettingsUpdate) error
}

var _ Service = (*service)(nil)

type service struct {
	userRepository userRepository.Repository
}

func NewService(userRepository userRepository.Repository) Service {
	return &service{
		userRepository: userRepository,
	}
}

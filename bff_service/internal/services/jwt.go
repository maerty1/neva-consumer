package services

import (
	"bff_service/internal/config"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

type JwtService interface {
	GenerateWebToken(userID int) (string, error)
	GenerateAlisaToken(userID int) (string, error)
}

var _ JwtService = (*service)(nil)

type service struct {
	jwtConfig config.JWTConfig
}

func (s *service) GenerateWebToken(userID int) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 72).Unix(), // Токен истекает через 72 часа
		"aud":     string(config.AudienceWebApp),
	}

	// Создаем токен с методом подписи HMAC и claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.jwtConfig.GetJWTSecret()))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
func (s *service) GenerateAlisaToken(userID int) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 240).Unix(), // Токен истекает через 30 дней
		"aud":     string(config.AudienceAliceSkill),
	}

	// Создаем токен с методом подписи HMAC и claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.jwtConfig.GetJWTSecret()))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func NewService(jwtConfig config.JWTConfig) JwtService {
	return &service{
		jwtConfig: jwtConfig,
	}
}

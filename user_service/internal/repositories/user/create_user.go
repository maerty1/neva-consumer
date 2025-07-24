package user

import (
	"context"
	"log"

	"user_service/internal/config/errors"
	"user_service/internal/models"
	"user_service/internal/utils/hasher"

	"github.com/jackc/pgconn"
)

func (r *repository) CreateUser(ctx context.Context, user models.UserRegisterRequest) error {

	hashedPassword, err := hasher.HashPassword(user.Password)
	if err != nil {
		log.Printf("ошибка хеширования пароля: %v", err)
		return errors.ErrInternalError
	}

	row := r.db.DB().QueryRow(
		ctx, `
			INSERT INTO users (email, password_hash)
			VALUES ($1, $2)
			RETURNING id;
			`,
		user.Email, hashedPassword)

	var userID string
	err = row.Scan(&userID)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgErr.Code == "23505" {
				log.Printf("ошибка: дублирующийся ключ для email %v", user.Email)
				return errors.DuplicateWithDetails(
					"Пользователь с таким email уже существует",
				)
			}
		}

		log.Printf("ошибка при добавлении пользователя: %v", err)
		return errors.ErrInternalError
	}

	return nil
}

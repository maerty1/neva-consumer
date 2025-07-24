package models

import (
	"errors"
	"user_service/internal/models/validators"
)

type UserRegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserAuthenticateRequest struct {
	Login    string `json:"email"`
	Password string `json:"password"`
}

type UserAuthenticateResponse struct {
	ID int `json:"id"`
}

func (c *UserRegisterRequest) UserRegisterValidate() error {
	if !validators.IsValidEmail(c.Email) {
		return errors.New("некорректный email")
	}

	if len(c.Password) == 0 || len(c.Password) <= 8 {
		return errors.New("пароль не может быть пустой строкой и длинна должна быть равна/больше 8 символам")
	}

	return nil
}

func (c *UserAuthenticateRequest) UserAuthenticateValidate() error {
	if len(c.Password) == 0 || len(c.Password) <= 8 {
		return errors.New("пароль не может быть пустой строкой и длинна должна быть равна/больше 8 символам")
	}

	return nil
}

type UserSettingsResponse struct {
	BoundingBoxTopLeftLat  float64 `json:"bounding_box_top_left_lat"`
	BoundingBoxTopLeftLon  float64 `json:"bounding_box_top_left_lon"`
	BoundingBoxTopRightLat float64 `json:"bounding_box_bottom_right_lat"`
	BoundingBoxTopRightLon float64 `json:"bounding_box_bottom_right_lon"`

	DefaultZoom      int     `json:"default_zoom"`
	DefaultCenterLat float64 `json:"default_center_lat"`
	DefaultCenterLon float64 `json:"default_center_lon"`
}

type UserSettingsUpdate struct {
	BoundingBoxTopLeftLat  float64 `json:"bounding_box_top_left_lat"`
	BoundingBoxTopLeftLon  float64 `json:"bounding_box_top_left_lon"`
	BoundingBoxTopRightLat float64 `json:"bounding_box_bottom_right_lat"`
	BoundingBoxTopRightLon float64 `json:"bounding_box_bottom_right_lon"`

	DefaultZoom      int     `json:"default_zoom"`
	DefaultCenterLat float64 `json:"default_center_lat"`
	DefaultCenterLon float64 `json:"default_center_lon"`
}

func (c *UserSettingsUpdate) UserSettingsUpdateValidate() error {

	if c.DefaultZoom < 0 || c.DefaultZoom > 16 {
		return errors.New("Default zoom должен быть от 0 до 16")
	}

	if c.DefaultCenterLat < -90 || c.DefaultCenterLat > 90 {
		return errors.New("DefaultCenterLat должна быть в диапазоне от -90 до 90")
	}
	if c.DefaultCenterLon < -180 || c.DefaultCenterLon > 180 {
		return errors.New("DefaultCenterLon должна быть в диапазоне от -180 до 180")
	}

	return nil
}

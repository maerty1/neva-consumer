package middleware

import (
	"bff_service/internal/config"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
)

// Роуты, куда доступен Алиса навык
var ALICE_ONLY_ROUTES = []string{"/core/api/v1/boiler_room_engineer_report/format/json", "/core/api/v2/boiler_room_engineer_report/format/json", "/core/api/v1/weather/with_forecast", "/core/api/v1/status/current"}
var COOKIE_AUTH = []string{}
var NO_HEADER_AUTH_TRIE = []string{"/users/api/v1/authenticate",
	"/bff/api/docs",

	"/users/api/docs",
	"/zulu/api/docs",
	"/core/api/docs",
	"/core/api/openapi.json",
	"/bff/api/microservices"}

func JwtMiddleware(cfg config.JWTConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Пропуск JWT-проверки для некоторых путей
		for _, prefix := range NO_HEADER_AUTH_TRIE {
			if strings.HasPrefix(c.Request.URL.Path, prefix) {
				c.Set("user_id", 0)
				c.Next()
				return
			}
		}

		var tokenString string

		// Проверка на наличие JWT в куки
		for _, path := range COOKIE_AUTH {
			if strings.HasPrefix(c.Request.URL.Path, path) {
				if cookie, err := c.Cookie("jwt"); err == nil {
					tokenString = cookie
				} else {
					c.String(http.StatusUnauthorized, "JWT не указан в файлах cookie")
					c.Abort()
					return
				}
				break
			}
		}

		// Проверка заголовка Authorization
		if tokenString == "" {
			authHeader := c.GetHeader("Authorization")
			if strings.HasPrefix(authHeader, "Bearer ") {
				tokenString = strings.TrimPrefix(authHeader, "Bearer ")
			} else {
				c.String(http.StatusUnauthorized, "Заголовок авторизации не найден или не начинается с Bearer")
				c.Abort()
				return
			}
		}

		// Проверка и разбор JWT токена
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("Некорректный метод подписи")
			}
			return []byte(cfg.GetJWTSecret()), nil
		})

		if err != nil {
			var errMsg string
			if errors.Is(err, jwt.ErrTokenExpired) {
				errMsg = "Просроченная подпись"
			} else if errors.Is(err, jwt.ErrSignatureInvalid) {
				errMsg = "Ошибка проверки подписи"
			} else {
				errMsg = "Недействительный JWT"
			}
			c.String(http.StatusUnauthorized, errMsg)
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			c.String(http.StatusUnauthorized, "Недействительный JWT")
			c.Abort()
			return
		}

		userID, exists := claims["user_id"]
		if !exists {
			c.String(http.StatusUnauthorized, "Отсутствует user_id в токене")
			c.Abort()
			return
		}

		audClaim, exists := claims["aud"]
		if !exists {
			c.String(http.StatusUnauthorized, "Отсутствует aud в токене")
			c.Abort()
			return
		}

		aud, ok := audClaim.(string)
		if !ok {
			c.String(http.StatusUnauthorized, "Некорректный формат aud")
			c.Abort()
			return
		}

		// Проверяем, можно ли Алисе идти в этот эндпоинт
		if aud == string(config.AudienceAliceSkill) {
			allowed := false
			for _, route := range ALICE_ONLY_ROUTES {
				if strings.HasPrefix(c.Request.URL.Path, route) {
					allowed = true
					break
				}
			}
			if !allowed {
				c.String(http.StatusForbidden, "Доступ только для веб-приложения")
				c.Abort()
				return
			}
		}

		// Сохраняем user_id в контекст
		c.Set("user_id", userID)

		c.Next()
	}
}

package users

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type AlisaTokenResponse struct {
	Access string `json:"access" example:"..."`
}

// @Router /users/api/v1/tokens/alisa [get]
// @Summary Получение токена для навыка Алисы.
// @Tags Users
// @Produce  json
// @Param request body UserAuthenticate true "Запрос аутентификации пользователя"
// @Success 200 {object} AlisaTokenResponse "Успешный ответ, содержащий JWT токен в поле 'access'"
func (f *facade) GetAlisaToken(c *gin.Context) {
	userID, exists := c.Get("user_id")

	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user_id не найден"})
	}

	userIDFloat, ok := userID.(float64)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user_id type"})
		return
	}

	userIDInt := int(userIDFloat)

	alisaToken, err := f.jwtService.GenerateAlisaToken(userIDInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, AlisaTokenResponse{
		Access: alisaToken,
	})
}

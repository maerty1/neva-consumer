package v1

import (
	"net/http"
	"strconv"

	"user_service/internal/config/errors"
	"user_service/internal/models"

	"user_service/internal/services/user"

	"github.com/gin-gonic/gin"
)

func RegisterUsersRouter(r *gin.Engine, userService user.Service) {
	user := r.Group("/users/api/v1")
	{
		user.GET("/settings", func(ctx *gin.Context) {
			getSettings(ctx, userService)
		})
		user.PUT("/settings", func(ctx *gin.Context) {
			updateSettings(ctx, userService)
		})

	}

}

// @Router /users/api/v1/settings [get]
// @Summary Получение настроек пользователя
// @Tags Settings
// @Produce  json
// @Success 200 {object} models.UserSettingsResponse "Настройки"
func getSettings(ctx *gin.Context, userService user.Service) {
	clientID := ctx.GetHeader("x-user-id")
	clientIDInt, err := strconv.Atoi(clientID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный client ID"})
		return
	}

	userSettings, err := userService.GetSettings(ctx, clientIDInt)
	if err != nil {
		resp := errors.GetHTTPStatus(err)
		ctx.JSON(resp.Status, gin.H{"error": resp.Message})
		return
	}
	ctx.JSON(http.StatusOK, userSettings)

}

// @Router /users/api/v1/settings [put]
// @Summary Обновление настроек пользователя
// @Tags Settings
// @Accept  json
// @Produce  json
// @Param request body models.UserSettingsUpdate true "Запрос на обновление"
// @Success 200
func updateSettings(ctx *gin.Context, userService user.Service) {
	var req models.UserSettingsUpdate
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	clientID := ctx.GetHeader("x-user-id")
	clientIDInt, err := strconv.Atoi(clientID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Некорректный client ID"})
		return
	}

	if err := req.UserSettingsUpdateValidate(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = userService.UpdateSettings(ctx, clientIDInt, req)
	if err != nil {
		resp := errors.GetHTTPStatus(err)
		ctx.JSON(resp.Status, gin.H{"error": resp.Message})
		return
	}

}

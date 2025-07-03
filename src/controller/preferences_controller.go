package controller

import (
	repository "gin/db/generated"
	"log"
	"net/http"
	"strconv"

	"gin/service"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type PreferencesController struct {
	service *service.PreferenceService
}

func NewPreferencesController(conn *pgx.Conn) *PreferencesController {
	service := service.NewPreferenceService(conn)

	return &PreferencesController{
		service: service,
	}
}

func (c *PreferencesController) GetUserPreferences(ctx *gin.Context) {
	var preferencesParams repository.GetUserPreferencesParams

	limit, offset := GetLimmitAndOffset(ctx)

	preferencesParams.Limit = int32(limit)
	preferencesParams.Offset = int32(offset)

	log.Printf("Limit %d", preferencesParams.Limit)

	userStr := ctx.Param("userId")

	userId, err := strconv.Atoi(userStr)
	if err != nil {
		limit = 10
	}

	preferencesParams.UserID = int32(userId)

	log.Printf("UserID %d", userStr )

	preferences, err := c.service.GetUserPreferences(preferencesParams)
	if err != nil {
		ctx.Error(gin.Error{
			Err:  err,
			Type: gin.ErrorTypePublic})
		return
	}

	ctx.JSON(http.StatusOK, preferences)

}
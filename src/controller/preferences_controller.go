package controller

import (
	repository "gin/db/generated"
	"log"
	"net/http"

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

	userToken, _ := ctx.Get("userID")
	preferencesParams.UserID = userToken.(int32)

	preferences, err := c.service.GetUserPreferences(preferencesParams)
	if err != nil {
		ctx.Error(gin.Error{
			Err:  err,
			Type: gin.ErrorTypePublic})
		return
	}

	ctx.JSON(http.StatusOK, preferences)

}

func (c *PreferencesController) AddPreference(ctx *gin.Context) {
	var addPreferenceParams repository.AddPreferenceParams

	if err := ctx.ShouldBind(&addPreferenceParams); err != nil {
		ctx.Error(gin.Error{
			Err:  err,
			Type: gin.ErrorTypePublic})
		return
	}

	userToken, _ := ctx.Get("userID")
	addPreferenceParams.UserID = userToken.(int32)

	err := c.service.AddPreference(addPreferenceParams)

	if err != nil {
		ctx.Error(gin.Error{
			Err:  err,
			Type: gin.ErrorTypePublic})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
	
}

func (c *PreferencesController) DeletePreference(ctx *gin.Context) {
	var deletePreferenceParams repository.DeletePreferenceParams

	if err := ctx.ShouldBind(&deletePreferenceParams); err != nil {
		ctx.Error(gin.Error{
			Err:  err,
			Type: gin.ErrorTypePublic})
		return
	}

	userToken, _ := ctx.Get("userID")
	deletePreferenceParams.UserID = userToken.(int32)

	err := c.service.DeletePreference(deletePreferenceParams)

	if err != nil {
		ctx.Error(gin.Error{
			Err:  err,
			Type: gin.ErrorTypePublic})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}
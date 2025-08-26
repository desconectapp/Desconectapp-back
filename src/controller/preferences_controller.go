package controller

import (
	"errors"
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

	preferencesParams.Limit = int32(limit) + 1
	preferencesParams.Offset = int32(offset)

	
	userToken, _ := ctx.Get("userID")
	preferencesParams.UserID = userToken.(int32)
	
	preferences, err := c.service.GetUserPreferences(preferencesParams)

	hasMore := len(preferences) == int(preferencesParams.Limit)

	if hasMore {
		preferences = preferences[:len(preferences)-1]
	}

	result := PaginatedPreferences{Preferences: preferences, HasMore: hasMore}

	if err != nil {
		ctx.Error(gin.Error{
			Err:  err,
			Type: gin.ErrorTypePublic})
		return
	}

	ctx.JSON(http.StatusOK, result)

}

func (c *PreferencesController) AddPreference(ctx *gin.Context) {
	var addPreferenceParams repository.AddPreferenceParams

	if err := ctx.ShouldBind(&addPreferenceParams); err != nil {
		ctx.Error(gin.Error{
			Err:  err,
			Type: gin.ErrorTypePublic})
		return
	}

	if addPreferenceParams.ActivityID == 0 {
		ctx.Error(gin.Error{
			Err:  errors.New("activity_id cannot be 0"),
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

	res := ActivityIdResponse{
		ActivityPreferenseID:  addPreferenceParams.ActivityID,
	}

	ctx.JSON(http.StatusOK, res)
	
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

	res := ActivityIdResponse{ActivityPreferenseID: deletePreferenceParams.ActivityID}

	ctx.JSON(http.StatusOK, res)
}

func (c *PreferencesController) BatchAddUserPreferences(ctx *gin.Context) {
	var userPreferences repository.BatchAddPreferencesParams

	if err := ctx.ShouldBindJSON(&userPreferences); err != nil {
		ctx.Error(gin.Error{
			Err:  err,
			Type: gin.ErrorTypePublic})
		return
	}

	log.Printf("preferences %v", userPreferences.ActivityIds)

	if len(userPreferences.ActivityIds) == 0 || len(userPreferences.ActivityIds) > 50 {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "preferences must contain between 1 and 50 activity IDs",
		})
		return
	}

	userToken, _ := ctx.Get("userID")
	userPreferences.UserID = userToken.(int32)

	err := c.service.BatchAddPreferences(userPreferences)

	if err != nil {
		ctx.Error(gin.Error{
			Err:  err,
			Type: gin.ErrorTypePublic})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"result": "preferences added successfully",
	})
}
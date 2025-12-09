package controller

import (
	"database/sql"
	"errors"
	repository "gin/db/generated"
	"log"
	"net/http"

	"gin/service"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PreferencesController struct {
	service *service.PreferenceService
}

func NewPreferencesController(conn *pgxpool.Pool) *PreferencesController {
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

	log.Println(preferencesParams)

	preferences, err := c.service.GetUserPreferences(preferencesParams)

	if err != nil {
		ErrorNoStatus(ctx, err)
		return
	}

	hasMore := len(preferences) == int(preferencesParams.Limit)

	if hasMore {
		preferences = preferences[:len(preferences)-1]
	}

	result := PaginatedPreferences{Preferences: preferences, HasMore: hasMore}

	ctx.JSON(http.StatusOK, result)

}

func (c *PreferencesController) AddPreference(ctx *gin.Context) {
	var addPreferenceParams repository.AddPreferenceParams

	if err := ctx.ShouldBind(&addPreferenceParams); err != nil {
		ErrorWithStatus(ctx, "Could not bind", http.StatusBadRequest)
		return
	}

	userToken, _ := ctx.Get("userID")
	addPreferenceParams.UserID = userToken.(int32)

	err := c.service.AddPreference(addPreferenceParams)

	if err != nil {
		ErrorWithStatus(ctx, err.Error(), http.StatusNotFound)
		return
	}
	res := ActivityIdResponse{
		ActivityPreferenseID: addPreferenceParams.ActivityID,
	}

	ctx.JSON(http.StatusOK, res)

}

func (c *PreferencesController) DeletePreference(ctx *gin.Context) {
	var deletePreferenceParams repository.DeletePreferenceParams

	if err := ctx.ShouldBind(&deletePreferenceParams); err != nil {
		ErrorWithStatus(ctx, "Could not bind", http.StatusBadRequest)
		return
	}

	userToken, _ := ctx.Get("userID")
	deletePreferenceParams.UserID = userToken.(int32)

	id, err := c.service.DeletePreference(deletePreferenceParams)

	if errors.Is(err, sql.ErrNoRows) {
		ErrorWithStatus(ctx, "Preference not found", http.StatusNotFound)
		return
	} else if err != nil {
		ErrorNoStatus(ctx, err)
		return
	}

	res := ActivityIdResponse{ActivityPreferenseID: id}

	ctx.JSON(http.StatusOK, res)
}

func (c *PreferencesController) BatchAddUserPreferences(ctx *gin.Context) {
	type batchRequest struct {
		ActivityIds      []int32  `json:"activity_ids"`
		CustomActivities []string `json:"custom_activities"`
	}

	var req batchRequest

	if err := ctx.ShouldBind(&req); err != nil {
		ErrorWithStatus(ctx, "Bad request", http.StatusBadRequest)
		return
	}

	var params repository.BatchAddPreferencesParams

	userToken, _ := ctx.Get("userID")
	params.UserID = userToken.(int32)
	params.ActivityIds = req.ActivityIds

	for _, customActivity := range req.CustomActivities {
		newActivityID, err := c.service.GenerateCustomActivity(customActivity)
		if err != nil {
			log.Println(err)
			ErrorWithStatus(ctx, err.Error(), http.StatusInternalServerError)
			return
		}
		params.ActivityIds = append(params.ActivityIds, newActivityID)
	}

	err := c.service.BatchAddPreferences(params)

	if err != nil {
		ErrorWithStatus(ctx, err.Error(), http.StatusNotFound)
		return
	}

	res := ActivityIdBatchResponse{
		ActivityIdBatchIDs: req.ActivityIds,
	}

	ctx.JSON(http.StatusOK, res)
}

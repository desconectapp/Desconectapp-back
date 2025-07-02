package controller

import (
	repository "gin/db/generated"
	"gin/service"
	"net/http"


	"log"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type ActivitiesController struct {
	service *service.ActivitiesRequestService
}

func NewActivitesController(conn *pgx.Conn) *ActivitiesController {
	service := service.NewActivitiesRequestService(conn)
	return &ActivitiesController{
		service: service,
	}
}

func (c *ActivitiesController) ListActivitiesRequests(ctx *gin.Context) {
	var activityParams repository.ListActivityRequestsParams

	limit, offset := GetLimmitAndOffset(ctx)

	activityParams.Limit = int32(limit)
	activityParams.Offset = int32(offset)

	activities, err := c.service.ListActivitiesRequests(activityParams)

	if err != nil {
		ctx.Error(gin.Error{
			Err:  err,
			Type: gin.ErrorTypePublic})
		return
	}
	ctx.JSON(http.StatusOK, activities)
}

func (c *ActivitiesController) CreateActivityRequest(ctx *gin.Context) {
	var activityParams repository.CreateActivityRequestParams
	if err := ctx.ShouldBind(&activityParams); err != nil {
		ctx.Error(gin.Error{
			Err:  err,
			Type: gin.ErrorTypePublic})
		return
	}

	activity, err := c.service.CreateActivityRequest(activityParams)
	if err != nil {
		ctx.Error(gin.Error{
			Err:  err,
			Type: gin.ErrorTypePublic})
		return
	}

	ctx.JSON(http.StatusOK, activity)
}

func (c *ActivitiesController) GetActivities(ctx *gin.Context) {
	var activityParams repository.GetActivitiesParams

	limit, offset := GetLimmitAndOffset(ctx)

	activityParams.Limit = int32(limit)
	activityParams.Offset = int32(offset)

	log.Printf("Limit = %d, Offset = %d", activityParams.Limit, activityParams.Offset)

	activities, err := c.service.GetActivities(activityParams)

	if err != nil {
		ctx.Error(gin.Error{
			Err:  err,
			Type: gin.ErrorTypePublic})
		return
	}

	ctx.JSON(http.StatusOK, activities)

}
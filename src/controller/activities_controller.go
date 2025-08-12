package controller

import (
	repository "gin/db/generated"
	"gin/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type ActivitiesService interface {
	ListActivities() ([]repository.ActivityRequest, error)
	CreateActivity(repository.CreateActivityRequestParams) (repository.ActivityRequest, error)
}

type ActivitiesController struct {
	service ActivitiesService
}

func NewActivitesController(conn *pgx.Conn) *ActivitiesController {
	service := service.NewActivitiesService(conn)
	return &ActivitiesController{
		service: service,
	}
}

func NewActivitiesControllerWithService(service ActivitiesService) *ActivitiesController {
	return &ActivitiesController{service: service}
}

func (c *ActivitiesController) ListActivities(ctx *gin.Context) {
	activities, err := c.service.ListActivities()
	if err != nil {
		ctx.Error(gin.Error{
			Err:  err,
			Type: gin.ErrorTypePublic})
		return
	}
	ctx.JSON(http.StatusOK, activities)
}

func (c *ActivitiesController) CreateActivity(ctx *gin.Context) {
	var activityParams repository.CreateActivityRequestParams
	if err := ctx.ShouldBind(&activityParams); err != nil {
		ctx.Error(gin.Error{
			Err:  err,
			Type: gin.ErrorTypePublic})
		return
	}

	activity, err := c.service.CreateActivity(activityParams)
	if err != nil {
		ctx.Error(gin.Error{
			Err:  err,
			Type: gin.ErrorTypePublic})
		return
	}

	ctx.JSON(http.StatusOK, activity)
}

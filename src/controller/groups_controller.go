package controller

import (
	repository "gin/db/generated"
	"gin/service"
	"net/http"



	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type GroupsController struct {
	service *service.GroupsService
}

func NewGroupsController(conn *pgx.Conn) *GroupsController {
	service := service.NewGroupsService(conn)
	return &GroupsController{
		service: service,
	}
}

func (c *GroupsController) ListGroups(ctx *gin.Context) {
	var groupParams repository.ListGroupsParams

	limit, offset := GetLimmitAndOffset(ctx)

	groupParams.Limit = int32(limit)
	groupParams.Offset = int32(offset)

	groupsList, err := c.service.ListGroups(groupParams)

	if err != nil {
		ctx.Error(gin.Error{
			Err:  err,
			Type: gin.ErrorTypePublic})
		return
	}
	ctx.JSON(http.StatusOK, groupsList)
}

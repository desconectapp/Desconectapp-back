package controller

import (
	repository "gin/db/generated"
	"gin/service"
	"net/http"
	"strconv"

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

type GroupInfo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Location    string `json:"location"`
	ActivityID  int32  `json:"activity_id"`
	MembersIds []int32 `json:"members_ids"`
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

func (c *GroupsController) GetGroup(ctx *gin.Context) {
	groupStr := ctx.Param("groupId")

	groupId, err := strconv.Atoi(groupStr)
	
	if err != nil {
		ctx.Error(gin.Error{
			Err:  err,
			Type: gin.ErrorTypePublic})
		return
	}

	group, err := c.service.GetGroup(int32(groupId))

	if err != nil {
		ctx.Error(gin.Error{
			Err:  err,
			Type: gin.ErrorTypePublic})
		return
	}

	ctx.JSON(http.StatusOK, group)	
}

func (c *GroupsController) CreateGroup(ctx *gin.Context) {

	groupParams, batchParams, err := getGroupParams(ctx)

	if err != nil {
		ctx.Error(gin.Error{
			Err:  err,
			Type: gin.ErrorTypePublic,})
		return
	}

	id, err := c.service.CreateGroup(groupParams, batchParams)
	if err != nil {
		ctx.Error(gin.Error{
			Err:  err,
			Type: gin.ErrorTypePublic,})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"updated": id,
	})
}

func getGroupParams(ctx *gin.Context) (repository.CreateGroupParams, repository.BatchAddUserToGroupParams, error) {
	var groupInfo GroupInfo
	var groupParams repository.CreateGroupParams
	var groupBatchParams repository.BatchAddUserToGroupParams
	
	
	groupId := ctx.Param("groupId")
	stringId, err := strconv.Atoi(groupId)
	if err != nil {
		return groupParams, groupBatchParams, err
	}

	if err = ctx.ShouldBind(&groupInfo); err != nil {
		return groupParams, groupBatchParams, err
	}

	groupParams.ActivityID = &groupInfo.ActivityID
	groupParams.Description = &groupInfo.Description
	groupParams.Location = &groupInfo.Location
	groupParams.Name = &groupInfo.Name

	groupBatchParams.GroupID = int32(stringId)
	groupBatchParams.UserID = groupInfo.MembersIds
	
	return groupParams, groupBatchParams, nil
}
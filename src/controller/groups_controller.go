package controller

import (
	repository "gin/db/generated"
	"gin/service"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type GroupsController struct {
	service *service.GroupsService
}

func NewGroupsController(conn *pgxpool.Pool) *GroupsController {
	service := service.NewGroupsService(conn)
	return &GroupsController{
		service: service,
	}
}

type GroupInfo struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Location    string  `json:"location"`
	ActivityID  int32   `json:"activity_id"`
	MembersIds  []int32 `json:"members_ids"`
}

func (c *GroupsController) ListGroups(ctx *gin.Context) {
	var groupParams repository.ListGroupsParams

	limit, offset := GetLimmitAndOffset(ctx)

	groupParams.Limit = int32(limit) + 1
	groupParams.Offset = int32(offset)

	groupsList, err := c.service.ListGroups(groupParams)

	hasMore := len(groupsList) == int(groupParams.Limit)

	if hasMore {
		groupsList = groupsList[:len(groupsList)-1]
	}

	result := PaginatedGroups{Groups: groupsList, HasMore: hasMore}

	if err != nil {
		ErrorNoStatus(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, result)
}

func (c *GroupsController) ListUserGroups(ctx *gin.Context) {
	var groupParams repository.ListUserGroupsParams

	limit, offset := GetLimmitAndOffset(ctx)

	groupParams.Limit = int32(limit) + 1
	groupParams.Offset = int32(offset)

	userToken, _ := ctx.Get("userID")
	groupParams.UserID = userToken.(int32)

	groupsList, err := c.service.ListUserGroups(groupParams)

	hasMore := len(groupsList) == int(groupParams.Limit)

	if hasMore {
		groupsList = groupsList[:len(groupsList)-1]
	}

	result := PaginatedUserGroup{Groups: groupsList, HasMore: hasMore}

	if err != nil {
		ErrorNoStatus(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, result)
}

func (c *GroupsController) GetGroup(ctx *gin.Context) {
	groupStr := ctx.Param("groupId")

	groupId, err := strconv.Atoi(groupStr)

	if err != nil {
		ErrorNoStatus(ctx, err)
		return
	}

	group, err := c.service.GetGroup(int32(groupId))

	if err != nil {
		ErrorNoStatus(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, group)
}

func (c *GroupsController) CreateGroup(ctx *gin.Context) {

	var groupInfo repository.CreateGroupParams

	if err := ctx.ShouldBind(&groupInfo); err != nil {
		ErrorWithStatus(ctx, "Could not bind", http.StatusBadRequest)
		return
	}

	group, err := c.service.CreateGroup(groupInfo)

	log.Println(err)

	if err != nil {
		ErrorNoStatus(ctx, err)
		return
	}

	res := NewGroup{
		ID:           group.ID,
		Name:         group.Name,
		Description:  group.Description,
		Location:     group.Location,
		ActivityID:   group.ActivityID,
		ActivityName: group.ActivityName,
		ActivityIcon: group.ActivityIcon,
		Members:      group.Members,
		Status:       group.Status,
	}

	ctx.JSON(http.StatusCreated, res)
}

func (c *GroupsController) ExitGroup(ctx *gin.Context) {
	var exitParams repository.ExitGroupParams

	log.Printf("Called")

	groupIdStr := ctx.Param("groupId")
	groupId, err := strconv.Atoi(groupIdStr)
	if err != nil {
		ErrorNoStatus(ctx, err)
		return
	}
	exitParams.GroupID = int32(groupId)

	userToken, _ := ctx.Get("userID")
	exitParams.UserID = userToken.(int32)

	err = c.service.ExitGroup(exitParams)
	if err != nil {
		ErrorNoStatus(ctx, err)
		return
	}

	res := GroupIdResponse{GroupId: exitParams.GroupID}

	ctx.JSON(http.StatusOK, res)
}

func (c *GroupsController) DeleteGroup(ctx *gin.Context) {
	groupIdStr := ctx.Param("groupId")
	groupId, err := strconv.Atoi(groupIdStr)
	if err != nil {
		ErrorNoStatus(ctx, err)
		return
	}

	id, err := c.service.DeleteGroup(int32(groupId))
	if err != nil {
		ErrorNoStatus(ctx, err)
		return
	}

	res := GroupIdResponse{GroupId: id}
	ctx.JSON(http.StatusOK, res)
}

func (c *GroupsController) UpdateGroupDescription(ctx *gin.Context) {
	var descriptionParams repository.UpdateGroupDescriptiomParams

	if err := ctx.ShouldBind(&descriptionParams); err != nil {
		ErrorWithStatus(ctx, "Could not bind", http.StatusBadRequest)
		return
	}

	groupIdStr := ctx.Param("groupId")
	groupId, err := strconv.Atoi(groupIdStr)
	if err != nil {
		ErrorNoStatus(ctx, err)
		return
	}
	descriptionParams.ID = int32(groupId)

	err = c.service.UpdateGroupDescription(descriptionParams)

	if err != nil {
		ErrorNoStatus(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}

func (c *GroupsController) ChangeGroupStatus(ctx *gin.Context) {
	var descriptionParams repository.ChangeGroupStatusParams

	if err := ctx.ShouldBind(&descriptionParams); err != nil {
		ErrorWithStatus(ctx, "Could not bind", http.StatusBadRequest)
		return
	}

	groupIdStr := ctx.Param("groupId")
	groupId, err := strconv.Atoi(groupIdStr)
	if err != nil {
		ErrorNoStatus(ctx, err)
		return
	}
	descriptionParams.ID = int32(groupId)

	err = c.service.ChangeGroupStatus(descriptionParams)

	if err != nil {
		ErrorNoStatus(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}

func (c *GroupsController) GetOpenGroups(ctx *gin.Context) {
	var filter service.ActivityFilter

	if err := ctx.ShouldBind(&filter); err != nil {
		ErrorWithStatus(ctx, "Could not bind", http.StatusBadRequest)
		return
	}

	limit, offset := GetLimmitAndOffset(ctx)

	filter.Limit = int32(limit) + 1
	filter.Offset = int32(offset)

	groups, err := c.service.GetStatusOpenGroups(filter)

	log.Println(groups)

	hasMore := len(groups) == int(filter.Limit)
	if hasMore {
		groups = groups[:len(groups)-1]
	}

	if err != nil {
		ErrorNoStatus(ctx, err)
		return
	}

	res := PaginatedOpenGroup{Groups: groups, HasMore: hasMore}

	ctx.JSON(http.StatusOK, res)
}

func (c *GroupsController) GetUserRecommendations(ctx *gin.Context) {
	var filter repository.GetPreferredGroupsParams

	userToken, _ := ctx.Get("userID")
	filter.UserID = userToken.(int32)

	limit, offset := GetLimmitAndOffset(ctx)

	filter.Limit = int32(limit) + 1
	filter.Offset = int32(offset)

	groups, err := c.service.GetUserRecommendations(filter)

	hasMore := len(groups) == int(filter.Limit)
	if hasMore {
		groups = groups[:len(groups)-1]
	}

	if err != nil {
		ErrorNoStatus(ctx, err)
		return
	}

	res := PaginatedOpenGroup{Groups: groups, HasMore: hasMore}

	ctx.JSON(http.StatusOK, res)
}
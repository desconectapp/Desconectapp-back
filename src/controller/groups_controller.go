package controller

import (
	repository "gin/db/generated"
	"gin/service"
	"log"
	"net/http"
	"strconv"
	"strings"

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

	if len(groupInfo.UserIds) == 0 {
		userToken, _ := ctx.Get("userID")
		groupInfo.UserIds = append(groupInfo.UserIds, userToken.(int32))
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
		LocationName: group.LocationName,
		ActivityID:   group.ActivityID,
		ActivityName: group.ActivityName,
		ActivityIcon: group.ActivityIcon,
		Members:      group.Members,
		Public:       group.Public,
	}

	log.Println(res)

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

func (c *GroupsController) JoinGroup(ctx *gin.Context) {
	var joinParams repository.AddUserToGroupParams

	log.Printf("Called")

	groupIdStr := ctx.Param("groupId")
	groupId, err := strconv.Atoi(groupIdStr)
	if err != nil {
		ErrorNoStatus(ctx, err)
		return
	}
	joinParams.GroupID = int32(groupId)

	userToken, _ := ctx.Get("userID")
	joinParams.UserID = userToken.(int32)

	err = c.service.JoinGroup(joinParams)
	if err != nil {
		ErrorNoStatus(ctx, err)
		return
	}

	res := GroupIdResponse{GroupId: joinParams.GroupID}

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

func (c *GroupsController) UpdateGroupName(ctx *gin.Context) {
	var descriptionParams repository.ChangeGroupNameParams

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

	err = c.service.ChangeGroupName(descriptionParams)

	if err != nil {
		ErrorNoStatus(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}

func (c *GroupsController) UpdateGroupLocation(ctx *gin.Context) {
	var descriptionParams repository.ChangeGroupLocationParams

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

	err = c.service.ChangeGroupLocation(descriptionParams)

	if err != nil {
		ErrorNoStatus(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}

func (c *GroupsController) ChangeGroupStatus(ctx *gin.Context) {
	var publicParam struct {
		PublicG *bool `json:"public_g" binding:"required"`
	}

	if err := ctx.ShouldBind(&publicParam); err != nil {
		ErrorWithStatus(ctx, "Could not bind", http.StatusBadRequest)
		return
	}

	var descriptionParams repository.ChangeGroupPublicParams

	groupIdStr := ctx.Param("groupId")
	groupId, err := strconv.Atoi(groupIdStr)
	if err != nil {
		ErrorNoStatus(ctx, err)
		return
	}

	descriptionParams.ID = int32(groupId)
	descriptionParams.Public = publicParam.PublicG

	err = c.service.ChangeGroupPublic(descriptionParams)

	if err != nil {
		ErrorNoStatus(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}

func (c *GroupsController) UpdateGroupAvatar(ctx *gin.Context) {
	var body struct {
		AvatarUrl *string `json:"avatar_url"`
	}
	if err := ctx.ShouldBind(&body); err != nil {
		ErrorWithStatus(ctx, "Could not bind", http.StatusBadRequest)
		return
	}
	groupIdStr := ctx.Param("groupId")
	groupId, err := strconv.Atoi(groupIdStr)
	if err != nil {
		ErrorNoStatus(ctx, err)
		return
	}
	params := repository.UpdateGroupAvatarParams{
		ID:        int32(groupId),
		AvatarUrl: body.AvatarUrl,
	}
	if err := c.service.UpdateGroupAvatar(params); err != nil {
		ErrorNoStatus(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{})
}

func (c *GroupsController) GetOpenGroups(ctx *gin.Context) {
	var filter service.ActivityFilter

	limit, offset := GetLimmitAndOffset(ctx)
	filter.Limit = int32(limit) + 1
	filter.Offset = int32(offset)

	// Get user ID from JWT token if available
	if userToken, exists := ctx.Get("userID"); exists {
		userID := userToken.(int32)
		filter.UserID = &userID
	}

	// Parse myPreferences flag
	if myPrefsStr := ctx.Query("myPreferences"); myPrefsStr == "true" {
		if filter.UserID == nil {
			ErrorWithStatus(ctx, "Authentication required for myPreferences", http.StatusUnauthorized)
			return
		}
		filter.MyPreferences = true
	}
	
	// Parse activities list
	if activitiesStr := ctx.Query("activities"); activitiesStr != "" {
		activityIds := strings.Split(activitiesStr, ",")
		for _, idStr := range activityIds {
			if activityId, err := strconv.ParseInt(strings.TrimSpace(idStr), 10, 32); err == nil {
				filter.ActivityIds = append(filter.ActivityIds, int32(activityId))
			}
		}
	}

	// Parse single activity_id from query params (for backwards compatibility)
	if activityIdStr := ctx.Query("activity_id"); activityIdStr != "" {
		if activityId, err := strconv.ParseInt(activityIdStr, 10, 32); err == nil {
			filter.ActivityId = int32(activityId)
		}
	}

	// Parse location parameters from query params
	if latStr := ctx.Query("latitude"); latStr != "" {
		if lat, err := strconv.ParseFloat(latStr, 64); err == nil {
			filter.Latitude = &lat
		}
	}

	if lngStr := ctx.Query("longitude"); lngStr != "" {
		if lng, err := strconv.ParseFloat(lngStr, 64); err == nil {
			filter.Longitude = &lng
		}
	}

	if radiusStr := ctx.Query("radius"); radiusStr != "" {
		if radius, err := strconv.ParseFloat(radiusStr, 64); err == nil {
			filter.Radius = &radius
		}
	}

	// Validate that if any location param is provided, all are provided
	if (filter.Latitude != nil || filter.Longitude != nil || filter.Radius != nil) &&
		(filter.Latitude == nil || filter.Longitude == nil || filter.Radius == nil) {
		ErrorWithStatus(ctx, "latitude, longitude, and radius must all be provided together", http.StatusBadRequest)
		return
	}

	groups, err := c.service.GetOpenGroups(filter)

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

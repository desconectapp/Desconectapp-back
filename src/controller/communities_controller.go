package controller

import (
	repository "gin/db/generated"
	"gin/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CommunitiesController struct {
	service *service.CommunitiesService
}

func NewCommunitiesController(conn *pgxpool.Pool) *CommunitiesController {
	service := service.NewCommunitiesService(conn)
	return &CommunitiesController{
		service: service,
	}
}

func (c *CommunitiesController) CreateCommunity(ctx *gin.Context) {
	var communityInfo repository.CreateCommunityParams

	if err := ctx.ShouldBind(&communityInfo); err != nil {
		ErrorWithStatus(ctx, "Could not bind", http.StatusBadRequest)
		return
	}

	userToken, _ := ctx.Get("userID")
	if len(communityInfo.UserIds) == 0 {
		communityInfo.UserIds = append(communityInfo.UserIds, userToken.(int32))
		communityInfo.AdminUserIds = append(communityInfo.AdminUserIds, userToken.(int32))
	}

	community, err := c.service.CreateCommunity(communityInfo)
	if err != nil {
		ErrorNoStatus(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, community)
}

func (c *CommunitiesController) ListUserCommunities(ctx *gin.Context) {
	var params repository.ListUserCommunitiesParams

	limit, offset := GetLimmitAndOffset(ctx)
	params.Limit = int32(limit) + 1
	params.Offset = int32(offset)

	userToken, _ := ctx.Get("userID")
	params.UserID = userToken.(int32)

	communities, err := c.service.ListUserCommunities(params)
	if err != nil {
		ErrorNoStatus(ctx, err)
		return
	}

	hasMore := len(communities) == int(params.Limit)
	if hasMore {
		communities = communities[:len(communities)-1]
	}

	ctx.JSON(http.StatusOK, gin.H{
		"communities": communities,
		"has_more":    hasMore,
	})
}

func (c *CommunitiesController) GetCommunity(ctx *gin.Context) {
	idStr := ctx.Param("communityId")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ErrorNoStatus(ctx, err)
		return
	}

	community, err := c.service.GetCommunity(int32(id))
	if err != nil {
		ErrorNoStatus(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, community)
}

func (c *CommunitiesController) DeleteCommunity(ctx *gin.Context) {
	idStr := ctx.Param("communityId")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ErrorNoStatus(ctx, err)
		return
	}

	deletedId, err := c.service.DeleteCommunity(int32(id))
	if err != nil {
		ErrorNoStatus(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"community_id": deletedId})
}

func (c *CommunitiesController) UpdateCommunityDescription(ctx *gin.Context) {
	var params repository.UpdateCommunityDescriptionParams

	if err := ctx.ShouldBind(&params); err != nil {
		ErrorWithStatus(ctx, "Could not bind", http.StatusBadRequest)
		return
	}

	idStr := ctx.Param("communityId")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ErrorNoStatus(ctx, err)
		return
	}
	params.ID = int32(id)

	err = c.service.UpdateCommunityDescription(params)
	if err != nil {
		ErrorNoStatus(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}

func (c *CommunitiesController) UpdateCommunityName(ctx *gin.Context) {
	var params repository.ChangeCommunityNameParams

	if err := ctx.ShouldBind(&params); err != nil {
		ErrorWithStatus(ctx, "Could not bind", http.StatusBadRequest)
		return
	}

	idStr := ctx.Param("communityId")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ErrorNoStatus(ctx, err)
		return
	}
	params.ID = int32(id)

	err = c.service.ChangeCommunityName(params)
	if err != nil {
		ErrorNoStatus(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}

func (c *CommunitiesController) UpdateCommunityLocation(ctx *gin.Context) {
	var params repository.ChangeCommunityLocationParams

	if err := ctx.ShouldBind(&params); err != nil {
		ErrorWithStatus(ctx, "Could not bind", http.StatusBadRequest)
		return
	}

	idStr := ctx.Param("communityId")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ErrorNoStatus(ctx, err)
		return
	}
	params.ID = int32(id)

	err = c.service.ChangeCommunityLocation(params)
	if err != nil {
		ErrorNoStatus(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}

func (c *CommunitiesController) UpdateCommunityAvatar(ctx *gin.Context) {
	var body struct {
		AvatarUrl *string `json:"avatar_url"`
	}

	if err := ctx.ShouldBind(&body); err != nil {
		ErrorWithStatus(ctx, "Could not bind", http.StatusBadRequest)
		return
	}

	idStr := ctx.Param("communityId")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ErrorNoStatus(ctx, err)
		return
	}

	params := repository.UpdateCommunityAvatarParams{
		ID:        int32(id),
		AvatarUrl: body.AvatarUrl,
	}

	err = c.service.UpdateCommunityAvatar(params)
	if err != nil {
		ErrorNoStatus(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{})
}

func (c *CommunitiesController) JoinCommunity(ctx *gin.Context) {
	var params repository.AddUserToCommunityParams

	idStr := ctx.Param("communityId")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ErrorNoStatus(ctx, err)
		return
	}
	params.CommunityID = int32(id)

	userToken, _ := ctx.Get("userID")
	params.UserID = userToken.(int32)

	err = c.service.AddUserToCommunity(params)
	if err != nil {
		ErrorNoStatus(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"community_id": params.CommunityID})
}

func (c *CommunitiesController) ExitCommunity(ctx *gin.Context) {
	var params repository.ExitCommunityParams

	idStr := ctx.Param("communityId")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ErrorNoStatus(ctx, err)
		return
	}
	params.CommunityID = int32(id)

	userToken, _ := ctx.Get("userID")
	params.UserID = userToken.(int32)

	err = c.service.ExitCommunity(params)
	if err != nil {
		ErrorNoStatus(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"community_id": params.CommunityID})
}

func (c *CommunitiesController) GetCommunitiesWithLocation(ctx *gin.Context) {
	var filter repository.GetCommunitiesWithLocationParams

	limit, offset := GetLimmitAndOffset(ctx)
	filter.Limit = int32(limit) + 1
	filter.Offset = int32(offset)

	if activityIdStr := ctx.Query("activity_id"); activityIdStr != "" {
		if activityId, err := strconv.ParseInt(activityIdStr, 10, 32); err == nil {
			activityId32 := int32(activityId)
			filter.ActivityID = &activityId32
		}
	}

	if latStr := ctx.Query("latitude"); latStr != "" {
		if lat, err := strconv.ParseFloat(latStr, 64); err == nil {
			filter.Latitude = lat
		}
	}
	if lngStr := ctx.Query("longitude"); lngStr != "" {
		if lng, err := strconv.ParseFloat(lngStr, 64); err == nil {
			filter.Longitude = lng
		}
	}
	if radiusStr := ctx.Query("radius"); radiusStr != "" {
		if radius, err := strconv.ParseFloat(radiusStr, 64); err == nil {
			filter.Radius = radius
		}
	}

	if (&filter.Latitude != nil || &filter.Longitude != nil || &filter.Radius != nil) &&
		(&filter.Latitude == nil || &filter.Longitude == nil || &filter.Radius == nil) {
		ErrorWithStatus(ctx, "latitude, longitude, and radius must all be provided together", http.StatusBadRequest)
		return
	}

	communities, err := c.service.GetCommunitiesWithLocation(filter)
	if err != nil {
		ErrorNoStatus(ctx, err)
		return
	}

	hasMore := len(communities) == int(filter.Limit)
	if hasMore {
		communities = communities[:len(communities)-1]
	}

	ctx.JSON(http.StatusOK, gin.H{
		"communities": communities,
		"has_more":    hasMore,
	})
}

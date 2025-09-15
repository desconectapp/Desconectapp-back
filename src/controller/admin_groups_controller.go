package controller

import (
	"gin/service"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type AdminGroupController struct {
	service *service.AdminGroupService
}

func NewAdminGroupController(conn *pgx.Conn) *AdminGroupController {
	svc := service.NewAdminGroupService(conn)
	return &AdminGroupController{
		service: svc,
	}
}

func (e *AdminGroupController) ListGroups(ctx *gin.Context) {
	limitEnd, _ := strconv.Atoi(ctx.DefaultQuery("_end", "25"))
	offsetStart, _ := strconv.Atoi(ctx.DefaultQuery("_start", "0"))

	sortField := ctx.DefaultQuery("_sort", "name")
	sortOrder := ctx.DefaultQuery("_order", "ASC")

	var namePtr *string
	if name := ctx.Query("name"); name != "" {
		namePtr = &name
	}

	groups, err := e.service.ListGroups(int32(limitEnd-offsetStart), int32(offsetStart), namePtr, sortField, sortOrder)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	total, err := e.service.CountGroups()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Header("X-Total-Count", strconv.Itoa(int(total)))
	ctx.JSON(http.StatusOK, groups)
}

func (c *AdminGroupController) GetGroup(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	group, err := c.service.GetGroup(int32(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, group)
}

func (c *AdminGroupController) CreateGroup(ctx *gin.Context) {
	var req struct {
		Name        string  `json:"name"`
		Description *string `json:"description"`
		Location    *string `json:"location"`
		ActivityID  *int32  `json:"activity_id"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	group, err := c.service.CreateGroup(req.Name, req.Description, req.Location, *req.ActivityID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, group)
}

// func (c *AdminGroupController) UpdateGroup(ctx *gin.Context) {
// 	id, _ := strconv.Atoi(ctx.Param("id"))
// 	var req struct {
// 		Name        string  `json:"name"`
// 		Description *string `json:"description"`
// 		Location    *string `json:"location"`
// 		ActivityID  *int32  `json:"activity_id"`
// 	}
// 	if err := ctx.ShouldBindJSON(&req); err != nil {
// 		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}
//
// 	err := c.service.UpdateGroup(int32(id), req.Name, req.Description, req.Location, *req.ActivityID)
// 	if err != nil {
// 		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
// 	ctx.Status(http.StatusNoContent)
// }

func (c *AdminGroupController) DeleteGroup(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	log.Println("Deleting group with ID:", id)
	if err := c.service.DeleteGroup(int32(id)); err != nil {
		println("Error deleting group:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.Status(http.StatusNoContent)
}

func (c *AdminGroupController) GetManyGroups(ctx *gin.Context) {
	ids := ctx.QueryArray("id")
	result := []interface{}{}

	for _, idStr := range ids {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			continue
		}
		group, err := c.service.GetGroup(int32(id))
		if err == nil {
			result = append(result, group)
		}
	}

	ctx.JSON(http.StatusOK, result)
}

func (c *AdminGroupController) ListGroupMembers(ctx *gin.Context) {
	groupID, _ := strconv.Atoi(ctx.Param("id"))
	log.Println("Group ID:", groupID)
	members, err := c.service.ListGroupMembers(int32(groupID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, members)
}

func (c *AdminGroupController) AddGroupMember(ctx *gin.Context) {
	groupID, _ := strconv.Atoi(ctx.Param("id"))
	var req struct {
		UserID *int32 `json:"user_id"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := c.service.AddGroupMember(int32(groupID), *req.UserID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.Status(http.StatusCreated)
}

func (c *AdminGroupController) RemoveGroupMember(ctx *gin.Context) {
	groupID, _ := strconv.Atoi(ctx.Param("id"))
	userID, _ := strconv.Atoi(ctx.Param("memberId"))
	if err := c.service.RemoveGroupMember(int32(groupID), int32(userID)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.Status(http.StatusNoContent)
}

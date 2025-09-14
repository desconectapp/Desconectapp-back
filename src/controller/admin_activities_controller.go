package controller

import (
	"gin/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type AdminActivityController struct {
	service *service.AdminActivityService
}

func NewAdminActivityController(conn *pgx.Conn) *AdminActivityController {
	service := service.NewAdminActivityService(conn)
	return &AdminActivityController{
		service: service,
	}
}

func (c *AdminActivityController) ListActivities(ctx *gin.Context) {
	limit, _ := strconv.Atoi(ctx.DefaultQuery("_end", "25"))
	offset, _ := strconv.Atoi(ctx.DefaultQuery("_start", "0"))

	sortField := ctx.DefaultQuery("_sort", "id")
	sortOrder := ctx.DefaultQuery("_order", "ASC")

	var namePtr, categoryPtr *string
	if name := ctx.Query("name"); name != "" {
		namePtr = &name
	}
	if category := ctx.Query("category"); category != "" {
		categoryPtr = &category
	}

	activities, err := c.service.ListActivities(int32(limit-offset), int32(offset), namePtr, categoryPtr, sortField, sortOrder)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	total, err := c.service.CountActivities(namePtr, categoryPtr)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Header("X-Total-Count", strconv.Itoa(int(total)))
	ctx.JSON(http.StatusOK, activities)
}

func (c *AdminActivityController) GetActivity(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	activity, err := c.service.GetActivity(int32(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, activity)
}

func (c *AdminActivityController) CreateActivity(ctx *gin.Context) {
	var req struct {
		Name     string  `json:"name"`
		Icon     *string `json:"icon"`
		Category *string `json:"category"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	icon := ""
	if req.Icon != nil {
		icon = *req.Icon
	}

	activity, err := c.service.CreateActivity(req.Name, icon, *req.Category)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, activity)
}

func (c *AdminActivityController) UpdateActivity(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	var req struct {
		Name     string  `json:"name"`
		Icon     *string `json:"icon"`
		Category *string `json:"category"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	icon := ""
	if req.Icon != nil {
		icon = *req.Icon
	}

	err := c.service.UpdateActivity(int32(id), req.Name, icon, *req.Category)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.Status(http.StatusNoContent)
}

func (c *AdminActivityController) DeleteActivity(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	if err := c.service.DeleteActivity(int32(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.Status(http.StatusNoContent)
}

func (c *AdminActivityController) GetManyActivities(ctx *gin.Context) {
	ids := ctx.QueryArray("id")
	result := []interface{}{}

	for _, idStr := range ids {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			continue
		}
		activity, err := c.service.GetActivity(int32(id))
		if err == nil {
			result = append(result, activity)
		}
	}

	ctx.JSON(http.StatusOK, result)
}

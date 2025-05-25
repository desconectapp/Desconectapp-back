package controller

import (
	"gin/service"
	"net/http"
	"github.com/gin-gonic/gin"
	"strconv"

	repository "gin/db/generated"
)

type Controller struct {
	service *service.Service
}

func NewController() *Controller {
	service := service.NewService()
	return &Controller{
		service: service,
	}
}


func (c *Controller) ListUsers(ctx *gin.Context) {
	users, err := c.service.ListUsers()
	if err != nil {
		ctx.Error(gin.Error{
			Err:  err,
			Type: gin.ErrorTypePublic,})
		return
	}
	ctx.JSON(http.StatusOK, users)
}

func (c *Controller) CreateUser(ctx *gin.Context) {
	var userParams repository.CreateUserParams
	if err := ctx.ShouldBind(&userParams); err != nil {
		ctx.Error(gin.Error{
			Err:  err,
			Type: gin.ErrorTypePublic,})
		return
	}

	user, err := c.service.CreateUser(userParams)
	if err != nil {
		ctx.Error(gin.Error{
			Err:  err,
			Type: gin.ErrorTypePublic,})
		return
	}
	ctx.JSON(http.StatusOK, user)
}

func (c *Controller) GetUser(ctx *gin.Context) {
	userId := ctx.Param("userId")
	stringId, err := strconv.Atoi(userId)
	if err != nil {
		ctx.Error(gin.Error{
			Err:  err,
			Type: gin.ErrorTypePublic,})
		return
	}

	user, err := c.service.GetUser(int32(stringId))
	if err != nil {
		ctx.Error(gin.Error{
			Err:  err,
			Type: gin.ErrorTypePublic,})
		return
	}
	ctx.JSON(http.StatusOK, user)
}

func (c *Controller) DeleteUser(ctx *gin.Context) {
	userId := ctx.Param("userId")
	stringId, err := strconv.Atoi(userId)
	if err != nil {
		ctx.Error(gin.Error{
			Err:  err,
			Type: gin.ErrorTypePublic,})
		return
	}

	id, err := c.service.DeleteUser(int32(stringId))
	if err != nil {
		ctx.Error(gin.Error{
			Err:  err,
			Type: gin.ErrorTypePublic,})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"deleted": id,
	})
}

func (c *Controller) UpdateUser(ctx *gin.Context) {
	userId := ctx.Param("userId")
	stringId, err := strconv.Atoi(userId)
	if err != nil {
		ctx.Error(gin.Error{
			Err:  err,
			Type: gin.ErrorTypePublic,})
		return
	}

	var userParams repository.UpdateUserParams
	if err := ctx.ShouldBind(&userParams); err != nil {
		ctx.Error(gin.Error{
			Err:  err,
			Type: gin.ErrorTypePublic,})
		return
	}

	userParams.ID = int32(stringId)

	id, err := c.service.UpdateUser(userParams)
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


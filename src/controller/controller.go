package controller

import (
	"errors"
	"gin/service"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"

	repository "gin/db/generated"
)

const (
	MIN_AGE = 15
	MAX_AGE = 100
)

type Controller struct {
	service *service.Service
}

func NewController(conn *pgx.Conn) *Controller {
	service := service.NewService(conn)
	return &Controller{
		service: service,
	}
}


func (c *Controller) ListUsers(ctx *gin.Context) {
	var userParams repository.ListUsersParams

	limit, offset := GetLimmitAndOffset(ctx)

	userParams.Limit = int32(limit) + 1
	userParams.Offset = int32(offset)

	users, err := c.service.ListUsers(userParams)

	hasMore := len(users) == int(userParams.Limit)

	if hasMore {
		users = users[:len(users)-1]
	}

	result := PaginatedUsers{Users: users, HasMore: hasMore}

	if err != nil {
		ctx.Error(gin.Error{
			Err:  err,
			Type: gin.ErrorTypePublic})
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (c *Controller) CreateProfile(ctx *gin.Context) {
	var profileData repository.CreateProfileParams
	if err := ctx.ShouldBind(&profileData); err != nil {
		ctx.Error(gin.Error{
			Err:  err,
			Type: gin.ErrorTypePublic})
		return
	}

	if profileData.Age < MIN_AGE || profileData.Age > MAX_AGE {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "age must be between 15 and 100",
		})
		return
	}

	userToken, _ := ctx.Get("userID")
	profileData.UserID = userToken.(int32)

	log.Printf("ID: %d", profileData.UserID)

	id, err := c.service.CreateProfile(profileData)
	if err != nil {
		ctx.Error(gin.Error{
			Err:  err,
			Type: gin.ErrorTypePublic})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"updated": id,
	})
}

func (c *Controller) GetUser(ctx *gin.Context) {
	userToken, _ := ctx.Get("userID")

	if userToken == 0 {
		ctx.Error(gin.Error{
			Err:  errors.New("user id cannot be 0"),
			Type: gin.ErrorTypePublic})
	}

	user, err := c.service.GetUser(userToken.(int32))
	if err != nil {
		ctx.Error(gin.Error{
			Err:  err,
			Type: gin.ErrorTypePublic})
		return
	}
	ctx.JSON(http.StatusOK, user)
}

func (c *Controller) DeleteUser(ctx *gin.Context) {
	userToken, _ := ctx.Get("userID")

	id, err := c.service.DeleteUser(userToken.(int32))
	if err != nil {
		ctx.Error(gin.Error{
			Err:  err,
			Type: gin.ErrorTypePublic})
		return
	}
	res := UserDeletedResponse{DeletedUserID: id}
	ctx.JSON(http.StatusOK, res)
}

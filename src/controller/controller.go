package controller

import (
	"errors"
	"gin/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"

	repository "gin/db/generated"
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

type UserUpdateInfo struct {
	Name             string  `json:"name"`
	Email            string  `json:"email"`
	Age              int32   `json:"age"`
	City             string  `json:"city"`
	CurrentSituation string  `json:"current_situation"`
	ActivityIDs      []int32 `json:"activity_ids"`
}

func (c *Controller) ListUsers(ctx *gin.Context) {
	var userParams repository.ListUsersParams

	limit, offset := GetLimmitAndOffset(ctx)

	userParams.Limit = int32(limit)
	userParams.Offset = int32(offset)

	users, err := c.service.ListUsers(userParams)
	if err != nil {
		ctx.Error(gin.Error{
			Err:  err,
			Type: gin.ErrorTypePublic})
		return
	}
	ctx.JSON(http.StatusOK, users)
}

func (c *Controller) CreateUser(ctx *gin.Context) {
	var userParams repository.CreateUserParams
	if err := ctx.ShouldBind(&userParams); err != nil {
		ctx.Error(gin.Error{
			Err:  err,
			Type: gin.ErrorTypePublic})
		return
	}

	user, err := c.service.CreateUser(userParams)
	if err != nil {
		ctx.Error(gin.Error{
			Err:  err,
			Type: gin.ErrorTypePublic})
		return
	}
	ctx.JSON(http.StatusOK, user)
}

func (c *Controller) AddUserPreferences(ctx *gin.Context) {
	var preferences []int32
	if err := ctx.ShouldBindJSON(&preferences); err != nil {
		ctx.Error(gin.Error{
			Err:  err,
			Type: gin.ErrorTypePublic})
		return
	}

	if len(preferences) == 0 || len(preferences) > 50 {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "preferences must contain between 1 and 50 activity IDs",
		})
		return
	}

	userId, err := strconv.Atoi(ctx.Param("userId"))
	if err != nil {
		ctx.Error(gin.Error{
			Err:  errors.New("userId is required and must be an integer"),
			Type: gin.ErrorTypePublic,
		})
		return
	}

	err = c.service.AddPreferences(int32(userId), preferences)
	if err != nil {
		ctx.Error(gin.Error{
			Err:  err,
			Type: gin.ErrorTypePublic})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"result": "preferences added successfully",
	})
}

func (c *Controller) CreateProfile(ctx *gin.Context) {
	var profileData repository.CreateProfileParams
	if err := ctx.ShouldBind(&profileData); err != nil {
		ctx.Error(gin.Error{
			Err:  err,
			Type: gin.ErrorTypePublic})
		return
	}

	if profileData.Age < 15 || profileData.Age > 100 {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "age must be between 15 and 100",
		})
		return
	}

	userId, err := strconv.Atoi(ctx.Param("userId"))
	if err != nil {
		ctx.Error(gin.Error{
			Err:  errors.New("userId is required and must be an integer"),
			Type: gin.ErrorTypePublic,
		})
		return
	}

	profileData.UserID = int32(userId)

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
	userId := ctx.Param("userId")
	stringId, err := strconv.Atoi(userId)
	if err != nil {
		ctx.Error(gin.Error{
			Err:  err,
			Type: gin.ErrorTypePublic})
		return
	}

	user, err := c.service.GetUser(int32(stringId))
	if err != nil {
		ctx.Error(gin.Error{
			Err:  err,
			Type: gin.ErrorTypePublic})
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
			Type: gin.ErrorTypePublic})
		return
	}

	id, err := c.service.DeleteUser(int32(stringId))
	if err != nil {
		ctx.Error(gin.Error{
			Err:  err,
			Type: gin.ErrorTypePublic})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"deleted": id,
	})
}

// func (c *Controller) UpdateUser(ctx *gin.Context) {
//
// 	userParams, userPreferences, err := getParams(ctx)
//
// 	if err != nil {
// 		ctx.Error(gin.Error{
// 			Err:  err,
// 			Type: gin.ErrorTypePublic})
// 		return
// 	}
//
// 	id, err := c.service.UpdateUser(userParams, userPreferences)
// 	if err != nil {
// 		ctx.Error(gin.Error{
// 			Err:  err,
// 			Type: gin.ErrorTypePublic})
// 		return
// 	}
// 	ctx.JSON(http.StatusOK, gin.H{
// 		"updated": id,
// 	})
// }

// func getParams(ctx *gin.Context) (repository.UpdateUserParams, repository.BatchAddPreferencesParams, error) {
// 	var userInfo UserUpdateInfo
// 	var userParams repository.UpdateUserParams
// 	var preferenceParams repository.BatchAddPreferencesParams
//
// 	userId := ctx.Param("userId")
// 	stringId, err := strconv.Atoi(userId)
// 	if err != nil {
// 		return userParams, preferenceParams, err
// 	}
//
// 	if err = ctx.ShouldBind(&userInfo); err != nil {
// 		return userParams, preferenceParams, err
// 	}
// 	userParams.ID = int32(stringId)
// 	userParams.Age = userInfo.Age
// 	userParams.City = userInfo.City
// 	userParams.CurrentSituation = userInfo.CurrentSituation
// 	userParams.Email = userInfo.Email
// 	userParams.Name = userInfo.Name
//
// 	preferenceParams.UserID = int32(stringId)
// 	preferenceParams.ActivityIds = userInfo.ActivityIDs
//
// 	return userParams, preferenceParams, nil
// }

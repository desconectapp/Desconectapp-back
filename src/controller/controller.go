package controller

import (
	"database/sql"
	"errors"
	"gin/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	repository "gin/db/generated"
)

const (
	MIN_AGE = 15
	MAX_AGE = 100
)

type Controller struct {
	service *service.Service
}

func NewController(conn *pgxpool.Pool) *Controller {
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

	if err != nil {
		ctx.Error(&gin.Error{
			Err:  err,
			Type: gin.ErrorTypePublic,
		})
		ctx.Abort()
		return
	}

	hasMore := len(users) == int(userParams.Limit)

	if hasMore {
		users = users[:len(users)-1]
	}

	result := PaginatedUsers{Users: users, HasMore: hasMore}

	ctx.JSON(http.StatusOK, result)
}

func (c *Controller) CreateProfile(ctx *gin.Context) {
	var profileData repository.CreateProfileParams
	if err := ctx.ShouldBind(&profileData); err != nil {
		ErrorWithStatus(ctx, "Could not bind", http.StatusBadRequest)
		return
	}

	if profileData.Age < MIN_AGE || profileData.Age > MAX_AGE {
		ErrorWithStatus(ctx, "Age must be between 15 and 100", http.StatusBadRequest)
		return
	}

	userToken, _ := ctx.Get("userID")
	profileData.UserID = userToken.(int32)

	user, err := c.service.CreateProfile(profileData)
	if err != nil {
		ErrorNoStatus(ctx, err)
		return
	}

	res := Profile{
		UserID:           user.UserID,
		Age:              user.Age,
		Name:             user.Name,
		City:             user.City,
		CurrentSituation: user.CurrentSituation,
		Gender:           user.Gender,
	}

	ctx.JSON(http.StatusOK, res)
}

func (c *Controller) GetUser(ctx *gin.Context) {
	userToken, _ := ctx.Get("userID")

	user, err := c.service.GetUser(userToken.(int32))

	if errors.Is(err, sql.ErrNoRows) {
		ErrorWithStatus(ctx, "The user does not exist", http.StatusNotFound)
		return
	} else if err != nil {
		ErrorNoStatus(ctx, err)
		return
	}

	res := Profile{
		UserID:           user.UserID,
		Age:              user.Age,
		Name:             user.Name,
		City:             user.City,
		CurrentSituation: user.CurrentSituation,
		Gender:           user.Gender,
	}

	ctx.JSON(http.StatusOK, res)
}

func (c *Controller) DeleteUser(ctx *gin.Context) {
	userToken, _ := ctx.Get("userID")

	id, err := c.service.DeleteUser(userToken.(int32))

	if errors.Is(err, sql.ErrNoRows) {
		ErrorWithStatus(ctx, "The user does not exist", http.StatusNotFound)
		return
	} else if err != nil {
		ErrorNoStatus(ctx, err)
		return
	}

	res := UserDeletedResponse{DeletedUserID: id}
	ctx.JSON(http.StatusOK, res)
}

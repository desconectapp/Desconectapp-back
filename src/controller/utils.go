package controller

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetLimmitAndOffset(ctx *gin.Context) (int32, int32) {
	limitStr := ctx.Query("limit")
	offsetStr := ctx.Query(("offset"))

	
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 0 {
		limit = 10
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	return int32(limit), int32(offset)
}

func ErrorNoStatus(ctx *gin.Context, err error) {
	ctx.Error(&gin.Error{
			Err:  err,
			Type: gin.ErrorTypePublic,
	})
	ctx.Abort()
}

func ErrorWithStatus(ctx *gin.Context, text string, status int) {
	ctx.Error(&gin.Error{
			Err:  errors.New(text),
			Type: gin.ErrorTypePublic,
		}).SetMeta(map[string]any{
			"status": status,
		})
	ctx.Abort()
}
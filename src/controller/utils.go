package controller

import (
	"strconv"
	"github.com/gin-gonic/gin"
)

func GetLimmitAndOffset(ctx *gin.Context) (int32, int32) {
	limitStr := ctx.Query("limit")
	offsetStr := ctx.Query(("offset"))

	
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 10
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		offset = 0
	}

	return int32(limit), int32(offset)
}
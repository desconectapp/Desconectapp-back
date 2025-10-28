package controller

import (
	repository "gin/db/generated"
	"gin/service"
	"log"
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


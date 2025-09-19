package controller

import (
	"gin/service"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type ChatsController struct {
	chatsService *service.ChatsService
}

func NewChatsController(conn *pgx.Conn) *ChatsController {
	service := service.NewChatsService(conn)
	return &ChatsController{
		chatsService: service,
	}
}

func (c *ChatsController) GetToken(ctx *gin.Context) {
	token, err := c.chatsService.GetToken()
	if err != nil {
		ctx.Error(gin.Error{
			Err:  err,
			Type: gin.ErrorTypePublic,
		})
	}
}

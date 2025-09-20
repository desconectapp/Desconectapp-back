package controller

import (
	"gin/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ChatsController struct {
	chatsService *service.ChatsService
}

func NewChatsController(conn *pgxpool.Pool) *ChatsController {
	service := service.NewChatsService(conn)
	return &ChatsController{
		chatsService: service,
	}
}

func (c *ChatsController) GetToken(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in token"})
		return
	}

	userIDInt32, ok := userID.(int32)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID format"})
		return
	}

	token, err := c.chatsService.GetToken(userIDInt32)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"supabase_token": token})
}

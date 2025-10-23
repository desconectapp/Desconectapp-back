package controller

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"gin/service"
	"net/http"
	"time"

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
	db      *repository.Queries
	ctx     context.Context
}

func NewController(conn *pgxpool.Pool) *Controller {
	service := service.NewService(conn)
	return &Controller{
		service: service,
		db:      repository.New(conn),
		ctx:     context.Background(),
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

func (c *Controller) UpdateProfile(ctx *gin.Context) {
	var profileData repository.UpdateProfileParams
	if err := ctx.ShouldBind(&profileData); err != nil {
		ErrorWithStatus(ctx, "Invalid json format", http.StatusBadRequest)
		return
	}

	if profileData.Age < MIN_AGE || profileData.Age > MAX_AGE {
		ErrorWithStatus(ctx, "Age must be between 15 and 100", http.StatusBadRequest)
		return
	}

	userId, _ := ctx.Get("userID")
	profileData.UserID = userId.(int32)

	user, err := c.service.UpdateProfile(profileData)
	if err != nil {
		ErrorWithStatus(ctx, "An error ocurred updating the profile", http.StatusBadRequest)
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

func (c *Controller) UpdateProfileAvatar(ctx *gin.Context) {
    var body struct {
        AvatarUrl *string `json:"avatar_url"`
    }
    if err := ctx.ShouldBind(&body); err != nil {
        ErrorWithStatus(ctx, "Invalid json format", http.StatusBadRequest)
        return
    }
    userId, _ := ctx.Get("userID")
    params := repository.UpdateProfileAvatarParams{
        UserID:    userId.(int32),
        AvatarUrl: body.AvatarUrl,
    }
    if err := c.service.UpdateProfileAvatar(params); err != nil {
        ErrorNoStatus(ctx, err)
        return
    }
    ctx.JSON(http.StatusOK, gin.H{})
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
		AvatarUrl:        user.AvatarUrl,
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

// RegisterPushToken registers a push token for a user
func (c *Controller) RegisterPushToken(ctx *gin.Context) {
	userID := ctx.GetInt("user_id")
	if userID == 0 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var req PushTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Store or update the push token
	_, err := c.db.CreatePushToken(c.ctx, repository.CreatePushTokenParams{
		UserID:   int32(userID),
		Token:    req.Token,
		Platform: req.Platform,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register push token"})
		return
	}

	ctx.JSON(http.StatusOK, PushTokenResponse{
		Success: true,
		Message: "Push token registered successfully",
	})
}

// UnregisterPushToken removes a push token
func (c *Controller) UnregisterPushToken(ctx *gin.Context) {
	token := ctx.Param("token")
	if token == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Token is required"})
		return
	}

	err := c.db.DeletePushToken(c.ctx, token)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unregister push token"})
		return
	}

	ctx.JSON(http.StatusOK, PushTokenResponse{
		Success: true,
		Message: "Push token unregistered successfully",
	})
}

// SendPushNotification sends a notification to a specific token
func (c *Controller) SendPushNotification(token, title, body string, data map[string]interface{}) error {
	// This is a placeholder - you'll need to implement the actual push notification sending
	// For now, we'll just log it
	fmt.Printf("Sending push notification to token %s: %s - %s\n", token, title, body)
	return nil
}

// SendPushNotificationToUser sends a notification to all tokens for a user
func (c *Controller) SendPushNotificationToUser(userID int32, title, body string, data map[string]interface{}) error {
	tokens, err := c.db.GetPushTokensByUser(c.ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to get user push tokens: %w", err)
	}

	for _, token := range tokens {
		err := c.SendPushNotification(token.Token, title, body, data)
		if err != nil {
			// Log error but continue with other tokens
			fmt.Printf("Failed to send notification to token %s: %v\n", token.Token, err)
		}
	}

	return nil
}

// SendPushNotificationToGroup sends a notification to all members of a group
func (c *Controller) SendPushNotificationToGroup(groupID int32, title, body string, data map[string]interface{}) error {
	tokens, err := c.db.GetPushTokensForGroup(c.ctx, groupID)
	if err != nil {
		return fmt.Errorf("failed to get group push tokens: %w", err)
	}

	for _, token := range tokens {
		err := c.SendPushNotification(token.Token, title, body, data)
		if err != nil {
			// Log error but continue with other tokens
			fmt.Printf("Failed to send notification to token %s: %v\n", token.Token, err)
		}
	}

	return nil
}

// TestPushNotification sends a test notification to the current user
func (c *Controller) TestPushNotification(ctx *gin.Context) {
	userID := ctx.GetInt("user_id")
	if userID == 0 {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Send test notification
	err := c.SendPushNotificationToUser(int32(userID), "Test Notification", "This is a test push notification from your app!", map[string]interface{}{
		"type": "test",
		"timestamp": time.Now().Unix(),
	})
	
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send test notification: " + err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Test notification sent successfully",
	})
}

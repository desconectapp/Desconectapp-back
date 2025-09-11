package controller

import (
	"errors"
	"fmt"
	"gin/service"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService  *service.AuthService
	emailService *service.EmailValidationService
}

func NewAuthController(authService *service.AuthService, emailService *service.EmailValidationService) *AuthController {
	return &AuthController{
		authService:  authService,
		emailService: emailService,
	}
}

func (c *AuthController) Login(ctx *gin.Context) {
	var loginReq LoginRequest
	if err := ctx.ShouldBind(&loginReq); err != nil {
		ErrorWithStatus(ctx, "Could not bind", http.StatusBadRequest)
		return
	}

	fmt.Println(loginReq.Email)

	session, err := c.authService.Login(loginReq.Email, loginReq.Password)
	if err != nil {
		ctx.Error(gin.Error{
			Err:  err,
			Type: gin.ErrorTypePublic,
		})
		return
	}

	ctx.JSON(http.StatusOK, AuthResponse{
		UserId:           session.UserId,
		Token:            session.Token,
		RefreshToken:     session.RefreshToken,
		ExpiresAt:        session.ExpiresAt,
		RefreshExpiresAt: session.RefreshExpiresAt,
	})
}

func (c *AuthController) Refresh(ctx *gin.Context) {
	var refreshReq RefreshRequest
	if err := ctx.ShouldBind(&refreshReq); err != nil {
		ErrorWithStatus(ctx, "Could not bind", http.StatusBadRequest)
		return
	}

	session, err := c.authService.RefreshToken(refreshReq.RefreshToken)
	if err != nil {
		ctx.Error(gin.Error{
			Err:  err,
			Type: gin.ErrorTypePublic,
		})
		return
	}

	ctx.JSON(http.StatusOK, AuthResponse{
		UserId:           session.UserId,
		Token:            session.Token,
		RefreshToken:     session.RefreshToken,
		ExpiresAt:        session.ExpiresAt,
		RefreshExpiresAt: session.RefreshExpiresAt,
	})
}

func (c *AuthController) AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("Authorization")
		if token == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "No token provided"})
			ctx.Abort()
			return
		}

		if len(token) > 7 && token[:7] == "Bearer " {
			token = token[7:]
		}

		userID, _, err := service.ValidateSession(token)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			ctx.Abort()
			return
		}

		ctx.Set("userID", userID)
		ctx.Next()
	}
}

func (c *AuthController) AdminMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := ctx.GetHeader("Authorization")
		if token == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "No token provided"})
			ctx.Abort()
			return
		}

		if len(token) > 7 && token[:7] == "Bearer " {
			token = token[7:]
		}

		userID, isAdmin, err := service.ValidateSession(token)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			ctx.Abort()
			return
		}

		ctx.Set("adminID", userID)
		ctx.Set("isAdmin", true)
		if !isAdmin {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
		} else {
			ctx.Next()
		}
	}
}

type SignupRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

func (c *AuthController) Signup(ctx *gin.Context) {
	var signupReq SignupRequest
	if err := ctx.ShouldBind(&signupReq); err != nil {
		ErrorWithStatus(ctx, "Could not bind", http.StatusBadRequest)
		return
	}

	if signupReq.Name == "" {
		signupReq.Name = "Test"
	}

	session, err := c.authService.Signup(signupReq.Name, signupReq.Email, signupReq.Password)
	if errors.Is(err, service.ErrUserExists) {
    	ctx.Error(&gin.Error{
			Err:  err,
			Type: gin.ErrorTypePublic,
		}).SetMeta(map[string]any{
			"status": http.StatusConflict,
		})
		return
	} else if err != nil {
		ctx.Error(gin.Error{
			Err:  err,
			Type: gin.ErrorTypePublic,
		})
		return
	}

	c.emailService.StartEmailVerification(session.UserId, signupReq.Email)

	ctx.JSON(http.StatusCreated, AuthResponse{
		UserId:           session.UserId,
		Token:            session.Token,
		RefreshToken:     session.RefreshToken,
		ExpiresAt:        session.ExpiresAt,
		RefreshExpiresAt: session.RefreshExpiresAt,
	})
}

type ValidateEmail struct {
	Code   string `json:"code" binding:"required,min=6"`
	UserID int32  `json:"user_id" binding:"required"`
}

func (c *AuthController) ValidateEmail(ctx *gin.Context) {
	var validation ValidateEmail
	if err := ctx.ShouldBind(&validation); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if validation.UserID == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}
	if validation.Code == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "code is required"})
		return
	}

	log.Println("Validating email for user ID:", validation.UserID, "with code:", validation.Code)

	err := c.emailService.ValidateEmailCode(validation.UserID, validation.Code)
	if err != nil {
		ctx.Error(gin.Error{
			Err:  err,
			Type: gin.ErrorTypePublic,
		})
		return
	}
}

func (c *AuthController) ResendValidationEmail(ctx *gin.Context) {
	_userId, exists := ctx.Get("userId")

	var userId int32
	var ok bool

	if exists {
		userId, ok = _userId.(int32)
	} else {
		qs := ctx.Query("user_id")
		if qs == "" {
			qs = ctx.Query("userId")
		}
		if qs != "" {
			if n, err := strconv.ParseInt(qs, 10, 32); err == nil {
				userId = int32(n)
				ok = true
			}
		}
	}

	if !ok || userId == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	err := c.emailService.ResendEmail(userId)
	if err != nil {
		ctx.Error(gin.Error{
			Err:  err,
			Type: gin.ErrorTypePublic,
		})
		return
	}
}

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

func (c *AuthController) ForgotPassword(ctx *gin.Context) {
	var forgotPasswordReq ForgotPasswordRequest
	if err := ctx.ShouldBind(&forgotPasswordReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userId, err := c.emailService.StartForgotPasswordFlow(forgotPasswordReq.Email)
	if err != nil {
		ctx.Error(gin.Error{
			Err:  err,
			Type: gin.ErrorTypePublic,
		})
		return
	}

	ctx.JSON(200, userId)
}

type UpdatePasswordRequest struct {
	Code        string `json:"code" binding:"required,min=6"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
	UserId      int32  `json:"user_id" binding:"required"`
}

func (c *AuthController) UpdatePassword(ctx *gin.Context) {
	var forgotPasswordReq UpdatePasswordRequest
	if err := ctx.ShouldBind(&forgotPasswordReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := c.emailService.UpdatePassword(forgotPasswordReq.UserId, forgotPasswordReq.NewPassword, forgotPasswordReq.Code)
	if err != nil {
		ctx.Error(gin.Error{
			Err:  err,
			Type: gin.ErrorTypePublic,
		})
		return
	}
}

package controller

import (
	"fmt"
	"gin/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService *service.AuthService
}



func NewAuthController(authService *service.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

func (c *AuthController) Login(ctx *gin.Context) {
	var loginReq LoginRequest
	if err := ctx.ShouldBind(&loginReq); err != nil {
		ctx.Error(gin.Error{
			Err:  err,
			Type: gin.ErrorTypePublic,
		})
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
		UserId: 		session.UserId,
		Token:            session.Token,
		RefreshToken:     session.RefreshToken,
		ExpiresAt:        session.ExpiresAt,
		RefreshExpiresAt: session.RefreshExpiresAt,
	})
}

func (c *AuthController) Refresh(ctx *gin.Context) {
	var refreshReq RefreshRequest
	if err := ctx.ShouldBind(&refreshReq); err != nil {
		ctx.Error(gin.Error{
			Err:  err,
			Type: gin.ErrorTypePublic,
		})
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
		UserId: 		session.UserId,
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

		userID, err := service.ValidateSession(token)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			ctx.Abort()
			return
		}

		ctx.Set("userID", userID)
		ctx.Next()
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
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if signupReq.Name == "" {
		signupReq.Name = "Test" 
	}

	session, err := c.authService.Signup(signupReq.Name, signupReq.Email, signupReq.Password)
	if err != nil {

		ctx.Error(gin.Error{
			Err:  err,
			Type: gin.ErrorTypePublic,
		})
		return
	}

	ctx.JSON(http.StatusCreated, AuthResponse{
		UserId: 		session.UserId,
		Token:            session.Token,
		RefreshToken:     session.RefreshToken,
		ExpiresAt:        session.ExpiresAt,
		RefreshExpiresAt: session.RefreshExpiresAt,
	})
}

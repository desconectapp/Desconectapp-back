package controller

import (
	"gin/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type AdminUserController struct {
	service *service.AdminUserService
}

func NewAdminUserController(conn *pgx.Conn) *AdminUserController {
	service := service.NewAdminUserService(conn)
	return &AdminUserController{
		service: service,
	}
}

func (c *AdminUserController) ListUsers(ctx *gin.Context) {
	limit, _ := strconv.Atoi(ctx.DefaultQuery("_end", "25"))
	offset, _ := strconv.Atoi(ctx.DefaultQuery("_start", "0"))
	email := ctx.Query("email")
	name := ctx.Query("name")

	emailValidated := ctx.Query("email_validated")
	var filterEmailValidated *bool
	if emailValidated != "" {
		val, err := strconv.ParseBool(emailValidated)
		if err == nil {
			filterEmailValidated = &val
		}
	}

	users, err := c.service.ListUsers(int32(limit-offset), int32(offset), &email, &name, filterEmailValidated)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

    total, err := c.service.CountUsers(&email, &name, nil)
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    ctx.Header("X-Total-Count", strconv.Itoa(int(total)))

	ctx.JSON(http.StatusOK, users)
}

func (c *AdminUserController) GetUser(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	user, err := c.service.GetUser(int32(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, user)
}

func (c *AdminUserController) CreateUser(ctx *gin.Context) {
	var req struct {
		Email     string `json:"email"`
		Password  string `json:"password"`
		Validated bool   `json:"validated"`
		Name      string `json:"name"`
		Age       int32  `json:"age"`
		City      string `json:"city"`
		Situation string `json:"situation"`
		Gender    string `json:"gender"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := c.service.CreateUserWithProfile(req.Email, req.Password, req.Validated, req.Name, req.Age, req.City, req.Situation, req.Gender)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, user)
}

func (c *AdminUserController) UpdateUser(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	var req struct {
		Email     string `json:"email"`
		Validated bool   `json:"validated"`
		Name      string `json:"name"`
		Age       int32  `json:"age"`
		City      string `json:"city"`
		Situation string `json:"situation"`
		Gender    string `json:"gender"`
		Complete  bool   `json:"complete"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := c.service.UpdateUserWithProfile(int32(id), req.Email, req.Validated, req.Name, req.Age, req.City, req.Situation, req.Gender, req.Complete)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.Status(http.StatusNoContent)
}

func (c *AdminUserController) DeleteUser(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	if err := c.service.DeleteUser(int32(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.Status(http.StatusNoContent)
}

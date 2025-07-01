package router

import (
	"context"
	"fmt"
	controller "gin/controller"
	"gin/middleware"
	"gin/service"
	"os"

	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type Router struct {
	controller           *controller.Controller
	activitiesController *controller.ActivitiesController
	authController       *controller.AuthController
	r                    *gin.Engine
}

func NewRouter() *Router {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	// defer conn.Close(context.Background())

	c := controller.NewController(conn)
	activitiesController := controller.NewActivitesController(conn)

	// Initialize auth service and controller
	authService := service.NewAuthService(conn)
	authController := controller.NewAuthController(authService)

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	r.Use(middleware.ErrorHandler())
	return &Router{
		controller:           c,
		activitiesController: activitiesController,
		authController:       authController,
		r:                    r,
	}
}

func (router *Router) Run(port string) {
	router.r.Run(port)
}

func (router *Router) SetupRoutes() {

	auth := router.r.Group("/auth")
	{
		auth.POST("/login", router.authController.Login)
		auth.POST("/logout", router.authController.Logout)
		auth.POST("/refresh", router.authController.Refresh)
	}

	users := router.r.Group("/users")
	users.Use(router.authController.AuthMiddleware())
	{
		users.GET("", router.controller.ListUsers)
		users.POST("", router.controller.CreateUser)
		users.DELETE("/:userId", router.controller.DeleteUser)
		users.GET("/:userId", router.controller.GetUser)
		users.PUT("/:userId", router.controller.UpdateUser)
	}

	activities := router.r.Group("/activities")
	activities.Use(router.authController.AuthMiddleware())
	{
		activities.GET("", router.activitiesController.ListActivitiesRequests)
		activities.POST("", router.activitiesController.CreateActivityRequest)
		activities.GET("/plain", router.activitiesController.GetActivities)
	}
}

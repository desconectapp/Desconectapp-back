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
	preferencesController *controller.PreferencesController
	groupsController		*controller.GroupsController
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
	emailValidation := service.NewEmailValidationService(conn)
	authController := controller.NewAuthController(authService, emailValidation)

	preferencesController := controller.NewPreferencesController(conn)

	groupsController := controller.NewGroupsController(conn)

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
		preferencesController: preferencesController,
		groupsController: groupsController,
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
		auth.POST("/signup", router.authController.Signup)
		auth.POST("/refresh", router.authController.Refresh)
	}

	users := router.r.Group("/users")
	users.Use(router.authController.AuthMiddleware())
	{
		users.GET("", router.controller.ListUsers)
		users.POST("/profile", router.controller.CreateProfile)
		users.DELETE("", router.controller.DeleteUser)
		users.GET("/user", router.controller.GetUser)
		// users.PUT("/:userId", router.controller.UpdateUser)
	}

	activities := router.r.Group("/activities")
	activities.Use(router.authController.AuthMiddleware())
	{
		activities.GET("/request", router.activitiesController.ListActivitiesRequests)
		activities.POST("/request", router.activitiesController.CreateActivityRequest)
		activities.GET("", router.activitiesController.GetActivities)
	}
	
	preferences := router.r.Group("/preferences")
	preferences.Use(router.authController.AuthMiddleware())
	
	{
		preferences.GET("", router.preferencesController.GetUserPreferences)
		preferences.POST("", router.preferencesController.AddPreference)
		preferences.DELETE("", router.preferencesController.DeletePreference)
		preferences.POST("/batch", router.preferencesController.BatchAddUserPreferences)
	}

	groups := router.r.Group("/groups")
	groups.Use(router.authController.AuthMiddleware())
	{
		groups.GET("/:groupId", router.groupsController.GetGroup)
		groups.GET("", router.groupsController.ListGroups)
		groups.GET("/user", router.groupsController.ListUserGroups)
		groups.POST("", router.groupsController.CreateGroup)
		groups.DELETE("/:groupId", router.groupsController.DeleteGroup)
		groups.DELETE("/user-from-group/:groupId", router.groupsController.ExitGroup)
	}
}

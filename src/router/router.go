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
	"github.com/jackc/pgx/v5/pgxpool"
)

type Router struct {
	controller              *controller.Controller
	activitiesController    *controller.ActivitiesController
	authController          *controller.AuthController
	preferencesController   *controller.PreferencesController
	groupsController        *controller.GroupsController
	adminUserController     *controller.AdminUserController
	adminActivityController *controller.AdminActivityController
	adminGroupController    *controller.AdminGroupController
	chatsController         *controller.ChatsController
	r                       *gin.Engine
}

func NewRouter() *Router {
	conn, err := pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
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
	adminUserController := controller.NewAdminUserController(conn)
	adminActivityController := controller.NewAdminActivityController(conn)
	adminGroupController := controller.NewAdminGroupController(conn)

	chatsController := controller.NewChatsController(conn)

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:              []string{"http://localhost:8081", "http://localhost:5173", "https://hoppscotch.io"},
		AllowMethods:              []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:              []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:             []string{"Content-Length", "X-Total-Count"},
		AllowCredentials:          true,
		MaxAge:                    12 * time.Hour,
		OptionsResponseStatusCode: 200,
	}))
	r.Use(middleware.ErrorHandler())
	return &Router{
		controller:              c,
		activitiesController:    activitiesController,
		authController:          authController,
		preferencesController:   preferencesController,
		groupsController:        groupsController,
		chatsController:         chatsController,
		adminUserController:     adminUserController,
		adminActivityController: adminActivityController,
		adminGroupController:    adminGroupController,
		r:                       r,
	}
}

func (router *Router) Run(port string) {
	router.r.Run(port)
}

func (router *Router) SetupRoutes() *gin.Engine {
	gin.SetMode(os.Getenv("GIN_MODE"))

	auth := router.r.Group("/auth")
	{
		auth.POST("/login", router.authController.Login)
		auth.POST("/signup", router.authController.Signup)
		auth.POST("/refresh", router.authController.Refresh)
		auth.POST("/email/verify", router.authController.ValidateEmail)
		auth.POST("/email/resend-verification", router.authController.ResendValidationEmail)
		auth.POST("/password/forgot", router.authController.ForgotPassword)
		auth.POST("/password/update", router.authController.UpdatePassword)
	}

	users := router.r.Group("/users")
	users.Use(router.authController.AuthMiddleware())
	{
		users.GET("", router.controller.ListUsers)
		users.GET("/user", router.controller.GetUser)
		users.DELETE("", router.controller.DeleteUser)
		users.POST("/profile", router.controller.CreateProfile)
		users.PUT("/profile", router.controller.UpdateProfile)
		users.PUT("/profile/avatar", router.controller.UpdateProfileAvatar)
	}

	activities := router.r.Group("/activities")
	activities.Use(router.authController.AuthMiddleware())
	{
		activities.GET("/request", router.activitiesController.ListActivitiesRequests)
		activities.POST("/request", router.activitiesController.CreateActivityRequest)
		activities.GET("", router.activitiesController.GetActivities)
		activities.DELETE("/request/:requestId", router.activitiesController.DeleteActivityRequest)
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
		groups.PUT("description/:groupId", router.groupsController.UpdateGroupDescription)
		groups.PUT("public/:groupId", router.groupsController.ChangeGroupStatus)
		groups.PUT("name/:groupId", router.groupsController.UpdateGroupName)
		groups.PUT("location/:groupId", router.groupsController.UpdateGroupLocation)
		groups.PUT("avatar/:groupId", router.groupsController.UpdateGroupAvatar)
		groups.PUT("add-user/:groupId", router.groupsController.JoinGroup)
		groups.GET("/open", router.groupsController.GetOpenGroups)
		groups.GET("/recs", router.groupsController.GetUserRecommendations)
	}
	chats := router.r.Group("/chats")
	chats.Use(router.authController.AuthMiddleware())
	{
		chats.GET("/token", router.chatsController.GetToken)
	}

	admin := router.r.Group("/admin")
	admin.Use(router.authController.AdminMiddleware())
	{
		admin.POST("/logout", router.adminUserController.LogOut)
		admin.GET("/me", router.adminUserController.GetMe)
		admin.GET("/users", router.adminUserController.ListUsers)
		admin.GET("/users/:id", router.adminUserController.GetUser)
		admin.POST("/users", router.adminUserController.CreateUser)
		admin.PUT("/users/:id", router.adminUserController.UpdateUser)
		admin.DELETE("/users/:id", router.adminUserController.DeleteUser)
		admin.POST("/users/:id/suspend", router.adminUserController.SuspendUser)
		admin.POST("/users/:id/unsuspend", router.adminUserController.UnsuspendUser)

		admin.POST("/users/password/reset", router.authController.ForgotPassword)
		admin.POST("/users/email/verify", router.authController.ResendValidationEmail)

		admin.GET("/activities", router.adminActivityController.ListActivities)
		admin.GET("/activities/:id", router.adminActivityController.GetActivity)
		admin.POST("/activities", router.adminActivityController.CreateActivity)
		admin.PUT("/activities/:id", router.adminActivityController.UpdateActivity)
		admin.DELETE("/activities/:id", router.adminActivityController.DeleteActivity)

		admin.GET("/groups", router.adminGroupController.ListGroups)
		admin.GET("/groups/:id", router.adminGroupController.GetGroup)
		admin.POST("/groups", router.adminGroupController.CreateGroup)
		admin.PUT("/groups/:id", router.adminGroupController.UpdateGroup)
		admin.DELETE("/groups/:id", router.adminGroupController.DeleteGroup)
		admin.GET("/groups/:id/members", router.adminGroupController.ListGroupMembers)
		admin.POST("/groups/:id/members", router.adminGroupController.AddGroupMember)
		admin.DELETE("/groups/:id/members/:memberId", router.adminGroupController.RemoveGroupMember)
		admin.GET("/groups/many", router.adminGroupController.GetManyGroups)
	}

	return router.r
}

package router

import (
	"context"
	"fmt"
	controller "gin/controller"
	"gin/middleware"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type Router struct {
	controller           *controller.Controller
	activitiesController *controller.ActivitiesController
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

	r := gin.Default()
	r.Use(middleware.ErrorHandler())
	return &Router{
		controller:           c,
		activitiesController: activitiesController,
		r:                    r,
	}
}

func (router *Router) Run(port string) {
	router.r.Run(port)
}

func (router *Router) SetupRoutes() {
	users := router.r.Group("/users")
	{
		users.GET("", router.controller.ListUsers)
		users.POST("", router.controller.CreateUser)
		users.DELETE("/:userId", router.controller.DeleteUser)
		users.GET("/:userId", router.controller.GetUser)
		users.PUT("/:userId", router.controller.UpdateUser)
	}

	activities := router.r.Group("/activities")
	{
		activities.GET("", router.activitiesController.ListActivities)
		activities.POST("", router.activitiesController.CreateActivity)
	}
}

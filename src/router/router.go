package router

import (
	"github.com/gin-gonic/gin"
	controller "gin/controller"
)

type Router struct {
	controller *controller.Controller
	r *gin.Engine
}


func NewRouter() *Router {
	controller := controller.NewController()
	r := gin.Default()
	return &Router{
		controller: controller,
		r:       r,
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
}
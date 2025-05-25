package main

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	repository "gin/db/generated"
	"gin/router"
)

func setupRouter(conn repository.DBTX) *gin.Engine {
	r := gin.Default()
	queries := repository.New(conn)
	ctx := context.Background()

	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	r.GET("/users", func(c *gin.Context) {
		users, err := queries.ListUsers(ctx)
		if err != nil {
			c.String(http.StatusBadRequest, "Error")
			return
		}
		log.Println(users)
		c.JSON(http.StatusOK, users)
	})

	r.POST("/users", func(c *gin.Context) {
		var userParams repository.CreateUserParams

		if err := c.ShouldBind(&userParams); err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}

		user, err := queries.CreateUser(ctx, userParams)
		if err != nil {
			log.Println(err)
			c.String(http.StatusBadRequest, "Error creating the user")
			return
		}

		c.JSON(http.StatusOK, user)
	})

	r.DELETE("/users/:userId", func(c *gin.Context) {
		userId := c.Param("userId")
		stringId, err := strconv.Atoi(userId)
		if err != nil {
			c.String(http.StatusBadRequest, "userId must be an integer")
			return
		}

		id, err := queries.DeleteUser(ctx, int32(stringId))
		if err != nil {
			log.Println(err)
			c.String(http.StatusBadRequest, "Error deleting the user")
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"deleted": id,
		})
	})

	r.PUT("/users/:userId", func(c *gin.Context) {
		userId := c.Param("userId")
		stringId, err := strconv.Atoi(userId)
		if err != nil {
			c.String(http.StatusBadRequest, "userId must be an integer")
			return
		}

		var userParams repository.UpdateUserParams
		if err := c.ShouldBind(&userParams); err != nil {
			c.String(http.StatusBadRequest, err.Error())
			return
		}

		userParams.ID = int32(stringId)

		id, err := queries.UpdateUser(ctx, userParams)
		if err != nil {
			log.Println(err)
			c.String(http.StatusBadRequest, "Error updating the user")
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"updated": id,
		})
	})

	return r
}

func main() {
	router := router.NewRouter()
	router.SetupRoutes()
	router.Run(":8080")
}

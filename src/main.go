package main

import (
	"context"
	"fmt"
	repository "gin/db/generated"

	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
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

	return r
}

func main() {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	fmt.Fprintf(os.Stdout, "Connection to database established successfully\n")

	r := setupRouter(conn)
	r.Run(":8080")
}

package main

import (
	"gin/router"
)

func main() {
	router := router.NewRouter()
	router.SetupRoutes()
	router.Run(":8080")
}

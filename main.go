package main

import (
	"os"
	"todolist/middleware"
	"todolist/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.New()
	Port := os.Getenv("PORT")

	if Port == "" {
		Port = "8080"
	}
	router.Use(gin.Logger())

	//User sign-up and log-in API
	routes.UserRoutes(router)

	//Todo routes
	routes.TodoRoutes(router)

	router.Use(middleware.CORSMiddleware())

	router.Run("0.0.0.0:" + Port)
}

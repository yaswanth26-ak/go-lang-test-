package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"

	"project/api"
	"project/config"
)

// @title Hierarchy API
// @version 1.0
// @description API for Excel upload hierarchy
// @host localhost:8080
// @BasePath /
func main() {
	cfg := config.LoadConfig()
	db := config.ConnectDatabase(cfg)

	router := gin.Default()
	router.Use(cors.Default())
	api.RegisterRoutes(router, db)

	log.Printf("Server running on port %s", cfg.AppPort)
	if err := router.Run(":" + cfg.AppPort); err != nil {
		log.Fatal("server failed to start: ", err)
	}
}
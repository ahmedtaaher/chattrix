package main

import (
	"chattrix/config"
	"chattrix/db"
	"chattrix/routes"
	"chattrix/utils"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()

	db.InitDB(cfg)
  
  jwtService := utils.NewJWTService(cfg.Auth.JWTSecret)
  
  router := gin.Default()

  router.Static("/uploads", "./uploads")

	routes.SetupRoutes(router, jwtService)

	log.Printf("Server starting on port %d", cfg.Server.Port)

	if err := router.Run(fmt.Sprintf(":%d", cfg.Server.Port)); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
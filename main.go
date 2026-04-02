package main

import (
	"chattrix/config"
	"chattrix/db"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()
	db.InitDB(cfg)
	router := gin.Default()
	log.Printf("Server starting on port %d", cfg.Server.Port)
	if err := router.Run(fmt.Sprintf(":%d", cfg.Server.Port)); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
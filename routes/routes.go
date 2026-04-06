package routes

import (
	"chattrix/db"
	"chattrix/handler"
	"chattrix/middleware"
	"chattrix/repository"
	"chattrix/service"
	"chattrix/utils"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, jwtService *utils.JWTService) {
  db := db.GetDB()
  userRepo := repository.NewUserRepository(db)
  authService := service.NewAuthService(userRepo, jwtService)
  authHandler := handler.NewAuthHandler(authService)
  
  v1 := router.Group("/chattrix/api")
  {
    auth := v1.Group("/auth")
    {
      auth.POST("/signup", authHandler.SignUp)
      auth.POST("/login", authHandler.Login)
    }

    protected := v1.Group("/")
    protected.Use(middleware.AuthMiddleware(jwtService))
    {
      protected.GET("/profile", authHandler.GetProfile)
      protected.PUT("/profile", authHandler.UpdateProfile)
      protected.POST("/change-password", authHandler.ChangePassword)
    }
  }
}
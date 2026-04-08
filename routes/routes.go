package routes

import (
	"chattrix/db"
	"chattrix/handler"
	"chattrix/middleware"
	"chattrix/repository"
	"chattrix/service"
	"chattrix/utils"
	"chattrix/websocket"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, jwtService *utils.JWTService) {
  db := db.GetDB()
  
  userRepo := repository.NewUserRepository(db)
  messageRepo := repository.NewMessageRepository(db)
  chatRepo := repository.NewChatRepository(db)
  
  hub := websocket.NewHub()
  
  authService := service.NewAuthService(userRepo, jwtService)
  messageService := service.NewMessageService(messageRepo, chatRepo, hub)

  authHandler := handler.NewAuthHandler(authService)
  wsHandler := websocket.NewWSHandler(hub, authService, jwtService, messageService)
  
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
      protected.POST("/upload-avatar", authHandler.UploadAvatar)
    }
  }

  router.GET("/ws", wsHandler.HandleConnection)
}
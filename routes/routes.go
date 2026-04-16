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
  notificationRepo := repository.NewNotificationRepository(db)
  
  hub := websocket.NewHub()
  
  authService := service.NewAuthService(userRepo, jwtService)
  notificationService := service.NewNotificationService(notificationRepo, authService, hub)
  messageService := service.NewMessageService(messageRepo, chatRepo, userRepo, notificationService, hub)
  chatService := service.NewChatService(chatRepo)
  
  authHandler := handler.NewAuthHandler(authService)
  messageHandler := handler.NewMessageHandler(messageService)
  chatHandler := handler.NewChatHandler(chatService)
  notificationHandler := handler.NewNotificationHandler(notificationService)
  
  wsHandler := websocket.NewWSHandler(hub, authService, jwtService, messageService, notificationService)
  
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
      protected.GET("/search/users", authHandler.SearchUsers)

      messages := protected.Group("/messages")
      {
	      messages.GET("/:chat_id", messageHandler.GetPaginatedMessages)
	      messages.PUT("/:id", messageHandler.EditMessage)
	      messages.DELETE("/:id", messageHandler.DeleteMessage)
        messages.GET("/unread", messageHandler.GetUnreadCounts)
      }

      chats := protected.Group("/chats")
      {
        chats.POST("", chatHandler.CreateChat)
        chats.GET("", chatHandler.GetUserChats)
        chats.POST("/:id/users", chatHandler.AddUsers)
        chats.DELETE("/:id/users/:user_id", chatHandler.RemoveUser)
        chats.DELETE("/:id/leave", chatHandler.LeaveChat)
        chats.PUT("/:id/pin", chatHandler.PinChat)
        chats.PUT("/:id/mute", chatHandler.MuteChat)
        chats.PUT("/:id/users/:user_id/role", chatHandler.ChangeUserRole)
        chats.DELETE("/:id", chatHandler.DeleteChat)
        chats.GET("/search", chatHandler.SearchChats)
        chats.POST("/:id/invite", chatHandler.CreateInvite)
      }

      notifications := protected.Group("/notifications")
      {
        notifications.GET("/notifications", notificationHandler.GetNotifications)
        notifications.PUT("/notifications/read", notificationHandler.MarkAllAsRead)
        notifications.PUT("/notifications/:id/read", notificationHandler.MarkOneAsRead)
        notifications.GET("/notifications/unread-count", notificationHandler.GetUnreadCount)
        notifications.DELETE("/:id", notificationHandler.DeleteNotification)
      }
      
      protected.POST("/join", chatHandler.JoinByInvite)
    }
  }

  router.GET("/ws", wsHandler.HandleConnection)
}
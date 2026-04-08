package websocket

import (
	"chattrix/dto"
	"chattrix/service"
	"chattrix/utils"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type WSHandler struct {
	hub         *Hub
	authService interface {
		SetOnline(userID uuid.UUID) error
		SetOffline(userID uuid.UUID) error
	}
	jwtService *utils.JWTService
	messageService *service.MessageService
}

func NewWSHandler(hub *Hub, authService interface {
	SetOnline(userID uuid.UUID) error
	SetOffline(userID uuid.UUID) error
}, jwtService *utils.JWTService, messageService *service.MessageService) *WSHandler {
	return &WSHandler{
		hub:         hub,
		authService: authService,
		jwtService:  jwtService,
		messageService: messageService,
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *WSHandler) HandleConnection(context *gin.Context) {
	token := context.Query("token")
	if token == "" {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
		return
	}

	claims, err := h.jwtService.ValidateToken(token)
	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	userID := claims.UserID

	conn, err := upgrader.Upgrade(context.Writer, context.Request, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}

	h.hub.AddUser(userID, conn)
	_ = h.authService.SetOnline(userID)
  
  unreadResponse, err := h.messageService.HandleUnreadMessages(userID)
  if err == nil {
    h.hub.SendToUsers([]uuid.UUID{userID}, unreadResponse)
  }

	log.Println("User connected:", userID)

	for {
		_, msgBytes, err := conn.ReadMessage()
		if err != nil {
			break
		}

		var wsMsg dto.WSMessage
		if err := json.Unmarshal(msgBytes, &wsMsg); err != nil {
			log.Println("Invalid message format")
			continue
		}

		switch wsMsg.Type {

		case "message":
			response, userIDs, err := h.messageService.HandleSendMessage(userID, wsMsg)
			if err != nil {
				log.Println("Message error:", err)
				continue
			}

			h.hub.SendToUsers(userIDs, response)
    
    case "seen":
      response, userIDs, err := h.messageService.HandleSeen(userID, wsMsg.ChatID)
      if err != nil {
        continue
      }
      
      h.hub.SendToUsers(userIDs, response)

    case "typing":
      response, userIDs, err := h.messageService.HandleTyping(userID, wsMsg.ChatID, true)
      if err != nil {
        continue
      }
      
      h.hub.SendToUsers(userIDs, response)

    case "stop_typing":
      response, userIDs, err := h.messageService.HandleTyping(userID, wsMsg.ChatID, false)
      if err != nil {
        continue
      }
      
      h.hub.SendToUsers(userIDs, response)
		}
	}

	h.hub.RemoveUser(userID)
	_ = h.authService.SetOffline(userID)

	log.Println("User disconnected:", userID)
}
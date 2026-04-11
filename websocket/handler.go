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
  defer conn.Close()

	h.hub.AddUser(userID, conn)
	_ = h.authService.SetOnline(userID)
  log.Println("user connected:", userID)
  
  unreadResponse, err := h.messageService.HandleUnreadMessages(userID)
  if err == nil {
    h.hub.SendToUsers([]uuid.UUID{userID}, unreadResponse)
  }

	for {
		_, msgBytes, err := conn.ReadMessage()
		if err != nil {
			break
		}

		var base struct {
      Type string `json:"type"`
    }

    if err := json.Unmarshal(msgBytes, &base); err != nil {
      log.Println("Invalid message format:", err)
      continue
    }

		switch base.Type {

		case "message":
      var wsMsg dto.WSMessage
      if err := json.Unmarshal(msgBytes, &wsMsg); err != nil {
        log.Println("Invalid message format:", err)
        continue
      }

			response, userIDs, err := h.messageService.HandleSendMessage(userID, wsMsg)
			if err != nil {
				log.Println("Message error:", err)
				continue
			}

			h.hub.SendToUsers(userIDs, response)
    
    case "seen":
      var body struct {
        ChatID uuid.UUID `json:"chat_id"`
      }
      if err := json.Unmarshal(msgBytes, &body); err != nil {
        log.Println("Invalid seen format:", err)
        continue
      }

      response, userIDs, err := h.messageService.HandleSeen(userID, body.ChatID)
      if err != nil {
        continue
      }
      
      h.hub.SendToUsers(userIDs, response)

    case "typing", "stop_typing":
      var body struct {
        ChatID uuid.UUID `json:"chat_id"`
      }
      if err := json.Unmarshal(msgBytes, &body); err != nil {
        log.Println("Invalid typing format:", err)
        continue
      }

      isTyping := base.Type == "typing"

      response, userIDs, err := h.messageService.HandleTyping(userID, body.ChatID, isTyping)
      if err != nil {
        continue
      }
      
      h.hub.SendToUsers(userIDs, response)

    case "reaction":
			var r dto.WSReaction
			if err := json.Unmarshal(msgBytes, &r); err != nil {
				continue
			}

			response, userIDs, err := h.messageService.HandleReactionRealtime(userID, r)
			if err != nil {
				continue
			}

			h.hub.SendToUsers(userIDs, response)

    case "edit":
			var body struct {
				MessageID uuid.UUID `json:"message_id"`
				Content   string    `json:"content"`
			}

			if err := json.Unmarshal(msgBytes, &body); err != nil {
				continue
			}

			response, members, err := h.messageService.EditMessageRealtime(userID, body.MessageID, body.Content)
			if err != nil {
				continue
			}

			h.hub.SendToUsers(members, response)

    case "delete":
			var body struct {
				MessageID uuid.UUID `json:"message_id"`
			}

			if err := json.Unmarshal(msgBytes, &body); err != nil {
				continue
			}

			response, members, err := h.messageService.DeleteMessageRealtime(userID, body.MessageID)
			if err != nil {
				continue
			}

			h.hub.SendToUsers(members, response)
		}
	}

	h.hub.RemoveUser(userID)
	_ = h.authService.SetOffline(userID)

	log.Println("User disconnected:", userID)
}
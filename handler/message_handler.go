package handler

import (
	"chattrix/service"
	"chattrix/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type MessageHandler struct {
	messageService *service.MessageService
}

func NewMessageHandler(messageService *service.MessageService) *MessageHandler {
  return &MessageHandler{messageService: messageService}
}

func (h *MessageHandler) GetMessages(context *gin.Context) {
	chatID, _ := uuid.Parse(context.Param("chat_id"))

	limit := 20
	offset := 0

	messages, err := h.messageService.GetMessages(chatID, limit, offset)
	if err != nil {
		utils.ErrorResponse(context, http.StatusInternalServerError, "failed to get messages")
		return
	}

	utils.SuccessResponse(context, http.StatusOK, "messages fetched", messages)
}

func (h *MessageHandler) EditMessage(context *gin.Context) {
	messageID, _ := uuid.Parse(context.Param("id"))

	var body struct {
		Content string `json:"content"`
	}

	if err := context.ShouldBindJSON(&body); err != nil {
		utils.ErrorResponse(context, http.StatusBadRequest, "invalid request")
		return
	}

	userVal, _ := context.Get("user_id")
	userID := userVal.(uuid.UUID)

	err := h.messageService.EditMessage(userID, messageID, body.Content)
	if err != nil {
		utils.ErrorResponse(context, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(context, http.StatusOK, "message updated", nil)
}

func (h *MessageHandler) DeleteMessage(context *gin.Context) {
	messageID, _ := uuid.Parse(context.Param("id"))

	userVal, _ := context.Get("user_id")
	userID := userVal.(uuid.UUID)

	err := h.messageService.DeleteMessage(userID, messageID)
	if err != nil {
		utils.ErrorResponse(context, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(context, http.StatusOK, "message deleted", nil)
}

func (h *MessageHandler) GetPaginatedMessages(context *gin.Context) {
  chatID, err := uuid.Parse(context.Param("chat_id"))
  if err != nil {
    context.JSON(400, gin.H{"error": "invalid chat id"})
      return
    }

  limit := 30

  beforeStr := context.Query("before")
  var before *time.Time

  if beforeStr != "" {
    t, err := time.Parse(time.RFC3339, beforeStr)
    if err == nil {
      before = &t
    }
  }

  msgs, err := h.messageService.GetPaginatedMessages(chatID, before, limit)
  if err != nil {
    context.JSON(500, gin.H{"error": err.Error()})
    return
  }

  context.JSON(200, msgs)
}
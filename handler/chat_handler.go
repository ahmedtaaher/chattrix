package handler

import (
	"chattrix/dto"
	"chattrix/service"
	"chattrix/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ChatHandler struct {
	chatService *service.ChatService
}

func NewChatHandler(chatService *service.ChatService) *ChatHandler {
  return &ChatHandler{chatService: chatService}
}

func (h *ChatHandler) CreateChat(context *gin.Context) {

	var chatRequest dto.CreateChatRequest

	if err := context.ShouldBindJSON(&chatRequest); err != nil {
		utils.ErrorResponse(context, http.StatusBadRequest, "invalid request")
		return
	}

	userIDStr, _ := context.Get("user_id")
	userID := userIDStr.(uuid.UUID)

	chat, err := h.chatService.CreateChat(userID, chatRequest)
	if err != nil {
		utils.ErrorResponse(context, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(context, http.StatusOK, "chat created", chat)
}

func (h *ChatHandler) GetUserChats(context *gin.Context) {

	userIDVal, _ := context.Get("user_id")
	userID := userIDVal.(uuid.UUID)

	chats, err := h.chatService.GetUserChats(userID)
	if err != nil {
		utils.ErrorResponse(context, http.StatusInternalServerError, "failed to fetch chats")
		return
	}

	utils.SuccessResponse(context, http.StatusOK, "chats fetched", chats)
}

func (h *ChatHandler) AddUsers(context *gin.Context) {
	chatID, _ := uuid.Parse(context.Param("id"))

	var req dto.AddUsersRequest
	if err := context.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(context, http.StatusBadRequest, "invalid request")
		return
	}

	userIDVal, _ := context.Get("user_id")
	requesterID := userIDVal.(uuid.UUID)

	err := h.chatService.AddUsers(requesterID, chatID, req.UserIDs)
	if err != nil {
		utils.ErrorResponse(context, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(context, http.StatusOK, "users added", nil)
}

func (h *ChatHandler) RemoveUser(context *gin.Context) {

	chatID, _ := uuid.Parse(context.Param("id"))
	userID, _ := uuid.Parse(context.Param("user_id"))

	requesterVal, _ := context.Get("user_id")
	requesterID := requesterVal.(uuid.UUID)

	err := h.chatService.RemoveUser(requesterID, chatID, userID)
	if err != nil {
		utils.ErrorResponse(context, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(context, http.StatusOK, "user removed", nil)
}

func (h *ChatHandler) LeaveChat(context *gin.Context) {
	chatID, err := uuid.Parse(context.Param("id"))
	if err != nil {
		utils.ErrorResponse(context, http.StatusBadRequest, "invalid chat id")
		return
	}

	userVal, _ := context.Get("user_id")
	userID := userVal.(uuid.UUID)

	err = h.chatService.LeaveChat(userID, chatID)
	if err != nil {
		utils.ErrorResponse(context, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(context, http.StatusOK, "left chat successfully", nil)
}

func (h *ChatHandler) PinChat(context *gin.Context) {
	chatID, _ := uuid.Parse(context.Param("id"))

	var body struct {
		IsPinned bool `json:"is_pinned"`
	}

	if err := context.ShouldBindJSON(&body); err != nil {
		utils.ErrorResponse(context, http.StatusBadRequest, "invalid request")
		return
	}

	userVal, _ := context.Get("user_id")
	userID := userVal.(uuid.UUID)

	err := h.chatService.PinChat(userID, chatID, body.IsPinned)
	if err != nil {
		utils.ErrorResponse(context, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(context, http.StatusOK, "chat updated", nil)
}

func (h *ChatHandler) MuteChat(context *gin.Context) {

	chatID, _ := uuid.Parse(context.Param("id"))

	var body struct {
		IsMuted bool `json:"is_muted"`
	}

	if err := context.ShouldBindJSON(&body); err != nil {
		utils.ErrorResponse(context, http.StatusBadRequest, "invalid request")
		return
	}

	userVal, _ := context.Get("user_id")
	userID := userVal.(uuid.UUID)

	err := h.chatService.MuteChat(userID, chatID, body.IsMuted)
	if err != nil {
		utils.ErrorResponse(context, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(context, http.StatusOK, "chat updated", nil)
}

func (h *ChatHandler) ChangeUserRole(context *gin.Context) {

	chatID, _ := uuid.Parse(context.Param("id"))
	userID, _ := uuid.Parse(context.Param("user_id"))

	var body struct {
		Role string `json:"role"`
	}

	if err := context.ShouldBindJSON(&body); err != nil {
		utils.ErrorResponse(context, http.StatusBadRequest, "invalid request")
		return
	}

	requesterVal, _ := context.Get("user_id")
	requesterID := requesterVal.(uuid.UUID)

	err := h.chatService.ChangeUserRole(requesterID, chatID, userID, body.Role)
	if err != nil {
		utils.ErrorResponse(context, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(context, http.StatusOK, "role updated", nil)
}

func (h *ChatHandler) DeleteChat(context *gin.Context) {

	chatID, err := uuid.Parse(context.Param("id"))
	if err != nil {
		utils.ErrorResponse(context, http.StatusBadRequest, "invalid chat id")
		return
	}

	userVal, _ := context.Get("user_id")
	userID := userVal.(uuid.UUID)

	err = h.chatService.DeleteChat(userID, chatID)
	if err != nil {
		utils.ErrorResponse(context, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(context, http.StatusOK, "chat deleted successfully", nil)
}

func (h *ChatHandler) SearchChats(context *gin.Context) {

	query := context.Query("q")
	if query == "" {
		utils.ErrorResponse(context, http.StatusBadRequest, "query is required")
		return
	}

	userVal, _ := context.Get("user_id")
	userID := userVal.(uuid.UUID)

	chats, err := h.chatService.SearchChats(userID, query)
	if err != nil {
		utils.ErrorResponse(context, http.StatusInternalServerError, "failed to search chats")
		return
	}

	utils.SuccessResponse(context, http.StatusOK, "search results", chats)
}

func (h *ChatHandler) CreateInvite(context *gin.Context) {
	chatID, _ := uuid.Parse(context.Param("id"))

	userVal, _ := context.Get("user_id")
	userID := userVal.(uuid.UUID)

	code, err := h.chatService.CreateInvite(userID, chatID)
	if err != nil {
		utils.ErrorResponse(context, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(context, http.StatusOK, "invite created", gin.H{
		"code": code,
	})
}

func (h *ChatHandler) JoinByInvite(context *gin.Context) {
	var body struct {
		Code string `json:"code"`
	}

	if err := context.ShouldBindJSON(&body); err != nil {
		utils.ErrorResponse(context, http.StatusBadRequest, "invalid request")
		return
	}

	userVal, _ := context.Get("user_id")
	userID := userVal.(uuid.UUID)

	err := h.chatService.JoinByInvite(userID, body.Code)
	if err != nil {
		utils.ErrorResponse(context, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(context, http.StatusOK, "joined chat", nil)
}
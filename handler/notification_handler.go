package handler

import (
	"chattrix/service"
	"chattrix/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type NotificationHandler struct {
	notificationService *service.NotificationService
}

func NewNotificationHandler(notificationService *service.NotificationService) *NotificationHandler {
	return &NotificationHandler{notificationService: notificationService}
}

func (h *NotificationHandler) GetNotifications(context *gin.Context) {
	userID := context.MustGet("user_id").(uuid.UUID)

  limit := 20
	offset := 0

	if l := context.Query("limit"); l != "" {
		limit, _ = strconv.Atoi(l)
	}

	if o := context.Query("offset"); o != "" {
		offset, _ = strconv.Atoi(o)
	}

	notifs, err := h.notificationService.GetUserNotificationsPaginated(userID, limit, offset)
	if err != nil {
		utils.ErrorResponse(context, http.StatusInternalServerError, "failed to fetch notifications")
		return
	}

	utils.SuccessResponse(context, http.StatusOK, "notifications fetched", notifs)
}

func (h *NotificationHandler) MarkAllAsRead(context *gin.Context) {
	userID := context.MustGet("user_id").(uuid.UUID)

	err := h.notificationService.MarkAllAsRead(userID)
	if err != nil {
		utils.ErrorResponse(context, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(context, http.StatusOK, "all notifications read", nil)
}

func (h *NotificationHandler) MarkOneAsRead(context *gin.Context) {
	userID := context.MustGet("user_id").(uuid.UUID)
	notifID, _ := uuid.Parse(context.Param("id"))

	err := h.notificationService.MarkOneAsRead(userID, notifID)
	if err != nil {
		utils.ErrorResponse(context, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(context, http.StatusOK, "notification read", nil)
}

func (h *NotificationHandler) GetUnreadCount(context *gin.Context) {
	userID := context.MustGet("user_id").(uuid.UUID)

	count, err := h.notificationService.GetUnreadCount(userID)
	if err != nil {
		utils.ErrorResponse(context, http.StatusBadRequest, err.Error())
		return
	}

	utils.SuccessResponse(context, http.StatusOK, "count retrieved", count)
}

func (h *NotificationHandler) DeleteNotification(context *gin.Context) {
	userID := context.MustGet("user_id").(uuid.UUID)
	notifID, _ := uuid.Parse(context.Param("id"))

	err := h.notificationService.DeleteNotification(userID, notifID)
	if err != nil {
		utils.ErrorResponse(context, http.StatusBadRequest, err.Error())
		return
	}

  utils.SuccessResponse(context, http.StatusOK, "notification deleted", nil)
}
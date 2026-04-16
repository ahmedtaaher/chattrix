package service

import (
	"chattrix/dto"
	"chattrix/mapper"
	"chattrix/models"
	"chattrix/repository"
	"encoding/json"

	"github.com/google/uuid"
)

type Notifier interface {
	SendToUsers(userIDs []uuid.UUID, message []byte)
}

type NotificationService struct {
	repo          *repository.NotificationRepository
	onlineChecker OnlineChecker
  notifier      Notifier
}

func NewNotificationService(repo *repository.NotificationRepository, onlineChecker OnlineChecker, notifier Notifier) *NotificationService {
	return &NotificationService{
		repo:          repo,
		onlineChecker: onlineChecker,
    notifier: notifier,
	}
}

func (s *NotificationService) CreateNotification(userID uuid.UUID, notifType string, refID *uuid.UUID) error {
	if s.repo.ExistsRecent(userID, notifType, refID) {
		return nil
	}
  
  notif := &models.Notification{
		UserID:      userID,
		Type:        notifType,
		ReferenceID: refID,
	}

	if err := s.repo.Create(notif); err != nil {
		return err
	}

	if s.onlineChecker.IsOnline(userID) {
		responseDTO := mapper.ToNotificationResponse(notif)

		response, _ := json.Marshal(map[string]interface{}{
			"type": "notification",
			"data": responseDTO,
		})

		s.notifier.SendToUsers([]uuid.UUID{userID}, response)
	}

  return nil
}

func (s *NotificationService) GetUserNotifications(userID uuid.UUID) ([]dto.NotificationResponse, error) {
	notifs, err := s.repo.GetUserNotifications(userID)
	if err != nil {
		return nil, err
	}

	return mapper.ToNotificationResponseList(notifs), nil
}

func (s *NotificationService) MarkAllAsRead(userID uuid.UUID) error {
	err := s.repo.MarkAllAsRead(userID)
	if err != nil {
		return err
	}

	if s.onlineChecker.IsOnline(userID) {
		response, _ := json.Marshal(map[string]interface{}{
			"type": "notifications_read_all",
		})

		s.notifier.SendToUsers([]uuid.UUID{userID}, response)
	}

	return nil
}

func (s *NotificationService) MarkOneAsRead(userID, notifID uuid.UUID) error {
	err := s.repo.MarkOneAsRead(userID, notifID)
	if err != nil {
		return err
	}

	if s.onlineChecker.IsOnline(userID) {
		response, _ := json.Marshal(map[string]interface{}{
			"type": "notification_read",
			"id":   notifID,
		})

		s.notifier.SendToUsers([]uuid.UUID{userID}, response)
	}

	return nil
}

func (s *NotificationService) GetUnreadCount(userID uuid.UUID) (int64, error) {
	return s.repo.GetUnreadCount(userID)
}

func (s *NotificationService) GetUserNotificationsPaginated(userID uuid.UUID, limit, offset int) ([]dto.NotificationResponse, error) {
	notifs, err := s.repo.GetUserNotificationsPaginated(userID, limit, offset)
	if err != nil {
		return nil, err
	}

	return mapper.ToNotificationResponseList(notifs), nil
}

func (s *NotificationService) DeleteNotification(userID, notifID uuid.UUID) error {
	return s.repo.Delete(userID, notifID)
}
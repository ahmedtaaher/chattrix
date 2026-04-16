package repository

import (
	"chattrix/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type NotificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) *NotificationRepository {
  return &NotificationRepository{db: db}
}

func (n *NotificationRepository) Create(notification *models.Notification) error {
	return n.db.Create(notification).Error
}

func (n *NotificationRepository) GetUserNotifications(userID uuid.UUID) ([]models.Notification, error) {
	var notifs []models.Notification

	err := n.db.
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(50).
		Find(&notifs).Error

	return notifs, err
}

func (n *NotificationRepository) MarkAllAsRead(userID uuid.UUID) error {
	return n.db.
		Model(&models.Notification{}).
		Where("user_id = ?", userID).
		Update("is_read", true).Error
}

func (n *NotificationRepository) MarkOneAsRead(userID, notifID uuid.UUID) error {
	return n.db.
		Model(&models.Notification{}).
		Where("id = ? AND user_id = ?", notifID, userID).
		Update("is_read", true).Error
}

func (n *NotificationRepository) GetUnreadCount(userID uuid.UUID) (int64, error) {
	var count int64

	err := n.db.
		Model(&models.Notification{}).
		Where("user_id = ? AND is_read = false", userID).
		Count(&count).Error

	return count, err
}

func (n *NotificationRepository) GetUserNotificationsPaginated(userID uuid.UUID, limit, offset int) ([]models.Notification, error) {
	var notifs []models.Notification

	err := n.db.
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&notifs).Error

	return notifs, err
}

func (n *NotificationRepository) Delete(userID, notifID uuid.UUID) error {
	return n.db.
		Where("id = ? AND user_id = ?", notifID, userID).
		Delete(&models.Notification{}).Error
}

func (n *NotificationRepository) ExistsRecent(userID uuid.UUID, notifType string,refID *uuid.UUID) bool {
	var count int64

	n.db.Model(&models.Notification{}).
		Where("user_id = ? AND type = ? AND reference_id = ? AND created_at > NOW() - INTERVAL '10 seconds'",
			userID, notifType, refID).
		Count(&count)

	return count > 0
}
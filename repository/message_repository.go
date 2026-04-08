package repository

import (
	"chattrix/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MessageRepository struct {
	db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) *MessageRepository {
  return &MessageRepository{db: db}
}

func (m *MessageRepository) CreateMessage(message *models.Message) error {
  return m.db.Create(message).Error
}

func (m *MessageRepository) CreateStatus(status *models.MessageStatus) error {
  return m.db.Create(status).Error
}

func (m *MessageRepository) UpdateStatus(messageID, userID uuid.UUID, status string) error {
  return m.db.Model(&models.MessageStatus{}).
    Where("message_id = ? AND user_id = ?", messageID, userID).
    Update("status", status).Error
}

func (m *MessageRepository) MarkMessagesAsSeen(chatID, userID uuid.UUID) error {
	return m.db.Exec(`
		UPDATE message_status ms
		SET status = 'seen', updated_at = NOW()
		FROM messages m
		WHERE ms.message_id = m.id
		AND m.chat_id = ?
		AND ms.user_id = ?
		AND ms.status != 'seen'
	`, chatID, userID).Error
}

func (m *MessageRepository) GetUnreadMessages(userID uuid.UUID) ([]models.Message, error) {
	var messages []models.Message

	err := m.db.
		Table("messages m").
		Select("m.*").
		Joins("JOIN message_status ms ON ms.message_id = m.id").
		Where("ms.user_id = ? AND ms.status != ?", userID, "seen").
		Order("m.sent_at ASC").
		Find(&messages).Error

	return messages, err
}

func (m *MessageRepository) MarkAsDelivered(userID uuid.UUID) error {
	return m.db.Exec(`
		UPDATE message_status
		SET status = 'delivered', updated_at = NOW()
		WHERE user_id = ?
		AND status = 'sent'
	`, userID).Error
}


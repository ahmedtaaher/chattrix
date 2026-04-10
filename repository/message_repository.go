package repository

import (
	"chattrix/models"
	"errors"
	"time"

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

func (m *MessageRepository) CreateAttachment(attachment *models.Attachment) error {
	return m.db.Create(attachment).Error
}

func (m *MessageRepository) GetMessageWithAttachments(messageID uuid.UUID) (*models.Message, error) {
	var msg models.Message

	err := m.db.
		Preload("Attachments").
		First(&msg, "id = ?", messageID).Error

	if err != nil {
		return nil, err
	}

	return &msg, nil
}

func (m *MessageRepository) GetMessageByID(messageID uuid.UUID) (*models.Message, error) {
	var msg models.Message

	err := m.db.First(&msg, "id = ?", messageID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("message not found")
		}
		return nil, err
	}

	return &msg, nil
}

func (m *MessageRepository) GetMessages(chatID uuid.UUID, limit, offset int) ([]models.Message, error) {
	var messages []models.Message

	err := m.db.
		Preload("Attachments").
		Where("chat_id = ?", chatID).
		Order("sent_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&messages).Error

	return messages, err
}

func (m *MessageRepository) UpdateMessageContent(messageID uuid.UUID, content string) error {
	return m.db.Model(&models.Message{}).
		Where("id = ?", messageID).
		Updates(map[string]interface{}{
			"content":   content,
			"edited_at": time.Now(),
		}).Error
}

func (m *MessageRepository) SoftDelete(messageID uuid.UUID) error {
	return m.db.Model(&models.Message{}).
		Where("id = ?", messageID).
		Update("is_deleted", true).Error
}

func (m *MessageRepository) AddReaction(reaction *models.MessageReaction) error {
	return m.db.Create(reaction).Error
}

func (m *MessageRepository) RemoveReaction(messageID, userID uuid.UUID, reaction string) error {
	return m.db.
		Where("message_id = ? AND user_id = ? AND reaction = ?", messageID, userID, reaction).
		Delete(&models.MessageReaction{}).Error
}

func (m *MessageRepository) GetMessageWithFullData(messageID uuid.UUID) (*models.Message, error) {
	var msg models.Message

	err := m.db.
		Preload("Attachments").
		Preload("Reactions").
		Preload("ReplyToMessage").
		Preload("ReplyToMessage.Attachments").
		Preload("ReplyToMessage.Reactions").
		First(&msg, "id = ?", messageID).Error

	return &msg, err
}
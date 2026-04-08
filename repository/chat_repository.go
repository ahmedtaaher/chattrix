package repository

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ChatRepository struct {
	db *gorm.DB
}

func NewChatRepository(db *gorm.DB) *ChatRepository {
  return &ChatRepository{db: db}
}

func (c *ChatRepository) GetChatMembers(chatID uuid.UUID) ([]uuid.UUID, error) {
	var userIDs []uuid.UUID

	err := c.db.Table("chat_members").
		Where("chat_id = ?", chatID).
		Pluck("user_id", &userIDs).Error

	return userIDs, err
}
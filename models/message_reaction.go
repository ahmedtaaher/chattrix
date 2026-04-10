package models

import (
	"time"

	"github.com/google/uuid"
)

type MessageReaction struct {
	MessageID uuid.UUID `gorm:"primaryKey"`
	UserID    uuid.UUID `gorm:"primaryKey"`
	Reaction  string    `gorm:"primaryKey"`
	CreatedAt time.Time
}

func (MessageReaction) TableName() string {
	return "message_reactions"
}
package models

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID               uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	ChatID           uuid.UUID  `gorm:"type:uuid;not null"`
	SenderID         uuid.UUID  `gorm:"type:uuid;not null"`
	Type             string     `gorm:"size:20;not null"`
	Content          *string    
	ReplyToMessageID *uuid.UUID 
	SentAt           time.Time  `gorm:"autoCreateTime"`
	EditedAt         *time.Time
	IsDeleted        bool `gorm:"default:false"`
  Attachments      []Attachment `gorm:"foreignKey:MessageID"`
  Reactions        []MessageReaction   `gorm:"foreignKey:MessageID"`
  ReplyToMessage   *Message `gorm:"foreignKey:ReplyToMessageID"`
}
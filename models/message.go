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
	Content          *string    `gorm:"type:text"`
	ReplyToMessageID *uuid.UUID `gorm:"type:uuid"`
	SentAt           time.Time  `gorm:"autoCreateTime"`
	EditedAt         *time.Time
	IsDeleted        bool `gorm:"default:false"`
}
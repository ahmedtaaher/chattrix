package models

import (
	"time"

	"github.com/google/uuid"
)

type MessageStatus struct {
	MessageID uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID    uuid.UUID `gorm:"type:uuid;primaryKey"`
	Status    string    `gorm:"type:varchar(20);not null"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (MessageStatus) TableName() string {
	return "message_status"
}
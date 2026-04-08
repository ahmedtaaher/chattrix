package models

import (
	"time"

	"github.com/google/uuid"
)

type MessageStatus struct {
	MessageID uuid.UUID `gorm:"type:uuid;primaryKey"`
	UserID    uuid.UUID `gorm:"type:uuid;primaryKey"`
	Status    string    `gorm:"size:20;not null"` // sent, delivered, seen
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}
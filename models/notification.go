package models

import (
	"time"

	"github.com/google/uuid"
)

type Notification struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID      uuid.UUID `gorm:"type:uuid;not null"`
	Type        string    `gorm:"type:varchar(50);not null"`
	ReferenceID *uuid.UUID
	IsRead      bool      `gorm:"default:false"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
}
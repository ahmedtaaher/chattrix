package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Username       string    `gorm:"size:50;uniqueIndex;not null"`
	Nickname       string    `gorm:"size:50;not null"`
	PasswordHash   string    `gorm:"type:text;not null"`
	AvatarURL      *string   `gorm:"type:text"`
	IsOnline       bool      `gorm:"default:false"`
  LastSeen       *time.Time
	CreatedAt      time.Time `gorm:"autoCreateTime"`
}
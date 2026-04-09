package models

import (
	"time"

	"github.com/google/uuid"
)

type ChatInvite struct {
	ID         uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	ChatID     uuid.UUID  `gorm:"type:uuid;not null"`
	InviteCode string     `gorm:"column:invite_code;size:50;unique;not null"`
	CreatedBy  *uuid.UUID `gorm:"type:uuid"`
	ExpiresAt  *time.Time
	CreatedAt  time.Time
  Chat       Chat `gorm:"foreignKey:ChatID"`
}

func (ChatInvite) TableName() string {
	return "chat_invites"
}
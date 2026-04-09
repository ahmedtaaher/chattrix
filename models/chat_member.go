package models

import (
	"time"

	"github.com/google/uuid"
)

type ChatMember struct {
	ChatID            uuid.UUID  `gorm:"type:uuid;primaryKey"`
	UserID            uuid.UUID  `gorm:"type:uuid;primaryKey"`
	Role              string     `gorm:"type:varchar(20);default:'member'"`
	JoinedAt          time.Time  `gorm:"autoCreateTime"`
	IsPinned          bool       `gorm:"default:false"`
	IsMuted           bool       `gorm:"default:false"`
	LastReadMessageID *uuid.UUID `gorm:"type:uuid"`
  User              User       `gorm:"foreignKey:UserID"`
	Chat              Chat       `gorm:"foreignKey:ChatID"`
}

func (ChatMember) TableName() string {
	return "chat_members"
}
package models

import (
	"time"

	"github.com/google/uuid"
)

type Chat struct {
	ID        uuid.UUID   `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	IsGroup   bool        `gorm:"column:is_group;not null;default:false"`
	Name      *string     `gorm:"column:name;size:50"`
	AvatarURL *string     `gorm:"column:avatar_url;type:text"`
	CreatedBy *uuid.UUID  `gorm:"column:created_by;type:uuid"`
	CreatedAt time.Time   `gorm:"column:created_at;autoCreateTime"`
	Members  []ChatMember `gorm:"foreignKey:ChatID"`
	Messages []Message    `gorm:"foreignKey:ChatID"`
}

func (Chat) TableName() string {
	return "chats"
}
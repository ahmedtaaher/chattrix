package models

import "github.com/google/uuid"

type Attachment struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	MessageID uuid.UUID `gorm:"type:uuid;not null"`
	FileURL   string    `gorm:"type:text;not null"`
	FileType  string    `gorm:"size:50"`
	FileSize  int64
	Message   Message `gorm:"foreignKey:MessageID"`
}

func (Attachment) TableName() string {
  return "attachments"
}
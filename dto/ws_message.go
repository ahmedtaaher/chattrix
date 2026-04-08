package dto

import "github.com/google/uuid"

type WSMessage struct {
	Type     string `json:"type"`
	ChatID   uuid.UUID `json:"chat_id"`
  Content  string `json:"content,omitempty"`
}
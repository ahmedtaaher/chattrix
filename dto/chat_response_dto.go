package dto

import "github.com/google/uuid"

type ChatResponse struct {
	ChatID      uuid.UUID `json:"chat_id"`
	Name        string    `json:"name"`
	IsGroup     bool      `json:"is_group"`
	LastMessage *string   `json:"last_message"`
	UnreadCount int64     `json:"unread_count"`
}
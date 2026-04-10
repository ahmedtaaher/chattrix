package dto

import "github.com/google/uuid"

type WSForward struct {
	MessageID uuid.UUID `json:"message_id"`
	ChatID    uuid.UUID `json:"chat_id"`
}
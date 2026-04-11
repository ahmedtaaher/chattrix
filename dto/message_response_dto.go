package dto

import (
	"time"

	"github.com/google/uuid"
)

type AttachmentResponse struct {
	FileURL  string `json:"file_url"`
	FileType string `json:"file_type"`
	FileSize int64  `json:"file_size"`
}

type ReactionResponse struct {
	Reaction string `json:"reaction"`
	Count    int    `json:"count"`
}

type MessageResponse struct {
	ID        uuid.UUID            `json:"id"`
	ChatID    uuid.UUID            `json:"chat_id"`
	SenderID  uuid.UUID            `json:"sender_id"`
	Type      string               `json:"type"`
	Content   *string              `json:"content"`
	SentAt    time.Time            `json:"sent_at"`
	EditedAt  *time.Time           `json:"edited_at"`
	IsDeleted bool                 `json:"is_deleted"`

	Attachments []AttachmentResponse `json:"attachments"`
	Reactions   []ReactionResponse   `json:"reactions"`

	ReplyTo *MessageResponse `json:"reply_to,omitempty"`
  ForwardFrom *MessageResponse `json:"forward_from,omitempty"`
}
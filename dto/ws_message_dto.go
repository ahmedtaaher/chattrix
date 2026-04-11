package dto

import "github.com/google/uuid"

type WSAttachment struct {
  FileURL  string `json:"file_url"`
  FileType string `json:"file_type"`
  FileSize int64  `json:"file_size"`
}

type WSMessage struct {
	Type       string `json:"type"`
	ChatID     uuid.UUID `json:"chat_id"`
  Content    string `json:"content,omitempty"`
  ReplyToID  *uuid.UUID `json:"reply_to_id,omitempty"`
  ForwardFromMessageID *uuid.UUID `json:"forward_from_message_id,omitempty"`
  Files      []WSAttachment `json:"files,omitempty"`
}

type WSReaction struct {
	Type      string    `json:"type"` 
	MessageID uuid.UUID `json:"message_id"`
	Reaction  string    `json:"reaction"`
}

type WSEdit struct {
  MessageID uuid.UUID `json:"message_id"`
  Content   string    `json:"content"`
}
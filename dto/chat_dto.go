package dto

import "github.com/google/uuid"

type CreateChatRequest struct {
	IsGroup    bool    `json:"is_group"`
	Name       *string `json:"name,omitempty"`
	UserIDs    []uuid.UUID `json:"user_ids"`
}
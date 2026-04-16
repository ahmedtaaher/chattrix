package dto

import (
	"time"

	"github.com/google/uuid"
)

type NotificationResponse struct {
	ID          uuid.UUID  `json:"id"`
	Type        string     `json:"type"`
	ReferenceID *uuid.UUID `json:"reference_id,omitempty"`
	IsRead      bool       `json:"is_read"`
	CreatedAt   time.Time  `json:"created_at"`
	Title string `json:"title"`
	Body  string `json:"body"`
}
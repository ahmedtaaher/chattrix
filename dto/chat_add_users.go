package dto

import "github.com/google/uuid"

type AddUsersRequest struct {
	UserIDs []uuid.UUID `json:"user_ids" binding:"required"`
}
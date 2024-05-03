package dto

import "github.com/google/uuid"

type UserIdInstancePermission struct {
	Instance   string    `json:"instance"`
	Permission string    `json:"permission"`
	UserId     uuid.UUID `json:"user_id"`
}

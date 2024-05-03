package dto

import (
	"github.com/google/uuid"
)

type UserIdToken struct {
	UserId uuid.UUID `json:"user_id"`
	Token  string    `json:"token"`
}

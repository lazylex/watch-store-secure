package dto

import (
	"github.com/google/uuid"
)

type SessionDTO struct {
	UserId uuid.UUID `json:"user_id"`
	Token  string    `json:"token"`
}

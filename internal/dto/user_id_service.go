package dto

import "github.com/google/uuid"

type UserIdService struct {
	UserId  uuid.UUID `json:"user_id"`
	Service string    `json:"service"`
}

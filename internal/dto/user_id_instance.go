package dto

import "github.com/google/uuid"

type UserIdInstance struct {
	UserId   uuid.UUID `json:"user_id"`
	Instance string    `json:"instance"`
}

package dto

import "github.com/google/uuid"

type InstanceNameWithUserIdDTO struct {
	UserId   uuid.UUID `json:"user_id"`
	Instance string    `json:"instance"`
}

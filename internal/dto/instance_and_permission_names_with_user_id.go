package dto

import "github.com/google/uuid"

type InstanceAndPermissionNamesWithUserIdDTO struct {
	UserId     uuid.UUID `json:"user_id"`
	Instance   string    `json:"instance"`
	Permission string    `json:"permission"`
}

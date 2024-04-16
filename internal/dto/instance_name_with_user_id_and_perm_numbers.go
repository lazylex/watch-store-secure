package dto

import "github.com/google/uuid"

type InstanceNameWithUserIdAndPermNumbersDTO struct {
	UserId            uuid.UUID `json:"user_id"`
	Instance          string    `json:"instance"`
	PermissionNumbers []int     `json:"permission_numbers"`
}

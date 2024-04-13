package dto

import "github.com/google/uuid"

type ServiceNameWithUserIdAndPermNumbersDTO struct {
	UserId            uuid.UUID `json:"user_id"`
	Service           string    `json:"service"`
	PermissionNumbers []int     `json:"permission_numbers"`
}

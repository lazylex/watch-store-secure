package dto

import "github.com/google/uuid"

type RoleServiceNamesWithUserIdDTO struct {
	UserId  uuid.UUID `json:"user_id"`
	Role    string    `json:"role"`
	Service string    `json:"service"`
}

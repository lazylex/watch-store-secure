package dto

import "github.com/google/uuid"

type GroupServiceNamesWithUserIdDTO struct {
	UserId  uuid.UUID `json:"user_id"`
	Group   string    `json:"group"`
	Service string    `json:"service"`
}

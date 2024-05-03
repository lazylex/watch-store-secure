package dto

import "github.com/google/uuid"

type UserIdGroupService struct {
	Group   string    `json:"group"`
	Service string    `json:"service"`
	UserId  uuid.UUID `json:"user_id"`
}

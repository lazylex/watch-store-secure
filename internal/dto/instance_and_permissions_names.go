package dto

import "github.com/google/uuid"

type InstanceAndPermissionsNamesDTO struct {
	UserId      uuid.UUID `json:"user_id"`
	Instance    string    `json:"instance"`
	Permissions []string  `json:"permission_numbers"`
}

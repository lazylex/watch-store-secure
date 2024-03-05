package dto

import "github.com/google/uuid"

type UserIdWithPasswordHashDTO struct {
	UserId uuid.UUID `json:"user_id"`
	Hash   string    `json:"hash"`
}

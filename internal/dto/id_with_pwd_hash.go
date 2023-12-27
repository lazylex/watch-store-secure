package dto

import "github.com/google/uuid"

type IdWithPasswordHashDTO struct {
	Id   uuid.UUID `json:"id"`
	Hash string    `json:"hash"`
}

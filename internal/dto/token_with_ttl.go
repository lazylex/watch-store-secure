package dto

import (
	"github.com/google/uuid"
	"github.com/lazylex/watch-store/secure/internal/domain/value_objects/ttl"
)

type SessionDTO struct {
	Id    uuid.UUID `json:"id"`
	Token string    `json:"token"`
	TTL   ttl.TTL   `json:"ttl"`
}

package dto

import (
	"github.com/google/uuid"
	"github.com/lazylex/watch-store/secure/internal/domain/value_objects/ttl"
)

type SessionDTO struct {
	UserId uuid.UUID `json:"user_id"`
	Token  string    `json:"token"`
	TTL    ttl.TTL   `json:"ttl"`
}

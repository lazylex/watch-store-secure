package dto

import (
	"github.com/lazylex/watch-store/secure/internal/domain/value_objects/id"
	"github.com/lazylex/watch-store/secure/internal/domain/value_objects/ttl"
)

type SessionDTO struct {
	Id    id.Id   `json:"id"`
	Token string  `json:"token"`
	TTL   ttl.TTL `json:"ttl"`
}

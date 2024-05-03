package dto

import "github.com/lazylex/watch-store/secure/internal/domain/value_objects/login"

type LoginHash struct {
	Login login.Login `json:"login"`
	Hash  string      `json:"hash"`
}

package dto

import (
	"github.com/google/uuid"
	"github.com/lazylex/watch-store/secure/internal/domain/value_objects/account_state"
	"github.com/lazylex/watch-store/secure/internal/domain/value_objects/login"
)

type UserIdLoginHashState struct {
	Login  login.Login         `json:"login"`
	UserId uuid.UUID           `json:"user_id"`
	Hash   string              `json:"hash"`
	State  account_state.State `json:"state"`
}

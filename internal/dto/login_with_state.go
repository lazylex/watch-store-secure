package dto

import (
	"github.com/lazylex/watch-store/secure/internal/domain/value_objects/account_state"
	"github.com/lazylex/watch-store/secure/internal/domain/value_objects/login"
)

type LoginStateDTO struct {
	Login login.Login         `json:"login"`
	State account_state.State `json:"state"`
}

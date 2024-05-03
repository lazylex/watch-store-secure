package dto

import (
	"github.com/lazylex/watch-store/secure/internal/domain/value_objects/login"
	"github.com/lazylex/watch-store/secure/internal/domain/value_objects/password"
)

type LoginPassword struct {
	Login    login.Login       `json:"login"`
	Password password.Password `json:"password"`
}

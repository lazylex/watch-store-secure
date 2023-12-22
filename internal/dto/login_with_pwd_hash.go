package dto

import "github.com/lazylex/watch-store/secure/internal/domain/value_objects/login"

type LoginWithPasswordHashDTO struct {
	Login login.Login `json:"login"`
	Hash  string      `json:"hash"`
}

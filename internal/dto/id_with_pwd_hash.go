package dto

import (
	"github.com/lazylex/watch-store/secure/internal/domain/value_objects/id"
)

type IdWithPasswordHashDTO struct {
	Id   id.Id  `json:"id"`
	Hash string `json:"hash"`
}

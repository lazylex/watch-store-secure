package persistent

import (
	"context"
	"github.com/lazylex/watch-store/secure/internal/domain/value_objects/login"
	"github.com/lazylex/watch-store/secure/internal/dto"
)

type LoginInterface interface {
	GetAccountLoginData(context.Context, login.Login) (dto.AccountLoginDataDTO, error)
	SetAccountLoginData(context.Context, dto.AccountLoginDataDTO) error
}

type RBACInterface interface{}

type Interface interface {
	LoginInterface
	RBACInterface
	Close()
}

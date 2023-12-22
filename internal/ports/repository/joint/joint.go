package joint

import (
	"context"
	"github.com/lazylex/watch-store/secure/internal/domain/value_objects/login"
	"github.com/lazylex/watch-store/secure/internal/dto"
)

type Interface interface {
	SaveSession(context.Context, dto.SessionDTO) error
	SaveLoginAndPasswordHash(context.Context, dto.LoginWithPasswordHashDTO) error
	GetIdAndPasswordHash(context.Context, login.Login) (dto.IdWithPasswordHashDTO, error)
}

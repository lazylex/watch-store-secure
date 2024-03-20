package in_memory

import (
	"context"
	"github.com/lazylex/watch-store/secure/internal/domain/value_objects/account_state"
	"github.com/lazylex/watch-store/secure/internal/domain/value_objects/login"
	"github.com/lazylex/watch-store/secure/internal/dto"
)

type Interface interface {
	SaveSession(context.Context, dto.SessionDTO) error
	GetUserIdAndPasswordHash(context.Context, login.Login) (dto.UserIdWithPasswordHashDTO, error)
	SetUserIdAndPasswordHash(context.Context, dto.UserLoginAndIdWithPasswordHashDTO)
	SetAccountStateByLogin(context.Context, login.Login, account_state.State)
	GetAccountStateByLogin(context.Context, login.Login) (account_state.State, error)
}

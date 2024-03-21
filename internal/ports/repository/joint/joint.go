package joint

import (
	"context"
	"github.com/lazylex/watch-store/secure/internal/domain/value_objects/account_state"
	"github.com/lazylex/watch-store/secure/internal/domain/value_objects/login"
	"github.com/lazylex/watch-store/secure/internal/dto"
)

type LoginInterface interface {
	SaveSession(context.Context, dto.SessionDTO) error
	SetAccountLoginData(context.Context, dto.AccountLoginDataDTO) error
	GetAccountLoginData(context.Context, login.Login) (dto.AccountLoginDataDTO, error)
	GetUserIdAndPasswordHash(context.Context, login.Login) (dto.UserIdWithPasswordHashDTO, error)
	SetAccountState(ctx context.Context, stateDTO dto.LoginStateDTO) error
	GetAccountState(context.Context, login.Login) (account_state.State, error)
}

type RBACInterface interface{}

type Interface interface {
	LoginInterface
	RBACInterface
}

package in_memory

import (
	"context"
	"github.com/lazylex/watch-store/secure/internal/domain/value_objects/account_state"
	"github.com/lazylex/watch-store/secure/internal/domain/value_objects/login"
	"github.com/lazylex/watch-store/secure/internal/dto"
)

type LoginInterface interface {
	SaveSession(context.Context, dto.SessionDTO) error
	SetUserIdAndPasswordHash(context.Context, dto.UserLoginAndIdWithPasswordHashDTO)
	GetUserIdAndPasswordHash(context.Context, login.Login) (dto.UserIdWithPasswordHashDTO, error)
	SetAccountState(ctx context.Context, stateDTO dto.LoginStateDTO) error
	GetAccountStateByLogin(context.Context, login.Login) (account_state.State, error)
}

type RBACInterface interface {
	SetServicePermissionsNumbersForAccount(context.Context, dto.ServiceNameWithUserIdAndPermNumbersDTO) error
	GetServicePermissionsNumbersForAccount(context.Context, dto.ServiceNameWithUserIdDTO) ([]int, error)
}

type Interface interface {
	LoginInterface
	RBACInterface
}

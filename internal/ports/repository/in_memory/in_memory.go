package in_memory

import (
	"context"
	"github.com/google/uuid"
	"github.com/lazylex/watch-store/secure/internal/domain/value_objects/account_state"
	"github.com/lazylex/watch-store/secure/internal/domain/value_objects/login"
	"github.com/lazylex/watch-store/secure/internal/dto"
)

type LoginInterface interface {
	SaveSession(context.Context, *dto.UserIdToken) error
	IsSessionActiveByUUID(context.Context, uuid.UUID) bool
	IsSessionActiveByToken(context.Context, string) bool
	DeleteSession(context.Context, uuid.UUID) error
	SetUserIdAndPasswordHash(context.Context, *dto.UserIdLoginHash)
	GetUserIdAndPasswordHash(context.Context, login.Login) (dto.UserIdHash, error)
	SetAccountState(ctx context.Context, stateDTO *dto.LoginState) error
	GetAccountStateByLogin(context.Context, login.Login) (account_state.State, error)
}

type RBACInterface interface {
	SetServicePermissionsNumbersForAccount(context.Context, *dto.UserIdServicePermNumbers) error
	GetServicePermissionsNumbersForAccount(context.Context, *dto.UserIdService) ([]int, error)

	SetInstancePermissionsNumbersForAccount(context.Context, *dto.UserIdInstancePermNumbers) error
	GetInstancePermissionsNumbersForAccount(context.Context, *dto.UserIdInstance) ([]int, error)

	ExistInstancePermissionsNumbersForAccount(context.Context, *dto.UserIdInstance) bool
	ExistServicePermissionsNumbersForAccount(context.Context, *dto.UserIdService) bool
}

type Interface interface {
	LoginInterface
	RBACInterface
}

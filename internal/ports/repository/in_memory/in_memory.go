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
	UserUUIDFromSession(ctx context.Context, sessionToken string) (uuid.UUID, error)
	SessionToken(context.Context, uuid.UUID) (string, error)
	IsSessionActiveByToken(context.Context, string) bool
	DeleteSession(context.Context, uuid.UUID) error
	SetUserIdAndPasswordHash(context.Context, *dto.UserIdLoginHash)
	UserIdAndPasswordHash(context.Context, login.Login) (dto.UserIdHash, error)
	SetAccountState(ctx context.Context, stateDTO *dto.LoginState) error
	AccountStateByLogin(context.Context, login.Login) (account_state.State, error)
}

type RBACInterface interface {
	SetServicePermissionsNumbersForAccount(context.Context, *dto.UserIdServicePermNumbers) error
	ServicePermissionsNumbersForAccount(context.Context, *dto.UserIdService) ([]int, error)

	SetInstancePermissionsNumbersForAccount(context.Context, *dto.UserIdInstancePermNumbers) error
	InstancePermissionsNumbersForAccount(context.Context, *dto.UserIdInstance) ([]int, error)

	ExistInstancePermissionsNumbersForAccount(context.Context, *dto.UserIdInstance) bool
	ExistServicePermissionsNumbersForAccount(context.Context, *dto.UserIdService) bool

	SetServiceNumberedPermissions(context.Context, string, *[]dto.NameNumber) error
	ServiceNumberedPermissions(context.Context, string) (*[]dto.NameNumber, error)
}

type InstanceInterface interface {
	SetInstanceServiceAndSecret(ctx context.Context, data *dto.NameServiceSecret) error
	SetInstanceServiceName(ctx context.Context, data *dto.NameService) error
	ServiceName(ctx context.Context, instanceName string) (string, error)
	SetInstanceSecret(ctx context.Context, data *dto.NameSecret) error
	InstanceSecret(ctx context.Context, name string) (string, error)
}

type Interface interface {
	LoginInterface
	RBACInterface
	InstanceInterface
}

package joint

import (
	"context"
	"github.com/google/uuid"
	"github.com/lazylex/watch-store/secure/internal/domain/value_objects/account_state"
	"github.com/lazylex/watch-store/secure/internal/domain/value_objects/login"
	"github.com/lazylex/watch-store/secure/internal/dto"
	"github.com/lazylex/watch-store/secure/internal/ports/common"
)

type ServiceInterface interface {
	CreateService(context.Context, *dto.NameDescription) error
	CreateOrUpdateInstance(context.Context, *dto.NameServiceSecret) error
}

type LoginInterface interface {
	SaveSession(context.Context, *dto.UserIdToken) error
	DeleteSession(context.Context, uuid.UUID) error
	GetSessionToken(context.Context, uuid.UUID) (string, error)
	SetAccountLoginData(context.Context, *dto.UserIdLoginHashState) error
	GetAccountLoginData(context.Context, login.Login) (dto.UserIdLoginHashState, error)
	GetUserIdAndPasswordHash(context.Context, login.Login) (dto.UserIdHash, error)
	GetUserUUIDFromSession(ctx context.Context, sessionToken string) (uuid.UUID, error)
	SetAccountState(context.Context, *dto.LoginState) error
	GetAccountState(context.Context, login.Login) (account_state.State, error)
}

type RBACInterface interface {
	common.RBACCreateInterface
	common.RBACAssignToAccountInterface
	common.RBACAssignInterface
	common.RBACDeleteInterface

	GetServicePermissionsForAccount(context.Context, *dto.UserIdService) ([]dto.NameNumberDescription, error)
	GetServicePermissionsNumbersForAccount(context.Context, *dto.UserIdService) ([]int, error)

	GetInstancePermissionsNumbersForAccount(context.Context, *dto.UserIdInstance) ([]int, error)
}

//go:generate mockgen -source=joint.go -destination=mocks/joint.go
type Interface interface {
	ServiceInterface
	LoginInterface
	RBACInterface
	GetInstanceSecret(context.Context, string) (string, error)
	GetServiceName(context.Context, string) (string, error)
}

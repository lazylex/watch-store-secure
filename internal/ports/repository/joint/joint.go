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
	CreateService(context.Context, *dto.NameWithDescriptionDTO) error
	CreateInstance(context.Context, *dto.NameAndServiceDTO) error
}

type LoginInterface interface {
	SaveSession(context.Context, *dto.SessionDTO) error
	DeleteSession(context.Context, uuid.UUID) error
	SetAccountLoginData(context.Context, *dto.AccountLoginDataDTO) error
	GetAccountLoginData(context.Context, login.Login) (dto.AccountLoginDataDTO, error)
	GetUserIdAndPasswordHash(context.Context, login.Login) (dto.UserIdWithPasswordHashDTO, error)
	SetAccountState(context.Context, *dto.LoginStateDTO) error
	GetAccountState(context.Context, login.Login) (account_state.State, error)
}

// TODO добавить удаление разрешений, ролей и групп

type RBACInterface interface {
	common.RBACCreateInterface
	common.RBACAssignToAccountInterface
	common.RBACAssignInterface

	GetServicePermissionsForAccount(context.Context, *dto.ServiceNameWithUserIdDTO) ([]dto.PermissionWithoutServiceDTO, error)
	GetServicePermissionsNumbersForAccount(context.Context, *dto.ServiceNameWithUserIdDTO) ([]int, error)

	GetInstancePermissionsNumbersForAccount(context.Context, *dto.InstanceNameWithUserIdDTO) ([]int, error)
}

//go:generate mockgen -source=joint.go -destination=mocks/joint.go
type Interface interface {
	ServiceInterface
	LoginInterface
	RBACInterface
}

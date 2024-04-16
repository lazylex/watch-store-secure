package joint

import (
	"context"
	"github.com/lazylex/watch-store/secure/internal/domain/value_objects/account_state"
	"github.com/lazylex/watch-store/secure/internal/domain/value_objects/login"
	"github.com/lazylex/watch-store/secure/internal/dto"
)

type ServiceInterface interface {
	AddService(context.Context, dto.NameWithDescriptionDTO) error
}

type LoginInterface interface {
	SaveSession(context.Context, dto.SessionDTO) error
	SetAccountLoginData(context.Context, dto.AccountLoginDataDTO) error
	GetAccountLoginData(context.Context, login.Login) (dto.AccountLoginDataDTO, error)
	GetUserIdAndPasswordHash(context.Context, login.Login) (dto.UserIdWithPasswordHashDTO, error)
	SetAccountState(ctx context.Context, stateDTO dto.LoginStateDTO) error
	GetAccountState(context.Context, login.Login) (account_state.State, error)
}

// TODO добавить удаление разрешений, ролей и групп

type RBACInterface interface {
	AddPermission(context.Context, dto.PermissionDTO) error
	AddRole(context.Context, dto.NameAndServiceWithDescriptionDTO) error
	AddGroup(context.Context, dto.NameAndServiceWithDescriptionDTO) error

	AssignRoleToGroup(context.Context, dto.GroupRoleServiceNamesDTO) error
	AssignRoleToAccount(context.Context, dto.RoleServiceNamesWithUserIdDTO) error
	AssignGroupToAccount(context.Context, dto.GroupServiceNamesWithUserIdDTO) error
	AssignPermissionToRole(context.Context, dto.PermissionRoleServiceNamesDTO) error
	AssignPermissionToGroup(context.Context, dto.GroupPermissionServiceNamesDTO) error

	GetServicePermissionsForAccount(context.Context, dto.ServiceNameWithUserIdDTO) ([]dto.PermissionWithoutServiceDTO, error)
	GetServicePermissionsNumbersForAccount(context.Context, dto.ServiceNameWithUserIdDTO) ([]int, error)

	GetInstancePermissionsNumbersForAccount(context.Context, dto.InstanceNameWithUserIdDTO) ([]int, error)
}

type Interface interface {
	ServiceInterface
	LoginInterface
	RBACInterface
}

package persistent

import (
	"context"
	"github.com/lazylex/watch-store/secure/internal/domain/value_objects/account_state"
	"github.com/lazylex/watch-store/secure/internal/domain/value_objects/login"
	"github.com/lazylex/watch-store/secure/internal/dto"
	"github.com/lazylex/watch-store/secure/internal/ports/repository/joint"
)

type LoginInterface interface {
	SetAccountState(context.Context, dto.LoginStateDTO) error
	GetAccountLoginData(context.Context, login.Login) (dto.AccountLoginDataDTO, error)
	SetAccountLoginData(context.Context, dto.AccountLoginDataDTO) error

	GetAccountsLoginsByState(context.Context, account_state.State) ([]login.Login, error)
}

type RBACInterface interface {
	AddPermission(context.Context, dto.PermissionDTO) error
	AddRole(context.Context, dto.NameAndServiceWithDescriptionDTO) error
	AddGroup(context.Context, dto.NameAndServiceWithDescriptionDTO) error

	AssignRoleToGroup(context.Context, dto.GroupRoleServiceNamesDTO) error
	AssignRoleToAccount(context.Context, dto.RoleServiceNamesWithUserIdDTO) error
	AssignGroupToAccount(context.Context, dto.GroupServiceNamesWithUserIdDTO) error
	AssignPermissionToRole(context.Context, dto.PermissionRoleServiceNamesDTO) error
	AssignPermissionToGroup(context.Context, dto.GroupPermissionServiceNamesDTO) error

	GetInstancePermissionsForAccount(context.Context, dto.InstanceNameWithUserIdDTO) ([]dto.PermissionWithoutServiceDTO, error)
	GetInstancePermissionsNumbersForAccount(context.Context, dto.InstanceNameWithUserIdDTO) ([]int, error)

	GetServicePermissionsForAccount(context.Context, dto.ServiceNameWithUserIdDTO) ([]dto.PermissionWithoutServiceDTO, error)
	GetServicePermissionsNumbersForAccount(context.Context, dto.ServiceNameWithUserIdDTO) ([]int, error)
}

type Interface interface {
	LoginInterface
	joint.ServiceInterface
	RBACInterface
	Close()
}

package persistent

import (
	"context"
	"github.com/lazylex/watch-store/secure/internal/domain/value_objects/account_state"
	"github.com/lazylex/watch-store/secure/internal/domain/value_objects/login"
	"github.com/lazylex/watch-store/secure/internal/dto"
	"github.com/lazylex/watch-store/secure/internal/ports/repository/joint"
)

type LoginInterface interface {
	SetAccountState(context.Context, *dto.LoginState) error
	GetAccountLoginData(context.Context, login.Login) (dto.UserIdLoginHashState, error)
	SetAccountLoginData(context.Context, *dto.UserIdLoginHashState) error

	GetAccountsLoginsByState(context.Context, account_state.State) ([]login.Login, error)
}

type RBACInterface interface {
	CreatePermission(context.Context, *dto.NameServiceDescription) error
	CreateRole(context.Context, *dto.NameServiceDescription) error
	CreateGroup(context.Context, *dto.NameServiceDescription) error

	AssignRoleToAccount(context.Context, *dto.UserIdRoleService) error
	AssignGroupToAccount(context.Context, *dto.UserIdGroupService) error
	AssignInstancePermissionToAccount(context.Context, *dto.UserIdInstancePermission) error

	AssignRoleToGroup(context.Context, *dto.GroupRoleService) error
	AssignPermissionToRole(context.Context, *dto.PermissionRoleService) error
	AssignPermissionToGroup(context.Context, *dto.GroupPermissionService) error

	GetInstancePermissionsForAccount(context.Context, *dto.UserIdInstance) ([]dto.NameNumberDescription, error)
	GetInstancePermissionsNumbersForAccount(context.Context, *dto.UserIdInstance) ([]int, error)

	GetServicePermissionsForAccount(context.Context, *dto.UserIdService) ([]dto.NameNumberDescription, error)
	GetServicePermissionsNumbersForAccount(context.Context, *dto.UserIdService) ([]int, error)

	GetPermissionNumber(ctx context.Context, permission string, instance string) (int, error)
}

type Interface interface {
	LoginInterface
	joint.ServiceInterface
	RBACInterface
	GetMaxConnections() int
	Close()
}

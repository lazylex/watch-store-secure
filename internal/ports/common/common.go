package common

import (
	"context"
	"github.com/lazylex/watch-store/secure/internal/dto"
)

type RBACCreateInterface interface {
	CreatePermission(context.Context, *dto.NameServiceDescription) error
	CreateRole(context.Context, *dto.NameServiceDescription) error
	CreateGroup(context.Context, *dto.NameServiceDescription) error
}

type RBACAssignToAccountInterface interface {
	AssignRoleToAccount(context.Context, *dto.UserIdRoleService) error
	AssignGroupToAccount(context.Context, *dto.UserIdGroupService) error
	AssignInstancePermissionToAccount(context.Context, *dto.UserIdInstancePermission) error
}

type RBACAssignInterface interface {
	AssignRoleToGroup(context.Context, *dto.GroupRoleService) error
	AssignPermissionToRole(context.Context, *dto.PermissionRoleService) error
	AssignPermissionToGroup(context.Context, *dto.GroupPermissionService) error
}

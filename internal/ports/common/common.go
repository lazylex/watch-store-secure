package common

import (
	"context"
	"github.com/lazylex/watch-store/secure/internal/dto"
)

type RBACCreateInterface interface {
	CreatePermission(context.Context, *dto.PermissionWithoutNumberDTO) error
	CreateRole(context.Context, *dto.NameAndServiceWithDescriptionDTO) error
	CreateGroup(context.Context, *dto.NameAndServiceWithDescriptionDTO) error
}

type RBACAssignToAccountInterface interface {
	AssignRoleToAccount(context.Context, *dto.RoleServiceNamesWithUserIdDTO) error
	AssignGroupToAccount(context.Context, *dto.GroupServiceNamesWithUserIdDTO) error
	AssignInstancePermissionToAccount(context.Context, *dto.InstanceAndPermissionNamesWithUserIdDTO) error
}

type RBACAssignInterface interface {
	AssignRoleToGroup(context.Context, *dto.GroupRoleServiceNamesDTO) error
	AssignPermissionToRole(context.Context, *dto.PermissionRoleServiceNamesDTO) error
	AssignPermissionToGroup(context.Context, *dto.GroupPermissionServiceNamesDTO) error
}

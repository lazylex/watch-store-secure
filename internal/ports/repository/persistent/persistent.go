package persistent

import (
	"context"
	"github.com/lazylex/watch-store/secure/internal/domain/value_objects/login"
	"github.com/lazylex/watch-store/secure/internal/dto"
)

type LoginInterface interface {
	SetAccountState(context.Context, dto.LoginStateDTO) error
	GetAccountLoginData(context.Context, login.Login) (dto.AccountLoginDataDTO, error)
	SetAccountLoginData(context.Context, dto.AccountLoginDataDTO) error
}

type ServiceInterface interface {
	AddService(context.Context, dto.NameWithDescriptionDTO) error
}

type RBACInterface interface {
	AddPermission(context.Context, dto.PermissionDTO) error
	AddRole(context.Context, dto.NameAndServiceWithDescriptionDTO) error
	AddGroup(context.Context, dto.NameAndServiceWithDescriptionDTO) error

	AssignRoleToGroup(context.Context, dto.GroupRoleServiceNamesDTO) error
	AssignPermissionToRole(context.Context, dto.PermissionRoleServiceNamesDTO) error
}

type Interface interface {
	LoginInterface
	ServiceInterface
	RBACInterface
	Close()
}

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
	AddService(ctx context.Context, descriptionDTO dto.NameWithDescriptionDTO) error
}

type RBACInterface interface {
	AddPermission(context.Context, dto.PermissionDTO) error
	AddRole(ctx context.Context, descriptionDTO dto.NameWithDescriptionDTO) error
	AddGroup(ctx context.Context, descriptionDTO dto.NameWithDescriptionDTO) error
}

type Interface interface {
	LoginInterface
	ServiceInterface
	RBACInterface
	Close()
}

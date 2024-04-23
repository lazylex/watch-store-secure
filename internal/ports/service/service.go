package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/lazylex/watch-store/secure/internal/dto"
	"github.com/lazylex/watch-store/secure/internal/service"
)

type Service interface {
	Login(context.Context, *dto.LoginPasswordDTO) (string, error)
	Logout(context.Context, uuid.UUID) error
	CreateAccount(context.Context, *dto.LoginPasswordDTO, service.AccountOptions) (uuid.UUID, error)

	CreatePermission(context.Context, dto.PermissionWithoutNumberDTO) error
	CreateRole(context.Context, dto.NameAndServiceWithDescriptionDTO) error
	CreateGroup(context.Context, dto.NameAndServiceWithDescriptionDTO) error
}

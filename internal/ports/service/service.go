package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/lazylex/watch-store/secure/internal/dto"
	"github.com/lazylex/watch-store/secure/internal/ports/common"
	"github.com/lazylex/watch-store/secure/internal/service"
)

//go:generate mockgen -source=service.go -destination=mocks/service.go
type Service interface {
	Login(context.Context, *dto.LoginPassword) (string, error)
	Logout(context.Context, uuid.UUID) error
	CreateAccount(context.Context, *dto.LoginPassword, service.AccountOptions) (uuid.UUID, error)

	RegisterInstance(context.Context, *dto.NameServiceSecret) error
	RegisterService(context.Context, *dto.NameDescription) error

	common.RBACCreateInterface
	common.RBACAssignToAccountInterface
	common.RBACAssignInterface
	common.RBACDeleteInterface

	CreateToken(context.Context, *dto.UserIdInstance) (string, error)
}

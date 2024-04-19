package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/lazylex/watch-store/secure/internal/dto"
)

type Service interface {
	Login(context.Context, *dto.LoginPasswordDTO) (dto.SessionDTO, error)
	Logout(context.Context, uuid.UUID) error
}

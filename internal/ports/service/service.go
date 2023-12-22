package service

import (
	"context"
	"github.com/lazylex/watch-store/secure/internal/dto"
)

type Service interface {
	Login(context.Context, *dto.LoginPasswordDTO) (dto.SessionDTO, error)
	// TODO добавить Logout
}

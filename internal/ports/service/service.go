package service

import (
	"github.com/lazylex/watch-store/secure/internal/dto"
)

type Service interface {
	Login(dto *dto.LoginPasswordDTO) (string, error)
	// TODO добавить Logout
}

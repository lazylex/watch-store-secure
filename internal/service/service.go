package service

import (
	"errors"
	"github.com/lazylex/watch-store/secure/internal/domain/value_objects/login"
	"github.com/lazylex/watch-store/secure/internal/dto"
)

type Service struct{}

var (
	ErrAuthenticationData = serviceError("Неправильный логин или пароль")
)

func serviceError(text string) error {
	return errors.New("service: " + text)
}

// Login совершает логин пользователя (сервиса) по переданным в dto логину и паролю. Воввращает токен сессии и ошибку
func (s *Service) Login(dto *dto.LoginPasswordDTO) (string, error) {
	if !s.isAuthenticated(dto) {
		return "", ErrAuthenticationData
	}
	token := s.createToken()
	go s.login(token, dto.Login)

	return token, nil
}

// login совершает логин пользователя, пароль которого прошел проверку. В функцию передается токен token, который будет
// действителен на время действия сессии пользователя с логином login
func (s *Service) login(token string, login login.Login) {
	// TODO implement
}

// isAuthenticated проверяет, является ли пользователь (сервис) аутентифицированным
func (s *Service) isAuthenticated(dto *dto.LoginPasswordDTO) bool {
	// TODO implement
	return true
}

// createToken создает токен для идентификации аутентифицированного пользователя (сервиса)
func (s *Service) createToken() string {
	// TODO implement
	return ""
}

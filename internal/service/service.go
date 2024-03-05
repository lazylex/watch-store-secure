package service

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/lazylex/watch-store/secure/internal/domain/value_objects/password"
	"github.com/lazylex/watch-store/secure/internal/dto"
	"github.com/lazylex/watch-store/secure/internal/ports/metrics/service"
	"github.com/lazylex/watch-store/secure/internal/repository/joint"
)

type Service struct {
	metrics    service.MetricsInterface
	repository joint.Repository
	salt       string
}

var (
	ErrAuthenticationData = serviceError("Неправильный логин или пароль")
)

func serviceError(text string) error {
	return errors.New("service: " + text)
}

// TODO заменить параметры на Option (при необходимости)
// TODO считывать salt с конфигурации

// New конструктор для сервиса
func New(metrics service.MetricsInterface, repository joint.Repository) *Service {

	return &Service{metrics: metrics, repository: repository}
}

// Login совершает логин пользователя (сервиса) по переданным в dto логину и паролю. Воввращает токен сессии и ошибку
func (s *Service) Login(dto *dto.LoginPasswordDTO) (string, error) {
	userId := s.getUserIdIfLoginAndPasswordCorrect(dto)
	if userId == uuid.Nil {
		s.metrics.AuthenticationErrorInc()
		return "", ErrAuthenticationData
	}

	token := s.createToken()

	go s.login(token, userId)

	return token, nil
}

// login совершает логин пользователя, пароль которого прошел проверку. В функцию передается токен token, который будет
// действителен на время действия сессии пользователя с идентификатором userId
func (s *Service) login(token string, userId uuid.UUID) {
	var err error

	ctx := context.Background()
	go s.metrics.LoginInc()

	// TODO сделать чтение TTL из конфигурации
	if err = s.repository.SaveSession(ctx, dto.SessionDTO{Token: token, UserId: userId, TTL: 600}); err != nil {
		return
	}
}

// getUserIdIfLoginAndPasswordCorrect возвращает uuid пользователя (сервиса), если он является аутентифицированным.
// Иначе - возвращает uuid.Nil
func (s *Service) getUserIdIfLoginAndPasswordCorrect(dto *dto.LoginPasswordDTO) uuid.UUID {
	// TODO implement
	ctx := context.Background()
	userIdAndPasswordHash, err := s.repository.GetUserIdAndPasswordHash(ctx, dto.Login)
	if err != nil {
		return uuid.Nil
	}

	if !s.isPasswordCorrect(dto.Password, userIdAndPasswordHash.Hash) {
		return uuid.Nil
	}

	return uuid.Nil
}

func (s *Service) isPasswordCorrect(password password.Password, hash string) bool {
	// TODO implement
	return true
}

// createToken создает токен для идентификации аутентифицированного пользователя (сервиса)
func (s *Service) createToken() string {
	// TODO implement
	// Возможно, стоит создавать JWT-токен, подпись которого известна только в данном сервисе. В токен записывать
	// ip адрес получателя
	return ""
}

package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"github.com/google/uuid"
	"github.com/lazylex/watch-store/secure/internal/config"
	"github.com/lazylex/watch-store/secure/internal/domain/value_objects/account_state"
	"github.com/lazylex/watch-store/secure/internal/domain/value_objects/password"
	"github.com/lazylex/watch-store/secure/internal/dto"
	"github.com/lazylex/watch-store/secure/internal/ports/metrics/service"
	"github.com/lazylex/watch-store/secure/internal/ports/repository/joint"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	metrics    service.MetricsInterface
	repository joint.Interface
	salt       string
	secure     config.Secure
}

var (
	ErrAuthenticationData = serviceError("incorrect login or password")
	ErrNotEnabledAccount  = serviceError("account is not active")
	ErrCreatePwdHash      = serviceError("error while hashing password")
	ErrCreateToken        = serviceError("error creating token")
	ErrLogout             = serviceError("error logout")
)

func serviceError(text string) error {
	return errors.New("service: " + text)
}

// New конструктор для сервиса
func New(metrics service.MetricsInterface, repository joint.Interface, cfg config.Secure) *Service {
	return &Service{metrics: metrics, repository: repository, secure: cfg}
}

// Login совершает логин пользователя (сервиса) по переданным в dto логину и паролю. Возвращает токен сессии и ошибку
func (s *Service) Login(ctx context.Context, data *dto.LoginPasswordDTO) (string, error) {
	var token string

	// TODO адаптировать ошибки приходящие из репозитория к ошибкам сервиса
	state, err := s.repository.GetAccountState(ctx, data.Login)
	if err != nil {
		return "", err
	}
	if state != account_state.Enabled {
		return "", ErrNotEnabledAccount
	}

	userId, errGetUsr := s.getUserId(ctx, data)
	if userId == uuid.Nil || errGetUsr != nil {
		s.metrics.AuthenticationErrorInc()
		return "", err
	}

	if token, err = s.createToken(); err != nil {
		return "", err
	}

	if err = s.repository.SaveSession(ctx, dto.SessionDTO{Token: token, UserId: userId}); err != nil {
		return "", err
	}

	go s.metrics.LoginInc()

	return token, nil
}

// Logout производит выход из сеанса путём удаления данных о сессии пользователя (сервиса)
func (s *Service) Logout(ctx context.Context, id uuid.UUID) error {
	if s.repository.DeleteSession(ctx, id) == nil {
		s.metrics.LogoutInc()
		return nil
	}

	return ErrLogout
}

// CreateAccount создаёт активную учетную запись
func (s *Service) CreateAccount(ctx context.Context, data *dto.LoginPasswordDTO) error {
	var hash string
	var err error

	if hash, err = s.createPasswordHash(data.Password); err != nil {
		return err
	}

	loginData := dto.AccountLoginDataDTO{
		Login:  data.Login,
		UserId: uuid.New(),
		Hash:   hash,
		State:  account_state.Enabled,
	}

	return s.repository.SetAccountLoginData(ctx, loginData)
}

// createPasswordHash создаёт хэш пароля
func (s *Service) createPasswordHash(pwd password.Password) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(pwd), s.secure.PasswordCreationCost)
	if err != nil {
		return "", ErrCreatePwdHash
	}

	return string(bytes), nil
}

// getUserId возвращает uuid пользователя (сервиса)
func (s *Service) getUserId(ctx context.Context, dto *dto.LoginPasswordDTO) (uuid.UUID, error) {
	// TODO адаптировать ошибки приходящие из репозитория к ошибкам сервиса
	userIdAndPasswordHash, err := s.repository.GetUserIdAndPasswordHash(ctx, dto.Login)
	if err != nil {
		return uuid.Nil, err
	}

	if !s.isPasswordCorrect(dto.Password, userIdAndPasswordHash.Hash) {
		return uuid.Nil, ErrAuthenticationData
	}

	return userIdAndPasswordHash.UserId, nil
}

// isPasswordCorrect возвращает true, если пароль соответствует хэшу
func (s *Service) isPasswordCorrect(pwd password.Password, hash string) bool {
	passwordHash, err := s.createPasswordHash(pwd)
	if err != nil {
		return false
	}

	return passwordHash == hash
}

// createToken создает токен сессии для идентификации аутентифицированного пользователя (сервиса)
func (s *Service) createToken() (string, error) {
	b := make([]byte, s.secure.LoginTokenLength)
	if _, err := rand.Read(b); err != nil {
		return "", ErrCreateToken
	}

	return hex.EncodeToString(b), nil
}

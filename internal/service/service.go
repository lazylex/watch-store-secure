package service

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/lazylex/watch-store/secure/internal/domain/value_objects/account_state"
	"github.com/lazylex/watch-store/secure/internal/domain/value_objects/password"
	"github.com/lazylex/watch-store/secure/internal/dto"
	"github.com/lazylex/watch-store/secure/internal/ports/metrics/service"
	"github.com/lazylex/watch-store/secure/internal/repository/joint"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
)

type Service struct {
	metrics    service.MetricsInterface
	repository joint.Repository
	salt       string
}

var (
	ErrAuthenticationData = serviceError("incorrect login or password")
	ErrNotEnabledAccount  = serviceError("account is not active")
	ErrCreatePwdHash      = serviceError("error while hashing password")
)

func serviceError(text string) error {
	return errors.New("service: " + text)
}

// TODO заменить параметры на Option (при необходимости)

// New конструктор для сервиса
func New(metrics service.MetricsInterface, repository joint.Repository) *Service {

	return &Service{metrics: metrics, repository: repository}
}

// Login совершает логин пользователя (сервиса) по переданным в dto логину и паролю. Воввращает токен сессии и ошибку
func (s *Service) Login(dto *dto.LoginPasswordDTO) (string, error) {
	// TODO адаптировать ошибки приходящие из репозитория к ошибкам сервиса
	state, err := s.repository.GetAccountState(context.Background(), dto.Login)
	if err != nil {
		return "", err
	}
	if state != account_state.Enabled {
		return "", ErrNotEnabledAccount
	}

	userId, errGetUsr := s.getUserId(dto)
	if userId == uuid.Nil || errGetUsr != nil {
		s.metrics.AuthenticationErrorInc()
		return "", err
	}

	token := s.createToken()

	go s.login(token, userId)

	return token, nil
}

// CreateAccount создаёт активную учетную запись
func (s *Service) CreateAccount(loginAndPwd *dto.LoginPasswordDTO) error {
	ctx := context.Background()
	// TODO определять стоимость создания пароля в конфиге
	bytes, err := bcrypt.GenerateFromPassword([]byte(loginAndPwd.Password), 14)
	if err != nil {
		return ErrCreatePwdHash
	}
	data := dto.AccountLoginDataDTO{
		Login:  loginAndPwd.Login,
		UserId: uuid.New(),
		Hash:   string(bytes),
		State:  account_state.Enabled,
	}

	return s.repository.SetAccountLoginData(ctx, data)
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

// getUserId возвращает uuid пользователя (сервиса)
func (s *Service) getUserId(dto *dto.LoginPasswordDTO) (uuid.UUID, error) {
	// TODO адаптировать ошибки приходящие из репозитория к ошибкам сервиса
	userIdAndPasswordHash, err := s.repository.GetUserIdAndPasswordHash(context.Background(), dto.Login)
	if err != nil {
		return uuid.Nil, err
	}

	if !s.isPasswordCorrect(dto.Password, userIdAndPasswordHash.Hash) {
		return uuid.Nil, ErrAuthenticationData
	}

	return userIdAndPasswordHash.UserId, nil
}

func (s *Service) isPasswordCorrect(password password.Password, hash string) bool {
	// TODO implement
	slog.Debug(serviceError("isPasswordCorrect not implemented").Error())
	return true
}

// createToken создает токен для идентификации аутентифицированного пользователя (сервиса)
func (s *Service) createToken() string {
	// TODO implement
	// Возможно, стоит создавать JWT-токен, подпись которого известна только в данном сервисе. В токен записывать
	// ip адрес получателя или его instance
	slog.Debug(serviceError("createToken not implemented").Error())
	return ""
}

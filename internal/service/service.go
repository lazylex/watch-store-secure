package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/google/uuid"
	"github.com/lazylex/watch-store/secure/internal/config"
	"github.com/lazylex/watch-store/secure/internal/domain/value_objects/account_state"
	"github.com/lazylex/watch-store/secure/internal/domain/value_objects/password"
	"github.com/lazylex/watch-store/secure/internal/dto"
	se "github.com/lazylex/watch-store/secure/internal/errors/service"
	"github.com/lazylex/watch-store/secure/internal/ports/metrics/service"
	"github.com/lazylex/watch-store/secure/internal/ports/repository/joint"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	metrics    service.MetricsInterface
	repository joint.Interface
	secure     config.Secure
}

type AccountOptions struct {
	Groups              []dto.NameAndServiceDTO
	Roles               []dto.NameAndServiceDTO
	InstancePermissions []dto.InstanceAndPermissionsNamesDTO
}

// New конструктор для сервиса
func New(metrics service.MetricsInterface, repository joint.Interface, cfg config.Secure) *Service {
	return &Service{metrics: metrics, repository: repository, secure: cfg}
}

// Login совершает логин пользователя (сервиса) по переданным в dto логину и паролю. Возвращает токен сессии и ошибку
func (s *Service) Login(ctx context.Context, data *dto.LoginPasswordDTO) (string, error) {
	var token string

	state, err := s.repository.GetAccountState(ctx, data.Login)
	if err != nil {
		return "", adaptErr(err)
	}
	if state != account_state.Enabled {
		return "", ErrNotEnabledAccount()
	}

	userId, errGetUsr := s.getUserId(ctx, data)
	if userId == uuid.Nil || errGetUsr != nil {
		s.metrics.AuthenticationErrorInc()
		return "", adaptErr(err)
	}

	if token, err = s.createToken(); err != nil {
		return "", adaptErr(err)
	}

	if err = s.repository.SaveSession(ctx, dto.SessionDTO{Token: token, UserId: userId}); err != nil {
		return "", adaptErr(err)
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

	return ErrLogout()
}

// CreateAccount создаёт активную учетную запись
func (s *Service) CreateAccount(ctx context.Context, data *dto.LoginPasswordDTO, options AccountOptions) (uuid.UUID, error) {
	var hash string
	var err error

	if hash, err = s.createPasswordHash(data.Password); err != nil {
		return uuid.Nil, adaptErr(err)
	}

	userId := uuid.New()

	loginData := dto.AccountLoginDataDTO{Login: data.Login, UserId: userId, Hash: hash, State: account_state.Enabled}

	if err = s.repository.SetAccountLoginData(ctx, loginData); err != nil {
		return uuid.Nil, adaptErr(err)
	}

	errAssignGroupToAccount := make(chan int)
	defer close(errAssignGroupToAccount)
	errAssignRoleToAccount := make(chan int)
	defer close(errAssignRoleToAccount)
	errAssignInstancePermToAccount := make(chan int)
	defer close(errAssignInstancePermToAccount)

	go s.assignGroupToAccount(ctx, options, userId, errAssignGroupToAccount)
	go s.assignRoleToAccount(ctx, options, userId, errAssignRoleToAccount)
	go s.assignInstancePermissionsToAccount(ctx, options, userId, errAssignInstancePermToAccount)

	errGroupCount := <-errAssignGroupToAccount
	errRoleCount := <-errAssignRoleToAccount
	errInstanceCount := <-errAssignInstancePermToAccount
	if errRoleCount == 0 && errGroupCount == 0 && errInstanceCount == 0 {
		return userId, nil
	}

	return uuid.Nil,
		adaptErr(se.NewServiceError(fmt.Sprintf("couldn’t create %d roles; %d groups; %d instance assignments;",
			errRoleCount, errGroupCount, errInstanceCount)))
}

// assignGroupToAccount привязывает группы к учетной записи
func (s *Service) assignGroupToAccount(ctx context.Context, options AccountOptions, userId uuid.UUID, c chan int) {
	var errCount int
	if options.Groups != nil && len(options.Groups) > 0 {
		for _, group := range options.Groups {
			if err := s.repository.AssignGroupToAccount(ctx, dto.GroupServiceNamesWithUserIdDTO{
				UserId:  userId,
				Group:   group.Name,
				Service: group.Service,
			}); err != nil {
				errCount++
			}
		}
	}
	c <- errCount
}

// assignRoleToAccount привязывает роли к учетной записи
func (s *Service) assignRoleToAccount(ctx context.Context, options AccountOptions, userId uuid.UUID, c chan int) {
	var errCount int
	if options.Roles != nil && len(options.Roles) > 0 {
		for _, role := range options.Roles {
			if err := s.repository.AssignRoleToAccount(ctx, dto.RoleServiceNamesWithUserIdDTO{
				UserId:  userId,
				Role:    role.Name,
				Service: role.Service,
			}); err != nil {
				errCount++
			}
		}
	}
	c <- errCount
}

// assignInstancePermissionsToAccount привязывает разрешения учетной записи к конкретному экземпляру
func (s *Service) assignInstancePermissionsToAccount(ctx context.Context, options AccountOptions, userId uuid.UUID, c chan int) {
	var errCount int
	if options.InstancePermissions != nil && len(options.InstancePermissions) > 0 {
		for _, inst := range options.InstancePermissions {
			if err := s.repository.AssignInstancePermissionToAccount(ctx, dto.InstanceAndPermissionNamesWithUserIdDTO{
				UserId:     userId,
				Instance:   inst.Instance,
				Permission: inst.Permission,
			}); err != nil {
				errCount++
			}
		}
	}
	c <- errCount
}

// createPasswordHash создаёт хэш пароля
func (s *Service) createPasswordHash(pwd password.Password) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(pwd), s.secure.PasswordCreationCost)
	if err != nil {
		return "", se.ErrCreatePwdHash
	}

	return string(bytes), nil
}

// getUserId возвращает uuid пользователя (сервиса)
func (s *Service) getUserId(ctx context.Context, dto *dto.LoginPasswordDTO) (uuid.UUID, error) {
	userIdAndPasswordHash, err := s.repository.GetUserIdAndPasswordHash(ctx, dto.Login)
	if err != nil {
		return uuid.Nil, adaptErr(err)
	}

	if !s.isPasswordCorrect(dto.Password, userIdAndPasswordHash.Hash) {
		return uuid.Nil, se.ErrAuthenticationData
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
		return "", se.ErrCreateToken
	}

	return hex.EncodeToString(b), nil
}

// CreatePermission создает разрешение
func (s *Service) CreatePermission(ctx context.Context, data dto.PermissionWithoutNumberDTO) error {
	return adaptErr(s.repository.CreatePermission(ctx, data))
}

// CreateRole создает роль
func (s *Service) CreateRole(ctx context.Context, data dto.NameAndServiceWithDescriptionDTO) error {
	return adaptErr(s.repository.CreateRole(ctx, data))
}

// CreateGroup создает группу
func (s *Service) CreateGroup(ctx context.Context, data dto.NameAndServiceWithDescriptionDTO) error {
	return adaptErr(s.repository.CreateGroup(ctx, data))
}

// RegisterInstance регистрирует название экземпляра сервиса
func (s *Service) RegisterInstance(ctx context.Context, data *dto.NameAndServiceDTO) error {
	return adaptErr(s.repository.CreateInstance(ctx, *data))
}

func (s *Service) RegisterService(ctx context.Context, data *dto.NameWithDescriptionDTO) error {
	return adaptErr(s.repository.CreateService(ctx, *data))
}

// AssignRoleToAccount прикрепляет роль к учетной записи
func (s *Service) AssignRoleToAccount(ctx context.Context, data dto.RoleServiceNamesWithUserIdDTO) error {
	return adaptErr(s.repository.AssignRoleToAccount(ctx, data))
}

// AssignGroupToAccount прикрепляет учетную запись к группе
func (s *Service) AssignGroupToAccount(ctx context.Context, data dto.GroupServiceNamesWithUserIdDTO) error {
	return adaptErr(s.repository.AssignGroupToAccount(ctx, data))
}

// AssignInstancePermissionToAccount прикрепляет к учетной записи разрешения для конкретного экземпляра сервиса
func (s *Service) AssignInstancePermissionToAccount(ctx context.Context, data dto.InstanceAndPermissionNamesWithUserIdDTO) error {
	return adaptErr(s.repository.AssignInstancePermissionToAccount(ctx, data))
}

// AssignRoleToGroup прикрепляет роль к группе
func (s *Service) AssignRoleToGroup(ctx context.Context, data dto.GroupRoleServiceNamesDTO) error {
	return adaptErr(s.repository.AssignRoleToGroup(ctx, data))
}

// AssignPermissionToRole прикрепляет разрешение к роли
func (s *Service) AssignPermissionToRole(ctx context.Context, data dto.PermissionRoleServiceNamesDTO) error {
	return adaptErr(s.repository.AssignPermissionToRole(ctx, data))
}

// AssignPermissionToGroup прикрепляет разрешение к группе
func (s *Service) AssignPermissionToGroup(ctx context.Context, data dto.GroupPermissionServiceNamesDTO) error {
	return adaptErr(s.repository.AssignPermissionToGroup(ctx, data))
}

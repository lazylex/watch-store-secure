package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/lazylex/watch-store/secure/internal/config"
	"github.com/lazylex/watch-store/secure/internal/domain/value_objects/account_state"
	"github.com/lazylex/watch-store/secure/internal/domain/value_objects/password"
	"github.com/lazylex/watch-store/secure/internal/dto"
	se "github.com/lazylex/watch-store/secure/internal/errors/service"
	"github.com/lazylex/watch-store/secure/internal/ports/metrics/service"
	"github.com/lazylex/watch-store/secure/internal/ports/repository/joint"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type Service struct {
	metrics    service.MetricsInterface
	repository joint.Interface
	secure     config.Secure
}

type AccountOptions struct {
	Groups              []dto.NameService
	Roles               []dto.NameService
	InstancePermissions []dto.InstancePermission
}

// New конструктор для сервиса
func New(metrics service.MetricsInterface, repository joint.Interface, cfg config.Secure) *Service {
	return &Service{metrics: metrics, repository: repository, secure: cfg}
}

// Login совершает логин пользователя (сервиса) по переданным в dto логину и паролю. Возвращает токен сессии и ошибку
func (s *Service) Login(ctx context.Context, data *dto.LoginPassword) (string, error) {
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
		return "", adaptErr(errGetUsr)
	}

	if token, err = s.createToken(); err != nil {
		return "", adaptErr(err)
	}

	if err = s.repository.SaveSession(ctx, &dto.UserIdToken{Token: token, UserId: userId}); err != nil {
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
func (s *Service) CreateAccount(ctx context.Context, data *dto.LoginPassword, options AccountOptions) (uuid.UUID, error) {
	var hash string
	var err error

	if hash, err = s.createPasswordHash(data.Password); err != nil {
		return uuid.Nil, adaptErr(err)
	}

	userId := uuid.New()

	loginData := dto.UserIdLoginHashState{Login: data.Login, UserId: userId, Hash: hash, State: account_state.Enabled}

	if err = s.repository.SetAccountLoginData(ctx, &loginData); err != nil {
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

	return userId,
		adaptErr(se.NewServiceError(fmt.Sprintf("couldn’t create %d roles; %d groups; %d instance assignments;",
			errRoleCount, errGroupCount, errInstanceCount)))
}

// assignGroupToAccount привязывает группы к учетной записи
func (s *Service) assignGroupToAccount(ctx context.Context, options AccountOptions, userId uuid.UUID, c chan int) {
	var errCount int
	if options.Groups != nil && len(options.Groups) > 0 {
		for _, group := range options.Groups {
			if err := s.repository.AssignGroupToAccount(ctx, &dto.UserIdGroupService{
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
			if err := s.repository.AssignRoleToAccount(ctx, &dto.UserIdRoleService{
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
			if err := s.repository.AssignInstancePermissionToAccount(ctx, &dto.UserIdInstancePermission{
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
func (s *Service) getUserId(ctx context.Context, data *dto.LoginPassword) (uuid.UUID, error) {
	userIdAndPasswordHash, err := s.repository.GetUserIdAndPasswordHash(ctx, data.Login)
	if err != nil {
		return uuid.Nil, adaptErr(err)
	}

	if !s.isPasswordCorrect(data.Password, userIdAndPasswordHash.Hash) {
		return uuid.Nil, se.ErrAuthenticationData
	}

	return userIdAndPasswordHash.UserId, nil
}

// isPasswordCorrect возвращает true, если пароль соответствует хэшу
func (s *Service) isPasswordCorrect(pwd password.Password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(pwd)) == nil
}

// createToken создает токен сессии для идентификации аутентифицированного пользователя (сервиса)
func (s *Service) createToken() (string, error) {
	b := make([]byte, s.secure.LoginTokenLength/2)
	if _, err := rand.Read(b); err != nil {
		return "", se.ErrCreateToken
	}

	return hex.EncodeToString(b), nil
}

// CreatePermission создает разрешение
func (s *Service) CreatePermission(ctx context.Context, data *dto.NameServiceDescription) error {
	return adaptErr(s.repository.CreatePermission(ctx, data))
}

// CreateRole создает роль
func (s *Service) CreateRole(ctx context.Context, data *dto.NameServiceDescription) error {
	return adaptErr(s.repository.CreateRole(ctx, data))
}

// CreateGroup создает группу
func (s *Service) CreateGroup(ctx context.Context, data *dto.NameServiceDescription) error {
	return adaptErr(s.repository.CreateGroup(ctx, data))
}

// RegisterInstance регистрирует название экземпляра сервиса
func (s *Service) RegisterInstance(ctx context.Context, data *dto.NameServiceSecret) error {
	return adaptErr(s.repository.CreateOrUpdateInstance(ctx, data))
}

func (s *Service) RegisterService(ctx context.Context, data *dto.NameDescription) error {
	return adaptErr(s.repository.CreateService(ctx, data))
}

// AssignRoleToAccount прикрепляет роль к учетной записи
func (s *Service) AssignRoleToAccount(ctx context.Context, data *dto.UserIdRoleService) error {
	return adaptErr(s.repository.AssignRoleToAccount(ctx, data))
}

// AssignGroupToAccount прикрепляет учетную запись к группе
func (s *Service) AssignGroupToAccount(ctx context.Context, data *dto.UserIdGroupService) error {
	return adaptErr(s.repository.AssignGroupToAccount(ctx, data))
}

// AssignInstancePermissionToAccount прикрепляет к учетной записи разрешения для конкретного экземпляра сервиса
func (s *Service) AssignInstancePermissionToAccount(ctx context.Context, data *dto.UserIdInstancePermission) error {
	return adaptErr(s.repository.AssignInstancePermissionToAccount(ctx, data))
}

// AssignRoleToGroup прикрепляет роль к группе
func (s *Service) AssignRoleToGroup(ctx context.Context, data *dto.GroupRoleService) error {
	return adaptErr(s.repository.AssignRoleToGroup(ctx, data))
}

// AssignPermissionToRole прикрепляет разрешение к роли
func (s *Service) AssignPermissionToRole(ctx context.Context, data *dto.PermissionRoleService) error {
	return adaptErr(s.repository.AssignPermissionToRole(ctx, data))
}

// AssignPermissionToGroup прикрепляет разрешение к группе
func (s *Service) AssignPermissionToGroup(ctx context.Context, data *dto.GroupPermissionService) error {
	return adaptErr(s.repository.AssignPermissionToGroup(ctx, data))
}

// CreateToken создает JWT-токен, содержащий номера разрешений пользователя (сервиса) для переданного экземпляра сервиса
func (s *Service) CreateToken(ctx context.Context, data *dto.UserIdInstance) (string, error) {
	var err error
	var permissions1, permissions2 []int
	var serviceName, secret string

	if secret, err = s.repository.GetInstanceSecret(ctx, data.Instance); err != nil {
		return "", adaptErr(err)
	}

	if permissions1, err = s.repository.GetInstancePermissionsNumbersForAccount(ctx, data); err != nil {
		return "", adaptErr(err)
	}

	if serviceName, err = s.repository.GetServiceName(ctx, data.Instance); err != nil {
		return "", adaptErr(err)
	}

	if permissions2, err = s.repository.GetServicePermissionsNumbersForAccount(ctx,
		&dto.UserIdService{UserId: data.UserId, Service: serviceName}); err != nil {
		return "", adaptErr(err)
	}

	permissions1 = append(permissions1, permissions2...)
	permissions2 = permissions2[:0]

	unique := make(map[int]bool, len(permissions1)/2)

	for _, v := range permissions1 {
		if _, ok := unique[v]; !ok {
			unique[v] = true
			permissions2 = append(permissions2, v)
		}
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"perm": permissions2,
		"exp":  time.Now().Add(s.secure.TokenTTL).Unix(),
	})

	return token.SignedString([]byte(secret))
}

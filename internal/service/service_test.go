package service

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/lazylex/watch-store/secure/internal/config"
	"github.com/lazylex/watch-store/secure/internal/domain/value_objects/account_state"
	"github.com/lazylex/watch-store/secure/internal/dto"
	"github.com/lazylex/watch-store/secure/internal/errors/joint"
	"github.com/lazylex/watch-store/secure/internal/errors/service"
	mockservice "github.com/lazylex/watch-store/secure/internal/ports/metrics/service/mocks"
	mockjoint "github.com/lazylex/watch-store/secure/internal/ports/repository/joint/mocks"
	"time"

	"testing"
)

var loginData = dto.LoginPassword{Login: "good", Password: "correct"}

func TestService_Login(t *testing.T) {
	ctx := context.Background()
	controller := gomock.NewController(t)
	repo := mockjoint.NewMockInterface(controller)
	metrics := mockservice.NewMockMetricsInterface(controller)
	s := New(metrics, repo, config.Secure{LoginTokenLength: 24, PasswordCreationCost: 14})

	idHash := dto.UserIdHash{UserId: uuid.New(), Hash: `$2a$14$YSZzgtT8U7a6WKLrvhyCxe4f5Cc.Gnpj/gLlIt1QrOwBGm6Uo16dm`}

	repo.EXPECT().GetAccountState(ctx, loginData.Login).Times(1).Return(account_state.State(account_state.Enabled), nil)
	repo.EXPECT().GetUserIdAndPasswordHash(ctx, loginData.Login).Times(1).Return(idHash, nil)
	repo.EXPECT().SaveSession(ctx, gomock.Any()).Times(1).Return(nil)
	metrics.EXPECT().LoginInc().AnyTimes()
	token, err := s.Login(ctx, &loginData)
	if len(token) != 24 || err != nil {
		t.Fail()
	}
}

func TestService_LoginErrGetAccountState(t *testing.T) {
	ctx := context.Background()
	controller := gomock.NewController(t)
	repo := mockjoint.NewMockInterface(controller)
	metrics := mockservice.NewMockMetricsInterface(controller)
	s := New(metrics, repo, config.Secure{LoginTokenLength: 24, PasswordCreationCost: 14})

	repo.EXPECT().GetAccountState(ctx, loginData.Login).Times(1).Return(account_state.State(0), joint.ErrEmptyResult)

	token, err := s.Login(ctx, &loginData)

	if len(token) != 0 || err == nil {
		t.Fail()
	}
}

func TestService_LoginErrNotEnabledAccount(t *testing.T) {
	ctx := context.Background()
	controller := gomock.NewController(t)
	repo := mockjoint.NewMockInterface(controller)
	metrics := mockservice.NewMockMetricsInterface(controller)
	s := New(metrics, repo, config.Secure{LoginTokenLength: 24, PasswordCreationCost: 14})

	repo.EXPECT().GetAccountState(ctx, loginData.Login).Times(1).Return(account_state.State(account_state.Disabled), nil)

	token, err := s.Login(ctx, &loginData)

	if len(token) != 0 || err != service.ErrNotEnabledAccount {
		t.Fail()
	}
}

func TestService_LoginErrGetUserId(t *testing.T) {
	ctx := context.Background()
	controller := gomock.NewController(t)
	repo := mockjoint.NewMockInterface(controller)
	metrics := mockservice.NewMockMetricsInterface(controller)
	s := New(metrics, repo, config.Secure{LoginTokenLength: 24, PasswordCreationCost: 14})

	idHash := dto.UserIdHash{UserId: uuid.Nil, Hash: `$2a$14$YSZzgtT8U7a6WKLrvhyCxe4f5Cc.Gnpj/gLlIt1QrOwBGm6Uo16dm`}

	repo.EXPECT().GetAccountState(ctx, loginData.Login).Times(1).Return(account_state.State(account_state.Enabled), nil)
	repo.EXPECT().GetUserIdAndPasswordHash(ctx, loginData.Login).Times(1).Return(idHash, nil)

	metrics.EXPECT().AuthenticationErrorInc().Times(1)
	token, err := s.Login(ctx, &loginData)
	if len(token) != 0 || err != nil {
		t.Fail()
	}
}

func TestService_LoginErrDataNotSaved(t *testing.T) {
	ctx := context.Background()
	controller := gomock.NewController(t)
	repo := mockjoint.NewMockInterface(controller)
	metrics := mockservice.NewMockMetricsInterface(controller)
	s := New(metrics, repo, config.Secure{LoginTokenLength: 24, PasswordCreationCost: 14})

	idHash := dto.UserIdHash{UserId: uuid.New(), Hash: `$2a$14$YSZzgtT8U7a6WKLrvhyCxe4f5Cc.Gnpj/gLlIt1QrOwBGm6Uo16dm`}

	repo.EXPECT().GetAccountState(ctx, loginData.Login).Times(1).Return(account_state.State(account_state.Enabled), nil)
	repo.EXPECT().GetUserIdAndPasswordHash(ctx, loginData.Login).Times(1).Return(idHash, nil)
	repo.EXPECT().SaveSession(ctx, gomock.Any()).Times(1).Return(joint.ErrDataNotSaved)

	token, err := s.Login(ctx, &loginData)
	if len(token) != 0 || err == nil {
		t.Fail()
	}
}

func TestService_LoginErrGetUserIdAndPasswordHash(t *testing.T) {
	ctx := context.Background()
	controller := gomock.NewController(t)
	repo := mockjoint.NewMockInterface(controller)
	metrics := mockservice.NewMockMetricsInterface(controller)
	s := New(metrics, repo, config.Secure{LoginTokenLength: 24, PasswordCreationCost: 14})

	repo.EXPECT().GetAccountState(ctx, loginData.Login).Times(1).Return(account_state.State(account_state.Enabled), nil)
	repo.EXPECT().GetUserIdAndPasswordHash(ctx, loginData.Login).Times(1).Return(dto.UserIdHash{}, joint.ErrEmptyResult)
	metrics.EXPECT().AuthenticationErrorInc().Times(1)

	token, err := s.Login(ctx, &loginData)
	if len(token) != 0 || err == nil {
		t.Fail()
	}
}

func TestService_LoginIncorrectPassword(t *testing.T) {
	ctx := context.Background()
	controller := gomock.NewController(t)
	repo := mockjoint.NewMockInterface(controller)
	metrics := mockservice.NewMockMetricsInterface(controller)
	s := New(metrics, repo, config.Secure{LoginTokenLength: 24, PasswordCreationCost: 14})

	repo.EXPECT().GetAccountState(ctx, loginData.Login).Times(1).Return(account_state.State(account_state.Enabled), nil)
	repo.EXPECT().GetUserIdAndPasswordHash(ctx, loginData.Login).Times(1).Return(dto.UserIdHash{Hash: "incorrect pwd"}, nil)
	metrics.EXPECT().AuthenticationErrorInc().Times(1)

	token, err := s.Login(ctx, &loginData)
	if len(token) != 0 || err == nil {
		t.Fail()
	}
}

func TestService_Logout(t *testing.T) {
	ctx := context.Background()
	controller := gomock.NewController(t)
	repo := mockjoint.NewMockInterface(controller)
	metrics := mockservice.NewMockMetricsInterface(controller)
	s := New(metrics, repo, config.Secure{LoginTokenLength: 24, PasswordCreationCost: 14})
	id := uuid.New()
	repo.EXPECT().DeleteSession(ctx, id).Times(1).Return(nil)
	metrics.EXPECT().LogoutInc().Times(1)

	if s.Logout(ctx, id) != nil {
		t.Fail()
	}
}

func TestService_LogoutError(t *testing.T) {
	ctx := context.Background()
	controller := gomock.NewController(t)
	repo := mockjoint.NewMockInterface(controller)
	metrics := mockservice.NewMockMetricsInterface(controller)
	s := New(metrics, repo, config.Secure{LoginTokenLength: 24, PasswordCreationCost: 14})
	id := uuid.New()
	repo.EXPECT().DeleteSession(ctx, id).Times(1).Return(errors.New(""))

	if s.Logout(ctx, id) != service.ErrLogout {
		t.Fail()
	}
}

func TestService_CreateAccount(t *testing.T) {
	ctx := context.Background()
	controller := gomock.NewController(t)
	repo := mockjoint.NewMockInterface(controller)
	metrics := mockservice.NewMockMetricsInterface(controller)
	s := New(metrics, repo, config.Secure{LoginTokenLength: 24, PasswordCreationCost: 14})

	repo.EXPECT().SetAccountLoginData(ctx, gomock.Any()).Times(1).Return(nil)
	repo.EXPECT().AssignGroupToAccount(ctx, gomock.Any()).Times(1).Return(nil)
	repo.EXPECT().AssignRoleToAccount(ctx, gomock.Any()).Times(1).Return(nil)
	repo.EXPECT().AssignInstancePermissionToAccount(ctx, gomock.Any()).Times(1).Return(nil)

	accountId, err := s.CreateAccount(ctx, &dto.LoginPassword{Login: "Homer Jay Simpson", Password: "donut"}, AccountOptions{
		Groups:              []dto.NameService{{"users", "tron"}},
		Roles:               []dto.NameService{{"admin", "tron"}},
		InstancePermissions: []dto.InstancePermission{{"node1", "delete"}},
	})
	if err != nil || accountId == uuid.Nil {
		t.Fail()
	}
}

func TestService_CreateAccountErrAssign(t *testing.T) {
	ctx := context.Background()
	controller := gomock.NewController(t)
	repo := mockjoint.NewMockInterface(controller)
	metrics := mockservice.NewMockMetricsInterface(controller)
	s := New(metrics, repo, config.Secure{LoginTokenLength: 24, PasswordCreationCost: 14})

	repo.EXPECT().SetAccountLoginData(ctx, gomock.Any()).Times(1).Return(nil)
	repo.EXPECT().AssignGroupToAccount(ctx, gomock.Any()).Times(1).Return(joint.ErrDataNotSaved)
	repo.EXPECT().AssignRoleToAccount(ctx, gomock.Any()).Times(1).Return(joint.ErrDataNotSaved)
	repo.EXPECT().AssignInstancePermissionToAccount(ctx, gomock.Any()).Times(1).Return(joint.ErrDataNotSaved)

	accountId, err := s.CreateAccount(ctx, &dto.LoginPassword{Login: "Homer Jay Simpson", Password: "donut"}, AccountOptions{
		Groups:              []dto.NameService{{"users", "tron"}},
		Roles:               []dto.NameService{{"admin", "tron"}},
		InstancePermissions: []dto.InstancePermission{{"node1", "delete"}},
	})
	if err == nil || accountId == uuid.Nil {
		t.Fail()
	}
}

func TestService_CreateAccountErrSetAccountData(t *testing.T) {
	ctx := context.Background()
	controller := gomock.NewController(t)
	repo := mockjoint.NewMockInterface(controller)
	metrics := mockservice.NewMockMetricsInterface(controller)
	s := New(metrics, repo, config.Secure{LoginTokenLength: 24, PasswordCreationCost: 14})

	repo.EXPECT().SetAccountLoginData(ctx, gomock.Any()).Times(1).Return(joint.ErrDataNotSaved)

	accountId, err := s.CreateAccount(ctx, &dto.LoginPassword{Login: "Homer Jay Simpson", Password: "donut"}, AccountOptions{
		Groups:              []dto.NameService{{"users", "tron"}},
		Roles:               []dto.NameService{{"admin", "tron"}},
		InstancePermissions: []dto.InstancePermission{{"node1", "delete"}},
	})
	if err == nil || accountId != uuid.Nil {
		t.Fail()
	}
}

func TestService_RegisterInstance(t *testing.T) {
	ctx := context.Background()
	controller := gomock.NewController(t)
	repo := mockjoint.NewMockInterface(controller)
	metrics := mockservice.NewMockMetricsInterface(controller)
	s := New(metrics, repo, config.Secure{LoginTokenLength: 24, PasswordCreationCost: 14})
	data := dto.NameServiceSecret{Name: "saver", Service: "tron", Secret: "secret"}

	repo.EXPECT().CreateOrUpdateInstance(ctx, &data).Times(1).Return(nil)

	if s.RegisterInstance(ctx, &data) != nil {
		t.Fail()
	}
}

func TestService_RegisterInstanceErr(t *testing.T) {
	ctx := context.Background()
	controller := gomock.NewController(t)
	repo := mockjoint.NewMockInterface(controller)
	metrics := mockservice.NewMockMetricsInterface(controller)
	s := New(metrics, repo, config.Secure{LoginTokenLength: 24, PasswordCreationCost: 14})
	data := dto.NameServiceSecret{Name: "saver", Service: "tron", Secret: "secret"}

	repo.EXPECT().CreateOrUpdateInstance(ctx, &data).Times(1).Return(joint.ErrDuplicateData)

	if s.RegisterInstance(ctx, &data) != service.ErrAlreadyExist {
		t.Fail()
	}
}

func TestService_RegisterService(t *testing.T) {
	ctx := context.Background()
	controller := gomock.NewController(t)
	repo := mockjoint.NewMockInterface(controller)
	metrics := mockservice.NewMockMetricsInterface(controller)
	s := New(metrics, repo, config.Secure{LoginTokenLength: 24, PasswordCreationCost: 14})
	data := dto.NameDescription{Name: "saver", Description: ""}

	repo.EXPECT().CreateService(ctx, &data).Times(1).Return(nil)

	if s.RegisterService(ctx, &data) != nil {
		t.Fail()
	}
}

func TestService_RegisterServiceErr(t *testing.T) {
	ctx := context.Background()
	controller := gomock.NewController(t)
	repo := mockjoint.NewMockInterface(controller)
	metrics := mockservice.NewMockMetricsInterface(controller)
	s := New(metrics, repo, config.Secure{LoginTokenLength: 24, PasswordCreationCost: 14})
	data := dto.NameDescription{Name: "saver", Description: ""}

	repo.EXPECT().CreateService(ctx, &data).Times(1).Return(joint.ErrDuplicateData)

	if s.RegisterService(ctx, &data) != service.ErrAlreadyExist {
		t.Fail()
	}
}

func TestService_CreatePermission(t *testing.T) {
	ctx := context.Background()
	controller := gomock.NewController(t)
	repo := mockjoint.NewMockInterface(controller)
	metrics := mockservice.NewMockMetricsInterface(controller)
	s := New(metrics, repo, config.Secure{LoginTokenLength: 24, PasswordCreationCost: 14})
	data := dto.NameServiceDescription{Name: "Flynn", Description: "", Service: "tron"}

	repo.EXPECT().CreatePermission(ctx, &data).Times(1).Return(nil)

	if s.CreatePermission(ctx, &data) != nil {
		t.Fail()
	}
}

func TestService_CreatePermissionErr(t *testing.T) {
	ctx := context.Background()
	controller := gomock.NewController(t)
	repo := mockjoint.NewMockInterface(controller)
	metrics := mockservice.NewMockMetricsInterface(controller)
	s := New(metrics, repo, config.Secure{LoginTokenLength: 24, PasswordCreationCost: 14})
	data := dto.NameServiceDescription{Name: "Flynn", Description: "", Service: "tron"}

	repo.EXPECT().CreatePermission(ctx, &data).Times(1).Return(joint.ErrDuplicateData)

	if s.CreatePermission(ctx, &data) != service.ErrAlreadyExist {
		t.Fail()
	}
}

func TestService_CreateRole(t *testing.T) {
	ctx := context.Background()
	controller := gomock.NewController(t)
	repo := mockjoint.NewMockInterface(controller)
	metrics := mockservice.NewMockMetricsInterface(controller)
	s := New(metrics, repo, config.Secure{LoginTokenLength: 24, PasswordCreationCost: 14})
	data := dto.NameServiceDescription{Name: "creator", Description: "", Service: "tron"}

	repo.EXPECT().CreateRole(ctx, &data).Times(1).Return(nil)

	if s.CreateRole(ctx, &data) != nil {
		t.Fail()
	}
}

func TestService_CreateRoleErr(t *testing.T) {
	ctx := context.Background()
	controller := gomock.NewController(t)
	repo := mockjoint.NewMockInterface(controller)
	metrics := mockservice.NewMockMetricsInterface(controller)
	s := New(metrics, repo, config.Secure{LoginTokenLength: 24, PasswordCreationCost: 14})
	data := dto.NameServiceDescription{Name: "creator", Description: "", Service: "tron"}

	repo.EXPECT().CreateRole(ctx, &data).Times(1).Return(joint.ErrDuplicateData)

	if s.CreateRole(ctx, &data) != service.ErrAlreadyExist {
		t.Fail()
	}
}

func TestService_CreateGroup(t *testing.T) {
	ctx := context.Background()
	controller := gomock.NewController(t)
	repo := mockjoint.NewMockInterface(controller)
	metrics := mockservice.NewMockMetricsInterface(controller)
	s := New(metrics, repo, config.Secure{LoginTokenLength: 24, PasswordCreationCost: 14})
	data := dto.NameServiceDescription{Name: "users", Description: "", Service: "tron"}

	repo.EXPECT().CreateGroup(ctx, &data).Times(1).Return(nil)

	if s.CreateGroup(ctx, &data) != nil {
		t.Fail()
	}
}

func TestService_CreateGroupErr(t *testing.T) {
	ctx := context.Background()
	controller := gomock.NewController(t)
	repo := mockjoint.NewMockInterface(controller)
	metrics := mockservice.NewMockMetricsInterface(controller)
	s := New(metrics, repo, config.Secure{LoginTokenLength: 24, PasswordCreationCost: 14})
	data := dto.NameServiceDescription{Name: "users", Description: "", Service: "tron"}

	repo.EXPECT().CreateGroup(ctx, &data).Times(1).Return(joint.ErrDuplicateData)

	if s.CreateGroup(ctx, &data) != service.ErrAlreadyExist {
		t.Fail()
	}
}

func TestService_AssignRoleToAccount(t *testing.T) {
	ctx := context.Background()
	controller := gomock.NewController(t)
	repo := mockjoint.NewMockInterface(controller)
	metrics := mockservice.NewMockMetricsInterface(controller)
	s := New(metrics, repo, config.Secure{LoginTokenLength: 24, PasswordCreationCost: 14})
	data := dto.UserIdRoleService{UserId: uuid.New(), Role: "visitor", Service: "tron"}

	repo.EXPECT().AssignRoleToAccount(ctx, &data).Times(1).Return(nil)

	if s.AssignRoleToAccount(ctx, &data) != nil {
		t.Fail()
	}
}

func TestService_AssignRoleToAccountErr(t *testing.T) {
	ctx := context.Background()
	controller := gomock.NewController(t)
	repo := mockjoint.NewMockInterface(controller)
	metrics := mockservice.NewMockMetricsInterface(controller)
	s := New(metrics, repo, config.Secure{LoginTokenLength: 24, PasswordCreationCost: 14})
	data := dto.UserIdRoleService{UserId: uuid.New(), Role: "visitor", Service: "tron"}

	repo.EXPECT().AssignRoleToAccount(ctx, &data).Times(1).Return(joint.ErrDataNotSaved)

	if s.AssignRoleToAccount(ctx, &data) == nil {
		t.Fail()
	}
}

func TestService_AssignGroupToAccount(t *testing.T) {
	ctx := context.Background()
	controller := gomock.NewController(t)
	repo := mockjoint.NewMockInterface(controller)
	metrics := mockservice.NewMockMetricsInterface(controller)
	s := New(metrics, repo, config.Secure{LoginTokenLength: 24, PasswordCreationCost: 14})
	data := dto.UserIdGroupService{UserId: uuid.New(), Group: "users", Service: "tron"}

	repo.EXPECT().AssignGroupToAccount(ctx, &data).Times(1).Return(nil)

	if s.AssignGroupToAccount(ctx, &data) != nil {
		t.Fail()
	}
}

func TestService_AssignGroupToAccountErr(t *testing.T) {
	ctx := context.Background()
	controller := gomock.NewController(t)
	repo := mockjoint.NewMockInterface(controller)
	metrics := mockservice.NewMockMetricsInterface(controller)
	s := New(metrics, repo, config.Secure{LoginTokenLength: 24, PasswordCreationCost: 14})
	data := dto.UserIdGroupService{UserId: uuid.New(), Group: "users", Service: "tron"}

	repo.EXPECT().AssignGroupToAccount(ctx, &data).Times(1).Return(joint.ErrDataNotSaved)

	if s.AssignGroupToAccount(ctx, &data) == nil {
		t.Fail()
	}
}

func TestService_AssignInstancePermissionToAccount(t *testing.T) {
	ctx := context.Background()
	controller := gomock.NewController(t)
	repo := mockjoint.NewMockInterface(controller)
	metrics := mockservice.NewMockMetricsInterface(controller)
	s := New(metrics, repo, config.Secure{LoginTokenLength: 24, PasswordCreationCost: 14})
	data := dto.UserIdInstancePermission{Instance: "node1", Permission: "delete", UserId: uuid.New()}

	repo.EXPECT().AssignInstancePermissionToAccount(ctx, &data).Times(1).Return(nil)

	if s.AssignInstancePermissionToAccount(ctx, &data) != nil {
		t.Fail()
	}
}

func TestService_AssignInstancePermissionToAccountErr(t *testing.T) {
	ctx := context.Background()
	controller := gomock.NewController(t)
	repo := mockjoint.NewMockInterface(controller)
	metrics := mockservice.NewMockMetricsInterface(controller)
	s := New(metrics, repo, config.Secure{LoginTokenLength: 24, PasswordCreationCost: 14})
	data := dto.UserIdInstancePermission{Instance: "node1", Permission: "delete", UserId: uuid.New()}

	repo.EXPECT().AssignInstancePermissionToAccount(ctx, &data).Times(1).Return(joint.ErrDataNotSaved)

	if s.AssignInstancePermissionToAccount(ctx, &data) == nil {
		t.Fail()
	}
}

func TestService_AssignRoleToGroup(t *testing.T) {
	ctx := context.Background()
	controller := gomock.NewController(t)
	repo := mockjoint.NewMockInterface(controller)
	metrics := mockservice.NewMockMetricsInterface(controller)
	s := New(metrics, repo, config.Secure{LoginTokenLength: 24, PasswordCreationCost: 14})
	data := dto.GroupRoleService{Group: "users", Role: "admin", Service: "tron"}

	repo.EXPECT().AssignRoleToGroup(ctx, &data).Times(1).Return(nil)

	if s.AssignRoleToGroup(ctx, &data) != nil {
		t.Fail()
	}
}

func TestService_AssignRoleToGroupErr(t *testing.T) {
	ctx := context.Background()
	controller := gomock.NewController(t)
	repo := mockjoint.NewMockInterface(controller)
	metrics := mockservice.NewMockMetricsInterface(controller)
	s := New(metrics, repo, config.Secure{LoginTokenLength: 24, PasswordCreationCost: 14})
	data := dto.GroupRoleService{Group: "users", Role: "admin", Service: "tron"}

	repo.EXPECT().AssignRoleToGroup(ctx, &data).Times(1).Return(joint.ErrDataNotSaved)

	if s.AssignRoleToGroup(ctx, &data) == nil {
		t.Fail()
	}
}

func TestService_AssignPermissionToRole(t *testing.T) {
	ctx := context.Background()
	controller := gomock.NewController(t)
	repo := mockjoint.NewMockInterface(controller)
	metrics := mockservice.NewMockMetricsInterface(controller)
	s := New(metrics, repo, config.Secure{LoginTokenLength: 24, PasswordCreationCost: 14})
	data := dto.PermissionRoleService{Permission: "delete", Role: "admin", Service: "tron"}

	repo.EXPECT().AssignPermissionToRole(ctx, &data).Times(1).Return(nil)

	if s.AssignPermissionToRole(ctx, &data) != nil {
		t.Fail()
	}
}

func TestService_AssignPermissionToRoleErr(t *testing.T) {
	ctx := context.Background()
	controller := gomock.NewController(t)
	repo := mockjoint.NewMockInterface(controller)
	metrics := mockservice.NewMockMetricsInterface(controller)
	s := New(metrics, repo, config.Secure{LoginTokenLength: 24, PasswordCreationCost: 14})
	data := dto.PermissionRoleService{Permission: "delete", Role: "admin", Service: "tron"}

	repo.EXPECT().AssignPermissionToRole(ctx, &data).Times(1).Return(joint.ErrDataNotSaved)

	if s.AssignPermissionToRole(ctx, &data) == nil {
		t.Fail()
	}
}

func TestService_AssignPermissionToGroup(t *testing.T) {
	ctx := context.Background()
	controller := gomock.NewController(t)
	repo := mockjoint.NewMockInterface(controller)
	metrics := mockservice.NewMockMetricsInterface(controller)
	s := New(metrics, repo, config.Secure{LoginTokenLength: 24, PasswordCreationCost: 14})
	data := dto.GroupPermissionService{Permission: "delete", Group: "users", Service: "tron"}

	repo.EXPECT().AssignPermissionToGroup(ctx, &data).Times(1).Return(nil)

	if s.AssignPermissionToGroup(ctx, &data) != nil {
		t.Fail()
	}
}

func TestService_AssignPermissionToGroupErr(t *testing.T) {
	ctx := context.Background()
	controller := gomock.NewController(t)
	repo := mockjoint.NewMockInterface(controller)
	metrics := mockservice.NewMockMetricsInterface(controller)
	s := New(metrics, repo, config.Secure{LoginTokenLength: 24, PasswordCreationCost: 14})
	data := dto.GroupPermissionService{Permission: "delete", Group: "users", Service: "tron"}

	repo.EXPECT().AssignPermissionToGroup(ctx, &data).Times(1).Return(joint.ErrDataNotSaved)

	if s.AssignPermissionToGroup(ctx, &data) == nil {
		t.Fail()
	}
}

func TestService_CreateToken(t *testing.T) {
	ctx := context.Background()
	controller := gomock.NewController(t)
	repo := mockjoint.NewMockInterface(controller)
	metrics := mockservice.NewMockMetricsInterface(controller)
	s := New(metrics, repo, config.Secure{LoginTokenLength: 24, PasswordCreationCost: 14, TokenTTL: 168 * time.Hour})

	repo.EXPECT().GetInstanceSecret(ctx, gomock.Any()).Times(1).Return("secret", nil)
	repo.EXPECT().GetInstancePermissionsNumbersForAccount(ctx, gomock.Any()).Times(1).Return([]int{1}, nil)
	repo.EXPECT().GetServiceName(ctx, gomock.Any()).Times(1).Return("", nil)
	repo.EXPECT().GetServicePermissionsNumbersForAccount(ctx, gomock.Any()).Times(1).Return([]int{4, 6}, nil)
	token, err := s.CreateToken(ctx, &dto.UserIdInstance{UserId: uuid.Nil, Instance: ""})
	if len(token) == 0 || err != nil {
		t.Fail()
	}
}

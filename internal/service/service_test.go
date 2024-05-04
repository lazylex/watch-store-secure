package service

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/lazylex/watch-store/secure/internal/config"
	"github.com/lazylex/watch-store/secure/internal/domain/value_objects/account_state"
	"github.com/lazylex/watch-store/secure/internal/dto"
	"github.com/lazylex/watch-store/secure/internal/errors/joint"
	"github.com/lazylex/watch-store/secure/internal/errors/service"
	mockservice "github.com/lazylex/watch-store/secure/internal/ports/metrics/service/mocks"
	mockjoint "github.com/lazylex/watch-store/secure/internal/ports/repository/joint/mocks"

	"testing"
)

var correctLoginData = dto.LoginPassword{Login: "good", Password: "correct"}

func TestService_Login(t *testing.T) {
	ctx := context.Background()
	controller := gomock.NewController(t)
	repo := mockjoint.NewMockInterface(controller)
	metrics := mockservice.NewMockMetricsInterface(controller)
	s := New(metrics, repo, config.Secure{LoginTokenLength: 24, PasswordCreationCost: 14})

	idHash := dto.UserIdHash{UserId: uuid.New(), Hash: `$2a$14$YSZzgtT8U7a6WKLrvhyCxe4f5Cc.Gnpj/gLlIt1QrOwBGm6Uo16dm`}

	repo.EXPECT().GetAccountState(ctx, correctLoginData.Login).Times(1).Return(account_state.State(account_state.Enabled), nil)
	repo.EXPECT().GetUserIdAndPasswordHash(ctx, correctLoginData.Login).Times(1).Return(idHash, nil)
	repo.EXPECT().SaveSession(ctx, gomock.Any()).Times(1).Return(nil)
	metrics.EXPECT().LoginInc().AnyTimes()
	token, err := s.Login(ctx, &correctLoginData)
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

	repo.EXPECT().GetAccountState(ctx, correctLoginData.Login).Times(1).Return(account_state.State(0), joint.ErrEmptyResult)

	token, err := s.Login(ctx, &correctLoginData)

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

	repo.EXPECT().GetAccountState(ctx, correctLoginData.Login).Times(1).Return(account_state.State(account_state.Disabled), nil)

	token, err := s.Login(ctx, &correctLoginData)

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

	repo.EXPECT().GetAccountState(ctx, correctLoginData.Login).Times(1).Return(account_state.State(account_state.Enabled), nil)
	repo.EXPECT().GetUserIdAndPasswordHash(ctx, correctLoginData.Login).Times(1).Return(idHash, nil)

	metrics.EXPECT().AuthenticationErrorInc().Times(1)
	token, err := s.Login(ctx, &correctLoginData)
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

	repo.EXPECT().GetAccountState(ctx, correctLoginData.Login).Times(1).Return(account_state.State(account_state.Enabled), nil)
	repo.EXPECT().GetUserIdAndPasswordHash(ctx, correctLoginData.Login).Times(1).Return(idHash, nil)
	repo.EXPECT().SaveSession(ctx, gomock.Any()).Times(1).Return(joint.ErrDataNotSaved)

	token, err := s.Login(ctx, &correctLoginData)
	if len(token) != 0 || err == nil {
		t.Fail()
	}
}

func TestService_Logout(t *testing.T) {
	t.Fail()
}

func TestService_CreateAccount(t *testing.T) {
	t.Fail()
}

func TestService_RegisterInstance(t *testing.T) {
	t.Fail()
}

func TestService_RegisterService(t *testing.T) {
	t.Fail()
}

func TestService_CreatePermission(t *testing.T) {
	t.Fail()
}

func TestService_CreateRole(t *testing.T) {
	t.Fail()
}

func TestService_CreateGroup(t *testing.T) {
	t.Fail()
}

func TestService_AssignRoleToAccount(t *testing.T) {
	t.Fail()
}

func TestService_AssignGroupToAccount(t *testing.T) {
	t.Fail()
}

func TestService_AssignInstancePermissionToAccount(t *testing.T) {
	t.Fail()
}

func TestService_AssignRoleToGroup(t *testing.T) {
	t.Fail()
}

func TestService_AssignPermissionToRole(t *testing.T) {
	t.Fail()
}

func TestService_AssignPermissionToGroup(t *testing.T) {
	t.Fail()
}

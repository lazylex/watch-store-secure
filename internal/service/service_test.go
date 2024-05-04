package service

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/lazylex/watch-store/secure/internal/config"
	"github.com/lazylex/watch-store/secure/internal/domain/value_objects/account_state"
	"github.com/lazylex/watch-store/secure/internal/dto"
	mock_service "github.com/lazylex/watch-store/secure/internal/ports/metrics/service/mocks"
	mock_joint "github.com/lazylex/watch-store/secure/internal/ports/repository/joint/mocks"

	"testing"
)

func TestService_Login(t *testing.T) {
	ctx := context.Background()
	controller := gomock.NewController(t)
	repo := mock_joint.NewMockInterface(controller)
	metrics := mock_service.NewMockMetricsInterface(controller)
	s := New(metrics, repo, config.Secure{LoginTokenLength: 24, PasswordCreationCost: 14})

	data := dto.LoginPassword{Login: "good", Password: "correct"}
	idHash := dto.UserIdHash{UserId: uuid.New(), Hash: `$2a$14$YSZzgtT8U7a6WKLrvhyCxe4f5Cc.Gnpj/gLlIt1QrOwBGm6Uo16dm`}

	repo.EXPECT().GetAccountState(ctx, data.Login).Times(1).Return(account_state.State(account_state.Enabled), nil)
	repo.EXPECT().GetUserIdAndPasswordHash(ctx, data.Login).Times(1).Return(idHash, nil)
	repo.EXPECT().SaveSession(ctx, gomock.Any()).Times(1).Return(nil)
	metrics.EXPECT().LoginInc().AnyTimes()
	token, err := s.Login(ctx, &data)
	if len(token) != 24 || err != nil {
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

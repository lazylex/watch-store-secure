// Code generated by MockGen. DO NOT EDIT.
// Source: service.go

// Package mock_service is a generated GoMock package.
package mock_service

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
	dto "github.com/lazylex/watch-store/secure/internal/dto"
	service "github.com/lazylex/watch-store/secure/internal/service"
)

// MockService is a mock of Service interface.
type MockService struct {
	ctrl     *gomock.Controller
	recorder *MockServiceMockRecorder
}

// MockServiceMockRecorder is the mock recorder for MockService.
type MockServiceMockRecorder struct {
	mock *MockService
}

// NewMockService creates a new mock instance.
func NewMockService(ctrl *gomock.Controller) *MockService {
	mock := &MockService{ctrl: ctrl}
	mock.recorder = &MockServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockService) EXPECT() *MockServiceMockRecorder {
	return m.recorder
}

// AssignGroupToAccount mocks base method.
func (m *MockService) AssignGroupToAccount(arg0 context.Context, arg1 *dto.UserIdGroupService) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AssignGroupToAccount", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// AssignGroupToAccount indicates an expected call of AssignGroupToAccount.
func (mr *MockServiceMockRecorder) AssignGroupToAccount(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AssignGroupToAccount", reflect.TypeOf((*MockService)(nil).AssignGroupToAccount), arg0, arg1)
}

// AssignInstancePermissionToAccount mocks base method.
func (m *MockService) AssignInstancePermissionToAccount(arg0 context.Context, arg1 *dto.UserIdInstancePermission) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AssignInstancePermissionToAccount", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// AssignInstancePermissionToAccount indicates an expected call of AssignInstancePermissionToAccount.
func (mr *MockServiceMockRecorder) AssignInstancePermissionToAccount(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AssignInstancePermissionToAccount", reflect.TypeOf((*MockService)(nil).AssignInstancePermissionToAccount), arg0, arg1)
}

// AssignPermissionToGroup mocks base method.
func (m *MockService) AssignPermissionToGroup(arg0 context.Context, arg1 *dto.GroupPermissionService) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AssignPermissionToGroup", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// AssignPermissionToGroup indicates an expected call of AssignPermissionToGroup.
func (mr *MockServiceMockRecorder) AssignPermissionToGroup(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AssignPermissionToGroup", reflect.TypeOf((*MockService)(nil).AssignPermissionToGroup), arg0, arg1)
}

// AssignPermissionToRole mocks base method.
func (m *MockService) AssignPermissionToRole(arg0 context.Context, arg1 *dto.PermissionRoleService) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AssignPermissionToRole", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// AssignPermissionToRole indicates an expected call of AssignPermissionToRole.
func (mr *MockServiceMockRecorder) AssignPermissionToRole(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AssignPermissionToRole", reflect.TypeOf((*MockService)(nil).AssignPermissionToRole), arg0, arg1)
}

// AssignRoleToAccount mocks base method.
func (m *MockService) AssignRoleToAccount(arg0 context.Context, arg1 *dto.UserIdRoleService) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AssignRoleToAccount", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// AssignRoleToAccount indicates an expected call of AssignRoleToAccount.
func (mr *MockServiceMockRecorder) AssignRoleToAccount(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AssignRoleToAccount", reflect.TypeOf((*MockService)(nil).AssignRoleToAccount), arg0, arg1)
}

// AssignRoleToGroup mocks base method.
func (m *MockService) AssignRoleToGroup(arg0 context.Context, arg1 *dto.GroupRoleService) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AssignRoleToGroup", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// AssignRoleToGroup indicates an expected call of AssignRoleToGroup.
func (mr *MockServiceMockRecorder) AssignRoleToGroup(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AssignRoleToGroup", reflect.TypeOf((*MockService)(nil).AssignRoleToGroup), arg0, arg1)
}

// CreateAccount mocks base method.
func (m *MockService) CreateAccount(arg0 context.Context, arg1 *dto.LoginPassword, arg2 service.AccountOptions) (uuid.UUID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateAccount", arg0, arg1, arg2)
	ret0, _ := ret[0].(uuid.UUID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateAccount indicates an expected call of CreateAccount.
func (mr *MockServiceMockRecorder) CreateAccount(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateAccount", reflect.TypeOf((*MockService)(nil).CreateAccount), arg0, arg1, arg2)
}

// CreateGroup mocks base method.
func (m *MockService) CreateGroup(arg0 context.Context, arg1 *dto.NameServiceDescription) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateGroup", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateGroup indicates an expected call of CreateGroup.
func (mr *MockServiceMockRecorder) CreateGroup(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateGroup", reflect.TypeOf((*MockService)(nil).CreateGroup), arg0, arg1)
}

// CreatePermission mocks base method.
func (m *MockService) CreatePermission(arg0 context.Context, arg1 *dto.NameServiceDescription) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreatePermission", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreatePermission indicates an expected call of CreatePermission.
func (mr *MockServiceMockRecorder) CreatePermission(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreatePermission", reflect.TypeOf((*MockService)(nil).CreatePermission), arg0, arg1)
}

// CreateRole mocks base method.
func (m *MockService) CreateRole(arg0 context.Context, arg1 *dto.NameServiceDescription) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateRole", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateRole indicates an expected call of CreateRole.
func (mr *MockServiceMockRecorder) CreateRole(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateRole", reflect.TypeOf((*MockService)(nil).CreateRole), arg0, arg1)
}

// CreateToken mocks base method.
func (m *MockService) CreateToken(arg0 context.Context, arg1 *dto.UserIdInstance) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateToken", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateToken indicates an expected call of CreateToken.
func (mr *MockServiceMockRecorder) CreateToken(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateToken", reflect.TypeOf((*MockService)(nil).CreateToken), arg0, arg1)
}

// DeleteGroup mocks base method.
func (m *MockService) DeleteGroup(arg0 context.Context, arg1 *dto.NameService) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteGroup", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteGroup indicates an expected call of DeleteGroup.
func (mr *MockServiceMockRecorder) DeleteGroup(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteGroup", reflect.TypeOf((*MockService)(nil).DeleteGroup), arg0, arg1)
}

// DeletePermission mocks base method.
func (m *MockService) DeletePermission(arg0 context.Context, arg1 *dto.NameService) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeletePermission", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeletePermission indicates an expected call of DeletePermission.
func (mr *MockServiceMockRecorder) DeletePermission(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeletePermission", reflect.TypeOf((*MockService)(nil).DeletePermission), arg0, arg1)
}

// DeleteRole mocks base method.
func (m *MockService) DeleteRole(arg0 context.Context, arg1 *dto.NameService) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteRole", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteRole indicates an expected call of DeleteRole.
func (mr *MockServiceMockRecorder) DeleteRole(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteRole", reflect.TypeOf((*MockService)(nil).DeleteRole), arg0, arg1)
}

// Login mocks base method.
func (m *MockService) Login(arg0 context.Context, arg1 *dto.LoginPassword) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Login", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Login indicates an expected call of Login.
func (mr *MockServiceMockRecorder) Login(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Login", reflect.TypeOf((*MockService)(nil).Login), arg0, arg1)
}

// Logout mocks base method.
func (m *MockService) Logout(arg0 context.Context, arg1 uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Logout", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Logout indicates an expected call of Logout.
func (mr *MockServiceMockRecorder) Logout(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Logout", reflect.TypeOf((*MockService)(nil).Logout), arg0, arg1)
}

// RegisterInstance mocks base method.
func (m *MockService) RegisterInstance(arg0 context.Context, arg1 *dto.NameServiceSecret) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RegisterInstance", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// RegisterInstance indicates an expected call of RegisterInstance.
func (mr *MockServiceMockRecorder) RegisterInstance(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterInstance", reflect.TypeOf((*MockService)(nil).RegisterInstance), arg0, arg1)
}

// RegisterService mocks base method.
func (m *MockService) RegisterService(arg0 context.Context, arg1 *dto.NameDescription) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RegisterService", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// RegisterService indicates an expected call of RegisterService.
func (mr *MockServiceMockRecorder) RegisterService(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RegisterService", reflect.TypeOf((*MockService)(nil).RegisterService), arg0, arg1)
}

// ServiceNumberedPermissions mocks base method.
func (m *MockService) ServiceNumberedPermissions(arg0 context.Context, arg1 string) (*[]dto.NameNumber, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ServiceNumberedPermissions", arg0, arg1)
	ret0, _ := ret[0].(*[]dto.NameNumber)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ServiceNumberedPermissions indicates an expected call of ServiceNumberedPermissions.
func (mr *MockServiceMockRecorder) ServiceNumberedPermissions(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ServiceNumberedPermissions", reflect.TypeOf((*MockService)(nil).ServiceNumberedPermissions), arg0, arg1)
}

// UserUUIDFromSession mocks base method.
func (m *MockService) UserUUIDFromSession(arg0 context.Context, arg1 string) (uuid.UUID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UserUUIDFromSession", arg0, arg1)
	ret0, _ := ret[0].(uuid.UUID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UserUUIDFromSession indicates an expected call of UserUUIDFromSession.
func (mr *MockServiceMockRecorder) UserUUIDFromSession(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UserUUIDFromSession", reflect.TypeOf((*MockService)(nil).UserUUIDFromSession), arg0, arg1)
}

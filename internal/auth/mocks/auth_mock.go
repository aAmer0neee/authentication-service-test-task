// Code generated by MockGen. DO NOT EDIT.
// Source: auth.go

// Package auth_mock is a generated GoMock package.
package auth_mock

import (
	reflect "reflect"

	domain "github.com/aAmer0neee/authentication-service-test-task/internal/domain"
	gomock "github.com/golang/mock/gomock"
)

// MockAuthService is a mock of AuthService interface.
type MockAuthService struct {
	ctrl     *gomock.Controller
	recorder *MockAuthServiceMockRecorder
}

// MockAuthServiceMockRecorder is the mock recorder for MockAuthService.
type MockAuthServiceMockRecorder struct {
	mock *MockAuthService
}

// NewMockAuthService creates a new mock instance.
func NewMockAuthService(ctrl *gomock.Controller) *MockAuthService {
	mock := &MockAuthService{ctrl: ctrl}
	mock.recorder = &MockAuthServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuthService) EXPECT() *MockAuthServiceMockRecorder {
	return m.recorder
}

// LoginUser mocks base method.
func (m *MockAuthService) LoginUser(user *domain.User) (*domain.Tokens, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LoginUser", user)
	ret0, _ := ret[0].(*domain.Tokens)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LoginUser indicates an expected call of LoginUser.
func (mr *MockAuthServiceMockRecorder) LoginUser(user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoginUser", reflect.TypeOf((*MockAuthService)(nil).LoginUser), user)
}

// RefreshToken mocks base method.
func (m *MockAuthService) RefreshToken(inputUser *domain.User) (*domain.Tokens, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RefreshToken", inputUser)
	ret0, _ := ret[0].(*domain.Tokens)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RefreshToken indicates an expected call of RefreshToken.
func (mr *MockAuthServiceMockRecorder) RefreshToken(inputUser interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RefreshToken", reflect.TypeOf((*MockAuthService)(nil).RefreshToken), inputUser)
}

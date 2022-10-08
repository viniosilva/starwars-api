// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/viniosilva/starwars-api/internal/service (interfaces: HealthService)

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockHealthService is a mock of HealthService interface.
type MockHealthService struct {
	ctrl     *gomock.Controller
	recorder *MockHealthServiceMockRecorder
}

// MockHealthServiceMockRecorder is the mock recorder for MockHealthService.
type MockHealthServiceMockRecorder struct {
	mock *MockHealthService
}

// NewMockHealthService creates a new mock instance.
func NewMockHealthService(ctrl *gomock.Controller) *MockHealthService {
	mock := &MockHealthService{ctrl: ctrl}
	mock.recorder = &MockHealthServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockHealthService) EXPECT() *MockHealthServiceMockRecorder {
	return m.recorder
}

// Ping mocks base method.
func (m *MockHealthService) Ping(arg0 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Ping", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Ping indicates an expected call of Ping.
func (mr *MockHealthServiceMockRecorder) Ping(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Ping", reflect.TypeOf((*MockHealthService)(nil).Ping), arg0)
}
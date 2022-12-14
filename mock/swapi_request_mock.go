// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/viniosilva/starwars-api/internal/request (interfaces: SwapiRequest)

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	model "github.com/viniosilva/starwars-api/internal/model"
)

// MockSwapiRequest is a mock of SwapiRequest interface.
type MockSwapiRequest struct {
	ctrl     *gomock.Controller
	recorder *MockSwapiRequestMockRecorder
}

// MockSwapiRequestMockRecorder is the mock recorder for MockSwapiRequest.
type MockSwapiRequestMockRecorder struct {
	mock *MockSwapiRequest
}

// NewMockSwapiRequest creates a new mock instance.
func NewMockSwapiRequest(ctrl *gomock.Controller) *MockSwapiRequest {
	mock := &MockSwapiRequest{ctrl: ctrl}
	mock.recorder = &MockSwapiRequestMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSwapiRequest) EXPECT() *MockSwapiRequestMockRecorder {
	return m.recorder
}

// GetFilms mocks base method.
func (m *MockSwapiRequest) GetFilms(arg0 context.Context, arg1 int) (*model.SwapiFilmsResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFilms", arg0, arg1)
	ret0, _ := ret[0].(*model.SwapiFilmsResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFilms indicates an expected call of GetFilms.
func (mr *MockSwapiRequestMockRecorder) GetFilms(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFilms", reflect.TypeOf((*MockSwapiRequest)(nil).GetFilms), arg0, arg1)
}

// GetPlanets mocks base method.
func (m *MockSwapiRequest) GetPlanets(arg0 context.Context, arg1 int) (*model.SwapiPlanetsResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPlanets", arg0, arg1)
	ret0, _ := ret[0].(*model.SwapiPlanetsResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPlanets indicates an expected call of GetPlanets.
func (mr *MockSwapiRequestMockRecorder) GetPlanets(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPlanets", reflect.TypeOf((*MockSwapiRequest)(nil).GetPlanets), arg0, arg1)
}

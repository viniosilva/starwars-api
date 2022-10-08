// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/viniosilva/starwars-api/internal/service (interfaces: FilmService)

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	model "github.com/viniosilva/starwars-api/internal/model"
)

// MockFilmService is a mock of FilmService interface.
type MockFilmService struct {
	ctrl     *gomock.Controller
	recorder *MockFilmServiceMockRecorder
}

// MockFilmServiceMockRecorder is the mock recorder for MockFilmService.
type MockFilmServiceMockRecorder struct {
	mock *MockFilmService
}

// NewMockFilmService creates a new mock instance.
func NewMockFilmService(ctrl *gomock.Controller) *MockFilmService {
	mock := &MockFilmService{ctrl: ctrl}
	mock.recorder = &MockFilmServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockFilmService) EXPECT() *MockFilmServiceMockRecorder {
	return m.recorder
}

// CreateFilms mocks base method.
func (m *MockFilmService) CreateFilms(arg0 context.Context, arg1 []model.Film) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateFilms", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateFilms indicates an expected call of CreateFilms.
func (mr *MockFilmServiceMockRecorder) CreateFilms(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateFilms", reflect.TypeOf((*MockFilmService)(nil).CreateFilms), arg0, arg1)
}

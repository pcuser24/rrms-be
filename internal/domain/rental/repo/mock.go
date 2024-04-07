// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/user2410/rrms-backend/internal/domain/rental/repo (interfaces: Repo)
//
// Generated by this command:
//
//	mockgen -package repo -destination internal/domain/rental/repo/mock.go github.com/user2410/rrms-backend/internal/domain/rental/repo Repo
//

// Package repo is a generated GoMock package.
package repo

import (
	context "context"
	reflect "reflect"

	dto "github.com/user2410/rrms-backend/internal/domain/rental/dto"
	model "github.com/user2410/rrms-backend/internal/domain/rental/model"
	gomock "go.uber.org/mock/gomock"
)

// MockRepo is a mock of Repo interface.
type MockRepo struct {
	ctrl     *gomock.Controller
	recorder *MockRepoMockRecorder
}

// MockRepoMockRecorder is the mock recorder for MockRepo.
type MockRepoMockRecorder struct {
	mock *MockRepo
}

// NewMockRepo creates a new mock instance.
func NewMockRepo(ctrl *gomock.Controller) *MockRepo {
	mock := &MockRepo{ctrl: ctrl}
	mock.recorder = &MockRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepo) EXPECT() *MockRepoMockRecorder {
	return m.recorder
}

// CreateRental mocks base method.
func (m *MockRepo) CreateRental(arg0 context.Context, arg1 *dto.CreateRental) (*model.RentalModel, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateRental", arg0, arg1)
	ret0, _ := ret[0].(*model.RentalModel)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateRental indicates an expected call of CreateRental.
func (mr *MockRepoMockRecorder) CreateRental(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateRental", reflect.TypeOf((*MockRepo)(nil).CreateRental), arg0, arg1)
}

// GetRental mocks base method.
func (m *MockRepo) GetRental(arg0 context.Context, arg1 int64) (*model.RentalModel, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRental", arg0, arg1)
	ret0, _ := ret[0].(*model.RentalModel)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRental indicates an expected call of GetRental.
func (mr *MockRepoMockRecorder) GetRental(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRental", reflect.TypeOf((*MockRepo)(nil).GetRental), arg0, arg1)
}

// UpdateRental mocks base method.
func (m *MockRepo) UpdateRental(arg0 context.Context, arg1 *dto.UpdateRental, arg2 int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateRental", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateRental indicates an expected call of UpdateRental.
func (mr *MockRepoMockRecorder) UpdateRental(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateRental", reflect.TypeOf((*MockRepo)(nil).UpdateRental), arg0, arg1, arg2)
}
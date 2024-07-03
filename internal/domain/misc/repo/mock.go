// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/user2410/rrms-backend/internal/domain/misc/repo (interfaces: Repo)
//
// Generated by this command:
//
//	mockgen -package repo -destination internal/domain/misc/repo/mock.go github.com/user2410/rrms-backend/internal/domain/misc/repo Repo
//

// Package repo is a generated GoMock package.
package repo

import (
	context "context"
	reflect "reflect"

	uuid "github.com/google/uuid"
	dto "github.com/user2410/rrms-backend/internal/domain/misc/dto"
	model "github.com/user2410/rrms-backend/internal/domain/misc/model"
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

// CreateNotification mocks base method.
func (m *MockRepo) CreateNotification(arg0 context.Context, arg1 *dto.CreateNotification) ([]model.Notification, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateNotification", arg0, arg1)
	ret0, _ := ret[0].([]model.Notification)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateNotification indicates an expected call of CreateNotification.
func (mr *MockRepoMockRecorder) CreateNotification(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateNotification", reflect.TypeOf((*MockRepo)(nil).CreateNotification), arg0, arg1)
}

// CreateNotificationDevice mocks base method.
func (m *MockRepo) CreateNotificationDevice(arg0 context.Context, arg1, arg2 uuid.UUID, arg3 *dto.CreateNotificationDevice) (model.NotificationDevice, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateNotificationDevice", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(model.NotificationDevice)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateNotificationDevice indicates an expected call of CreateNotificationDevice.
func (mr *MockRepoMockRecorder) CreateNotificationDevice(arg0, arg1, arg2, arg3 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateNotificationDevice", reflect.TypeOf((*MockRepo)(nil).CreateNotificationDevice), arg0, arg1, arg2, arg3)
}

// DeleteExpiredTokens mocks base method.
func (m *MockRepo) DeleteExpiredTokens(arg0 context.Context, arg1 int32) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteExpiredTokens", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteExpiredTokens indicates an expected call of DeleteExpiredTokens.
func (mr *MockRepoMockRecorder) DeleteExpiredTokens(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteExpiredTokens", reflect.TypeOf((*MockRepo)(nil).DeleteExpiredTokens), arg0, arg1)
}

// GetNotificationDevice mocks base method.
func (m *MockRepo) GetNotificationDevice(arg0 context.Context, arg1, arg2 uuid.UUID, arg3, arg4 string) ([]model.NotificationDevice, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNotificationDevice", arg0, arg1, arg2, arg3, arg4)
	ret0, _ := ret[0].([]model.NotificationDevice)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetNotificationDevice indicates an expected call of GetNotificationDevice.
func (mr *MockRepoMockRecorder) GetNotificationDevice(arg0, arg1, arg2, arg3, arg4 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNotificationDevice", reflect.TypeOf((*MockRepo)(nil).GetNotificationDevice), arg0, arg1, arg2, arg3, arg4)
}

// GetNotificationsOfUser mocks base method.
func (m *MockRepo) GetNotificationsOfUser(arg0 context.Context, arg1 uuid.UUID, arg2 dto.GetNotificationsOfUserQuery) ([]model.Notification, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNotificationsOfUser", arg0, arg1, arg2)
	ret0, _ := ret[0].([]model.Notification)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetNotificationsOfUser indicates an expected call of GetNotificationsOfUser.
func (mr *MockRepoMockRecorder) GetNotificationsOfUser(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNotificationsOfUser", reflect.TypeOf((*MockRepo)(nil).GetNotificationsOfUser), arg0, arg1, arg2)
}

// UpdateNotification mocks base method.
func (m *MockRepo) UpdateNotification(arg0 context.Context, arg1 *dto.UpdateNotification) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateNotification", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateNotification indicates an expected call of UpdateNotification.
func (mr *MockRepoMockRecorder) UpdateNotification(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateNotification", reflect.TypeOf((*MockRepo)(nil).UpdateNotification), arg0, arg1)
}

// UpdateNotificationDeviceTokenTimestamp mocks base method.
func (m *MockRepo) UpdateNotificationDeviceTokenTimestamp(arg0 context.Context, arg1, arg2 uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateNotificationDeviceTokenTimestamp", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateNotificationDeviceTokenTimestamp indicates an expected call of UpdateNotificationDeviceTokenTimestamp.
func (mr *MockRepoMockRecorder) UpdateNotificationDeviceTokenTimestamp(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateNotificationDeviceTokenTimestamp", reflect.TypeOf((*MockRepo)(nil).UpdateNotificationDeviceTokenTimestamp), arg0, arg1, arg2)
}

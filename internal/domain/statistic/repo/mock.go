// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/user2410/rrms-backend/internal/domain/statistic/repo (interfaces: Repo)
//
// Generated by this command:
//
//	mockgen -package repo -destination internal/domain/statistic/repo/mock.go github.com/user2410/rrms-backend/internal/domain/statistic/repo Repo
//

// Package repo is a generated GoMock package.
package repo

import (
	context "context"
	reflect "reflect"
	time "time"

	uuid "github.com/google/uuid"
	dto "github.com/user2410/rrms-backend/internal/domain/statistic/dto"
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

// GetLeastRentedProperties mocks base method.
func (m *MockRepo) GetLeastRentedProperties(arg0 context.Context, arg1 uuid.UUID, arg2, arg3 int32) ([]dto.ExtremelyRentedPropertyItem, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLeastRentedProperties", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].([]dto.ExtremelyRentedPropertyItem)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLeastRentedProperties indicates an expected call of GetLeastRentedProperties.
func (mr *MockRepoMockRecorder) GetLeastRentedProperties(arg0, arg1, arg2, arg3 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLeastRentedProperties", reflect.TypeOf((*MockRepo)(nil).GetLeastRentedProperties), arg0, arg1, arg2, arg3)
}

// GetLeastRentedUnits mocks base method.
func (m *MockRepo) GetLeastRentedUnits(arg0 context.Context, arg1 uuid.UUID, arg2, arg3 int32) ([]dto.ExtremelyRentedUnitItem, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLeastRentedUnits", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].([]dto.ExtremelyRentedUnitItem)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLeastRentedUnits indicates an expected call of GetLeastRentedUnits.
func (mr *MockRepoMockRecorder) GetLeastRentedUnits(arg0, arg1, arg2, arg3 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLeastRentedUnits", reflect.TypeOf((*MockRepo)(nil).GetLeastRentedUnits), arg0, arg1, arg2, arg3)
}

// GetMaintenanceRequests mocks base method.
func (m *MockRepo) GetMaintenanceRequests(arg0 context.Context, arg1 uuid.UUID, arg2 time.Time) ([]int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMaintenanceRequests", arg0, arg1, arg2)
	ret0, _ := ret[0].([]int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMaintenanceRequests indicates an expected call of GetMaintenanceRequests.
func (mr *MockRepoMockRecorder) GetMaintenanceRequests(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMaintenanceRequests", reflect.TypeOf((*MockRepo)(nil).GetMaintenanceRequests), arg0, arg1, arg2)
}

// GetManagedProperties mocks base method.
func (m *MockRepo) GetManagedProperties(arg0 context.Context, arg1 uuid.UUID) ([]uuid.UUID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetManagedProperties", arg0, arg1)
	ret0, _ := ret[0].([]uuid.UUID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetManagedProperties indicates an expected call of GetManagedProperties.
func (mr *MockRepoMockRecorder) GetManagedProperties(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetManagedProperties", reflect.TypeOf((*MockRepo)(nil).GetManagedProperties), arg0, arg1)
}

// GetManagedPropertiesByRole mocks base method.
func (m *MockRepo) GetManagedPropertiesByRole(arg0 context.Context, arg1 uuid.UUID, arg2 string) ([]uuid.UUID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetManagedPropertiesByRole", arg0, arg1, arg2)
	ret0, _ := ret[0].([]uuid.UUID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetManagedPropertiesByRole indicates an expected call of GetManagedPropertiesByRole.
func (mr *MockRepoMockRecorder) GetManagedPropertiesByRole(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetManagedPropertiesByRole", reflect.TypeOf((*MockRepo)(nil).GetManagedPropertiesByRole), arg0, arg1, arg2)
}

// GetManagedUnits mocks base method.
func (m *MockRepo) GetManagedUnits(arg0 context.Context, arg1 uuid.UUID) ([]uuid.UUID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetManagedUnits", arg0, arg1)
	ret0, _ := ret[0].([]uuid.UUID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetManagedUnits indicates an expected call of GetManagedUnits.
func (mr *MockRepoMockRecorder) GetManagedUnits(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetManagedUnits", reflect.TypeOf((*MockRepo)(nil).GetManagedUnits), arg0, arg1)
}

// GetMostRentedProperties mocks base method.
func (m *MockRepo) GetMostRentedProperties(arg0 context.Context, arg1 uuid.UUID, arg2, arg3 int32) ([]dto.ExtremelyRentedPropertyItem, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMostRentedProperties", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].([]dto.ExtremelyRentedPropertyItem)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMostRentedProperties indicates an expected call of GetMostRentedProperties.
func (mr *MockRepoMockRecorder) GetMostRentedProperties(arg0, arg1, arg2, arg3 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMostRentedProperties", reflect.TypeOf((*MockRepo)(nil).GetMostRentedProperties), arg0, arg1, arg2, arg3)
}

// GetMostRentedUnits mocks base method.
func (m *MockRepo) GetMostRentedUnits(arg0 context.Context, arg1 uuid.UUID, arg2, arg3 int32) ([]dto.ExtremelyRentedUnitItem, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMostRentedUnits", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].([]dto.ExtremelyRentedUnitItem)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMostRentedUnits indicates an expected call of GetMostRentedUnits.
func (mr *MockRepoMockRecorder) GetMostRentedUnits(arg0, arg1, arg2, arg3 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMostRentedUnits", reflect.TypeOf((*MockRepo)(nil).GetMostRentedUnits), arg0, arg1, arg2, arg3)
}

// GetNewApplications mocks base method.
func (m *MockRepo) GetNewApplications(arg0 context.Context, arg1 uuid.UUID, arg2 time.Time) ([]int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNewApplications", arg0, arg1, arg2)
	ret0, _ := ret[0].([]int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetNewApplications indicates an expected call of GetNewApplications.
func (mr *MockRepoMockRecorder) GetNewApplications(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNewApplications", reflect.TypeOf((*MockRepo)(nil).GetNewApplications), arg0, arg1, arg2)
}

// GetOccupiedProperties mocks base method.
func (m *MockRepo) GetOccupiedProperties(arg0 context.Context, arg1 uuid.UUID) ([]uuid.UUID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOccupiedProperties", arg0, arg1)
	ret0, _ := ret[0].([]uuid.UUID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOccupiedProperties indicates an expected call of GetOccupiedProperties.
func (mr *MockRepoMockRecorder) GetOccupiedProperties(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOccupiedProperties", reflect.TypeOf((*MockRepo)(nil).GetOccupiedProperties), arg0, arg1)
}

// GetOccupiedUnits mocks base method.
func (m *MockRepo) GetOccupiedUnits(arg0 context.Context, arg1 uuid.UUID) ([]uuid.UUID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOccupiedUnits", arg0, arg1)
	ret0, _ := ret[0].([]uuid.UUID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOccupiedUnits indicates an expected call of GetOccupiedUnits.
func (mr *MockRepoMockRecorder) GetOccupiedUnits(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOccupiedUnits", reflect.TypeOf((*MockRepo)(nil).GetOccupiedUnits), arg0, arg1)
}

// GetPropertiesHavingListing mocks base method.
func (m *MockRepo) GetPropertiesHavingListing(arg0 context.Context, arg1 uuid.UUID) ([]uuid.UUID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPropertiesHavingListing", arg0, arg1)
	ret0, _ := ret[0].([]uuid.UUID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPropertiesHavingListing indicates an expected call of GetPropertiesHavingListing.
func (mr *MockRepoMockRecorder) GetPropertiesHavingListing(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPropertiesHavingListing", reflect.TypeOf((*MockRepo)(nil).GetPropertiesHavingListing), arg0, arg1)
}

// GetPropertiesWithActiveListing mocks base method.
func (m *MockRepo) GetPropertiesWithActiveListing(arg0 context.Context, arg1 uuid.UUID) ([]uuid.UUID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPropertiesWithActiveListing", arg0, arg1)
	ret0, _ := ret[0].([]uuid.UUID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPropertiesWithActiveListing indicates an expected call of GetPropertiesWithActiveListing.
func (mr *MockRepoMockRecorder) GetPropertiesWithActiveListing(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPropertiesWithActiveListing", reflect.TypeOf((*MockRepo)(nil).GetPropertiesWithActiveListing), arg0, arg1)
}

// GetRentalPaymentArrears mocks base method.
func (m *MockRepo) GetRentalPaymentArrears(arg0 context.Context, arg1 uuid.UUID, arg2 dto.RentalPaymentStatisticQuery) ([]dto.RentalPayment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRentalPaymentArrears", arg0, arg1, arg2)
	ret0, _ := ret[0].([]dto.RentalPayment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRentalPaymentArrears indicates an expected call of GetRentalPaymentArrears.
func (mr *MockRepoMockRecorder) GetRentalPaymentArrears(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRentalPaymentArrears", reflect.TypeOf((*MockRepo)(nil).GetRentalPaymentArrears), arg0, arg1, arg2)
}

// GetRentalPaymentIncomes mocks base method.
func (m *MockRepo) GetRentalPaymentIncomes(arg0 context.Context, arg1 uuid.UUID, arg2 dto.RentalPaymentStatisticQuery) (float32, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRentalPaymentIncomes", arg0, arg1, arg2)
	ret0, _ := ret[0].(float32)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRentalPaymentIncomes indicates an expected call of GetRentalPaymentIncomes.
func (mr *MockRepoMockRecorder) GetRentalPaymentIncomes(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRentalPaymentIncomes", reflect.TypeOf((*MockRepo)(nil).GetRentalPaymentIncomes), arg0, arg1, arg2)
}

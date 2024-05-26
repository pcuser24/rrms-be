// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/user2410/rrms-backend/internal/domain/listing/repo (interfaces: Repo)
//
// Generated by this command:
//
//	mockgen -package repo -destination internal/domain/listing/repo/mock.go github.com/user2410/rrms-backend/internal/domain/listing/repo Repo
//

// Package repo is a generated GoMock package.
package repo

import (
	context "context"
	reflect "reflect"

	uuid "github.com/google/uuid"
	dto "github.com/user2410/rrms-backend/internal/domain/listing/dto"
	model "github.com/user2410/rrms-backend/internal/domain/listing/model"
	model0 "github.com/user2410/rrms-backend/internal/domain/payment/model"
	service "github.com/user2410/rrms-backend/internal/domain/payment/service"
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

// CheckListingExpired mocks base method.
func (m *MockRepo) CheckListingExpired(arg0 context.Context, arg1 uuid.UUID) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckListingExpired", arg0, arg1)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckListingExpired indicates an expected call of CheckListingExpired.
func (mr *MockRepoMockRecorder) CheckListingExpired(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckListingExpired", reflect.TypeOf((*MockRepo)(nil).CheckListingExpired), arg0, arg1)
}

// CheckListingOwnership mocks base method.
func (m *MockRepo) CheckListingOwnership(arg0 context.Context, arg1, arg2 uuid.UUID) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckListingOwnership", arg0, arg1, arg2)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckListingOwnership indicates an expected call of CheckListingOwnership.
func (mr *MockRepoMockRecorder) CheckListingOwnership(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckListingOwnership", reflect.TypeOf((*MockRepo)(nil).CheckListingOwnership), arg0, arg1, arg2)
}

// CheckListingVisibility mocks base method.
func (m *MockRepo) CheckListingVisibility(arg0 context.Context, arg1, arg2 uuid.UUID) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckListingVisibility", arg0, arg1, arg2)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckListingVisibility indicates an expected call of CheckListingVisibility.
func (mr *MockRepoMockRecorder) CheckListingVisibility(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckListingVisibility", reflect.TypeOf((*MockRepo)(nil).CheckListingVisibility), arg0, arg1, arg2)
}

// CheckValidUnitForListing mocks base method.
func (m *MockRepo) CheckValidUnitForListing(arg0 context.Context, arg1, arg2 uuid.UUID) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckValidUnitForListing", arg0, arg1, arg2)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckValidUnitForListing indicates an expected call of CheckValidUnitForListing.
func (mr *MockRepoMockRecorder) CheckValidUnitForListing(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckValidUnitForListing", reflect.TypeOf((*MockRepo)(nil).CheckValidUnitForListing), arg0, arg1, arg2)
}

// CreateListing mocks base method.
func (m *MockRepo) CreateListing(arg0 context.Context, arg1 *dto.CreateListing) (*model.ListingModel, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateListing", arg0, arg1)
	ret0, _ := ret[0].(*model.ListingModel)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateListing indicates an expected call of CreateListing.
func (mr *MockRepoMockRecorder) CreateListing(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateListing", reflect.TypeOf((*MockRepo)(nil).CreateListing), arg0, arg1)
}

// DeleteListing mocks base method.
func (m *MockRepo) DeleteListing(arg0 context.Context, arg1 uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteListing", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteListing indicates an expected call of DeleteListing.
func (mr *MockRepoMockRecorder) DeleteListing(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteListing", reflect.TypeOf((*MockRepo)(nil).DeleteListing), arg0, arg1)
}

// FilterVisibleListings mocks base method.
func (m *MockRepo) FilterVisibleListings(arg0 context.Context, arg1 []uuid.UUID, arg2 uuid.UUID) ([]uuid.UUID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FilterVisibleListings", arg0, arg1, arg2)
	ret0, _ := ret[0].([]uuid.UUID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FilterVisibleListings indicates an expected call of FilterVisibleListings.
func (mr *MockRepoMockRecorder) FilterVisibleListings(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FilterVisibleListings", reflect.TypeOf((*MockRepo)(nil).FilterVisibleListings), arg0, arg1, arg2)
}

// GetListingByID mocks base method.
func (m *MockRepo) GetListingByID(arg0 context.Context, arg1 uuid.UUID) (*model.ListingModel, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetListingByID", arg0, arg1)
	ret0, _ := ret[0].(*model.ListingModel)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetListingByID indicates an expected call of GetListingByID.
func (mr *MockRepoMockRecorder) GetListingByID(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetListingByID", reflect.TypeOf((*MockRepo)(nil).GetListingByID), arg0, arg1)
}

// GetListingPayments mocks base method.
func (m *MockRepo) GetListingPayments(arg0 context.Context, arg1 uuid.UUID) ([]model0.PaymentModel, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetListingPayments", arg0, arg1)
	ret0, _ := ret[0].([]model0.PaymentModel)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetListingPayments indicates an expected call of GetListingPayments.
func (mr *MockRepoMockRecorder) GetListingPayments(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetListingPayments", reflect.TypeOf((*MockRepo)(nil).GetListingPayments), arg0, arg1)
}

// GetListingPaymentsByType mocks base method.
func (m *MockRepo) GetListingPaymentsByType(arg0 context.Context, arg1 uuid.UUID, arg2 service.PAYMENTTYPE) ([]model0.PaymentModel, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetListingPaymentsByType", arg0, arg1, arg2)
	ret0, _ := ret[0].([]model0.PaymentModel)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetListingPaymentsByType indicates an expected call of GetListingPaymentsByType.
func (mr *MockRepoMockRecorder) GetListingPaymentsByType(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetListingPaymentsByType", reflect.TypeOf((*MockRepo)(nil).GetListingPaymentsByType), arg0, arg1, arg2)
}

// GetListingsByIds mocks base method.
func (m *MockRepo) GetListingsByIds(arg0 context.Context, arg1 []uuid.UUID, arg2 []string) ([]model.ListingModel, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetListingsByIds", arg0, arg1, arg2)
	ret0, _ := ret[0].([]model.ListingModel)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetListingsByIds indicates an expected call of GetListingsByIds.
func (mr *MockRepoMockRecorder) GetListingsByIds(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetListingsByIds", reflect.TypeOf((*MockRepo)(nil).GetListingsByIds), arg0, arg1, arg2)
}

// SearchListingCombination mocks base method.
func (m *MockRepo) SearchListingCombination(arg0 context.Context, arg1 *dto.SearchListingCombinationQuery) (*dto.SearchListingCombinationResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchListingCombination", arg0, arg1)
	ret0, _ := ret[0].(*dto.SearchListingCombinationResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SearchListingCombination indicates an expected call of SearchListingCombination.
func (mr *MockRepoMockRecorder) SearchListingCombination(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchListingCombination", reflect.TypeOf((*MockRepo)(nil).SearchListingCombination), arg0, arg1)
}

// UpdateListing mocks base method.
func (m *MockRepo) UpdateListing(arg0 context.Context, arg1 uuid.UUID, arg2 *dto.UpdateListing) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateListing", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateListing indicates an expected call of UpdateListing.
func (mr *MockRepoMockRecorder) UpdateListing(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateListing", reflect.TypeOf((*MockRepo)(nil).UpdateListing), arg0, arg1, arg2)
}

// UpdateListingExpiration mocks base method.
func (m *MockRepo) UpdateListingExpiration(arg0 context.Context, arg1 uuid.UUID, arg2 int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateListingExpiration", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateListingExpiration indicates an expected call of UpdateListingExpiration.
func (mr *MockRepoMockRecorder) UpdateListingExpiration(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateListingExpiration", reflect.TypeOf((*MockRepo)(nil).UpdateListingExpiration), arg0, arg1, arg2)
}

// UpdateListingPriority mocks base method.
func (m *MockRepo) UpdateListingPriority(arg0 context.Context, arg1 uuid.UUID, arg2 int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateListingPriority", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateListingPriority indicates an expected call of UpdateListingPriority.
func (mr *MockRepoMockRecorder) UpdateListingPriority(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateListingPriority", reflect.TypeOf((*MockRepo)(nil).UpdateListingPriority), arg0, arg1, arg2)
}

// UpdateListingStatus mocks base method.
func (m *MockRepo) UpdateListingStatus(arg0 context.Context, arg1 uuid.UUID, arg2 bool) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateListingStatus", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateListingStatus indicates an expected call of UpdateListingStatus.
func (mr *MockRepoMockRecorder) UpdateListingStatus(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateListingStatus", reflect.TypeOf((*MockRepo)(nil).UpdateListingStatus), arg0, arg1, arg2)
}

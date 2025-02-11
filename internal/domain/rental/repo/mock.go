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

	uuid "github.com/google/uuid"
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

// CheckPreRentalVisibility mocks base method.
func (m *MockRepo) CheckPreRentalVisibility(arg0 context.Context, arg1 int64, arg2 uuid.UUID) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckPreRentalVisibility", arg0, arg1, arg2)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckPreRentalVisibility indicates an expected call of CheckPreRentalVisibility.
func (mr *MockRepoMockRecorder) CheckPreRentalVisibility(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckPreRentalVisibility", reflect.TypeOf((*MockRepo)(nil).CheckPreRentalVisibility), arg0, arg1, arg2)
}

// CheckRentalVisibility mocks base method.
func (m *MockRepo) CheckRentalVisibility(arg0 context.Context, arg1 int64, arg2 uuid.UUID) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckRentalVisibility", arg0, arg1, arg2)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckRentalVisibility indicates an expected call of CheckRentalVisibility.
func (mr *MockRepoMockRecorder) CheckRentalVisibility(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckRentalVisibility", reflect.TypeOf((*MockRepo)(nil).CheckRentalVisibility), arg0, arg1, arg2)
}

// CreateContract mocks base method.
func (m *MockRepo) CreateContract(arg0 context.Context, arg1 *dto.CreateContract) (*model.ContractModel, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateContract", arg0, arg1)
	ret0, _ := ret[0].(*model.ContractModel)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateContract indicates an expected call of CreateContract.
func (mr *MockRepoMockRecorder) CreateContract(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateContract", reflect.TypeOf((*MockRepo)(nil).CreateContract), arg0, arg1)
}

// CreatePreRental mocks base method.
func (m *MockRepo) CreatePreRental(arg0 context.Context, arg1 *dto.CreateRental) (model.RentalModel, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreatePreRental", arg0, arg1)
	ret0, _ := ret[0].(model.RentalModel)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreatePreRental indicates an expected call of CreatePreRental.
func (mr *MockRepoMockRecorder) CreatePreRental(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreatePreRental", reflect.TypeOf((*MockRepo)(nil).CreatePreRental), arg0, arg1)
}

// CreateRental mocks base method.
func (m *MockRepo) CreateRental(arg0 context.Context, arg1 *dto.CreateRental) (model.RentalModel, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateRental", arg0, arg1)
	ret0, _ := ret[0].(model.RentalModel)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateRental indicates an expected call of CreateRental.
func (mr *MockRepoMockRecorder) CreateRental(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateRental", reflect.TypeOf((*MockRepo)(nil).CreateRental), arg0, arg1)
}

// CreateRentalComplaint mocks base method.
func (m *MockRepo) CreateRentalComplaint(arg0 context.Context, arg1 *dto.CreateRentalComplaint) (model.RentalComplaint, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateRentalComplaint", arg0, arg1)
	ret0, _ := ret[0].(model.RentalComplaint)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateRentalComplaint indicates an expected call of CreateRentalComplaint.
func (mr *MockRepoMockRecorder) CreateRentalComplaint(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateRentalComplaint", reflect.TypeOf((*MockRepo)(nil).CreateRentalComplaint), arg0, arg1)
}

// CreateRentalComplaintReply mocks base method.
func (m *MockRepo) CreateRentalComplaintReply(arg0 context.Context, arg1 *dto.CreateRentalComplaintReply) (model.RentalComplaintReply, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateRentalComplaintReply", arg0, arg1)
	ret0, _ := ret[0].(model.RentalComplaintReply)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateRentalComplaintReply indicates an expected call of CreateRentalComplaintReply.
func (mr *MockRepoMockRecorder) CreateRentalComplaintReply(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateRentalComplaintReply", reflect.TypeOf((*MockRepo)(nil).CreateRentalComplaintReply), arg0, arg1)
}

// CreateRentalPayment mocks base method.
func (m *MockRepo) CreateRentalPayment(arg0 context.Context, arg1 *dto.CreateRentalPayment) (model.RentalPayment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateRentalPayment", arg0, arg1)
	ret0, _ := ret[0].(model.RentalPayment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateRentalPayment indicates an expected call of CreateRentalPayment.
func (mr *MockRepoMockRecorder) CreateRentalPayment(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateRentalPayment", reflect.TypeOf((*MockRepo)(nil).CreateRentalPayment), arg0, arg1)
}

// FilterVisibleRentals mocks base method.
func (m *MockRepo) FilterVisibleRentals(arg0 context.Context, arg1 uuid.UUID, arg2 []int64) ([]int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FilterVisibleRentals", arg0, arg1, arg2)
	ret0, _ := ret[0].([]int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FilterVisibleRentals indicates an expected call of FilterVisibleRentals.
func (mr *MockRepoMockRecorder) FilterVisibleRentals(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FilterVisibleRentals", reflect.TypeOf((*MockRepo)(nil).FilterVisibleRentals), arg0, arg1, arg2)
}

// GetContractByID mocks base method.
func (m *MockRepo) GetContractByID(arg0 context.Context, arg1 int64) (*model.ContractModel, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetContractByID", arg0, arg1)
	ret0, _ := ret[0].(*model.ContractModel)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetContractByID indicates an expected call of GetContractByID.
func (mr *MockRepoMockRecorder) GetContractByID(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetContractByID", reflect.TypeOf((*MockRepo)(nil).GetContractByID), arg0, arg1)
}

// GetContractByRentalID mocks base method.
func (m *MockRepo) GetContractByRentalID(arg0 context.Context, arg1 int64) (*model.ContractModel, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetContractByRentalID", arg0, arg1)
	ret0, _ := ret[0].(*model.ContractModel)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetContractByRentalID indicates an expected call of GetContractByRentalID.
func (mr *MockRepoMockRecorder) GetContractByRentalID(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetContractByRentalID", reflect.TypeOf((*MockRepo)(nil).GetContractByRentalID), arg0, arg1)
}

// GetContractsByIds mocks base method.
func (m *MockRepo) GetContractsByIds(arg0 context.Context, arg1 []int64, arg2 []string) ([]model.ContractModel, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetContractsByIds", arg0, arg1, arg2)
	ret0, _ := ret[0].([]model.ContractModel)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetContractsByIds indicates an expected call of GetContractsByIds.
func (mr *MockRepoMockRecorder) GetContractsByIds(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetContractsByIds", reflect.TypeOf((*MockRepo)(nil).GetContractsByIds), arg0, arg1, arg2)
}

// GetManagedPreRentals mocks base method.
func (m *MockRepo) GetManagedPreRentals(arg0 context.Context, arg1 uuid.UUID, arg2 *dto.GetPreRentalsQuery) ([]model.RentalModel, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetManagedPreRentals", arg0, arg1, arg2)
	ret0, _ := ret[0].([]model.RentalModel)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetManagedPreRentals indicates an expected call of GetManagedPreRentals.
func (mr *MockRepoMockRecorder) GetManagedPreRentals(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetManagedPreRentals", reflect.TypeOf((*MockRepo)(nil).GetManagedPreRentals), arg0, arg1, arg2)
}

// GetManagedRentalPayments mocks base method.
func (m *MockRepo) GetManagedRentalPayments(arg0 context.Context, arg1 uuid.UUID, arg2 *dto.GetManagedRentalPaymentsQuery) ([]model.RentalPayment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetManagedRentalPayments", arg0, arg1, arg2)
	ret0, _ := ret[0].([]model.RentalPayment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetManagedRentalPayments indicates an expected call of GetManagedRentalPayments.
func (mr *MockRepoMockRecorder) GetManagedRentalPayments(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetManagedRentalPayments", reflect.TypeOf((*MockRepo)(nil).GetManagedRentalPayments), arg0, arg1, arg2)
}

// GetManagedRentals mocks base method.
func (m *MockRepo) GetManagedRentals(arg0 context.Context, arg1 uuid.UUID, arg2 *dto.GetRentalsQuery) ([]int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetManagedRentals", arg0, arg1, arg2)
	ret0, _ := ret[0].([]int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetManagedRentals indicates an expected call of GetManagedRentals.
func (mr *MockRepoMockRecorder) GetManagedRentals(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetManagedRentals", reflect.TypeOf((*MockRepo)(nil).GetManagedRentals), arg0, arg1, arg2)
}

// GetMyRentals mocks base method.
func (m *MockRepo) GetMyRentals(arg0 context.Context, arg1 uuid.UUID, arg2 *dto.GetRentalsQuery) ([]int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMyRentals", arg0, arg1, arg2)
	ret0, _ := ret[0].([]int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMyRentals indicates an expected call of GetMyRentals.
func (mr *MockRepoMockRecorder) GetMyRentals(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMyRentals", reflect.TypeOf((*MockRepo)(nil).GetMyRentals), arg0, arg1, arg2)
}

// GetPaymentsOfRental mocks base method.
func (m *MockRepo) GetPaymentsOfRental(arg0 context.Context, arg1 int64) ([]model.RentalPayment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPaymentsOfRental", arg0, arg1)
	ret0, _ := ret[0].([]model.RentalPayment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPaymentsOfRental indicates an expected call of GetPaymentsOfRental.
func (mr *MockRepoMockRecorder) GetPaymentsOfRental(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPaymentsOfRental", reflect.TypeOf((*MockRepo)(nil).GetPaymentsOfRental), arg0, arg1)
}

// GetPreRental mocks base method.
func (m *MockRepo) GetPreRental(arg0 context.Context, arg1 int64) (model.RentalModel, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPreRental", arg0, arg1)
	ret0, _ := ret[0].(model.RentalModel)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPreRental indicates an expected call of GetPreRental.
func (mr *MockRepoMockRecorder) GetPreRental(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPreRental", reflect.TypeOf((*MockRepo)(nil).GetPreRental), arg0, arg1)
}

// GetPreRentalsToTenant mocks base method.
func (m *MockRepo) GetPreRentalsToTenant(arg0 context.Context, arg1 uuid.UUID, arg2 *dto.GetPreRentalsQuery) ([]model.RentalModel, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPreRentalsToTenant", arg0, arg1, arg2)
	ret0, _ := ret[0].([]model.RentalModel)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPreRentalsToTenant indicates an expected call of GetPreRentalsToTenant.
func (mr *MockRepoMockRecorder) GetPreRentalsToTenant(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPreRentalsToTenant", reflect.TypeOf((*MockRepo)(nil).GetPreRentalsToTenant), arg0, arg1, arg2)
}

// GetRental mocks base method.
func (m *MockRepo) GetRental(arg0 context.Context, arg1 int64) (model.RentalModel, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRental", arg0, arg1)
	ret0, _ := ret[0].(model.RentalModel)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRental indicates an expected call of GetRental.
func (mr *MockRepoMockRecorder) GetRental(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRental", reflect.TypeOf((*MockRepo)(nil).GetRental), arg0, arg1)
}

// GetRentalComplaint mocks base method.
func (m *MockRepo) GetRentalComplaint(arg0 context.Context, arg1 int64) (model.RentalComplaint, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRentalComplaint", arg0, arg1)
	ret0, _ := ret[0].(model.RentalComplaint)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRentalComplaint indicates an expected call of GetRentalComplaint.
func (mr *MockRepoMockRecorder) GetRentalComplaint(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRentalComplaint", reflect.TypeOf((*MockRepo)(nil).GetRentalComplaint), arg0, arg1)
}

// GetRentalComplaintReplies mocks base method.
func (m *MockRepo) GetRentalComplaintReplies(arg0 context.Context, arg1 int64, arg2, arg3 int32) ([]model.RentalComplaintReply, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRentalComplaintReplies", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].([]model.RentalComplaintReply)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRentalComplaintReplies indicates an expected call of GetRentalComplaintReplies.
func (mr *MockRepoMockRecorder) GetRentalComplaintReplies(arg0, arg1, arg2, arg3 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRentalComplaintReplies", reflect.TypeOf((*MockRepo)(nil).GetRentalComplaintReplies), arg0, arg1, arg2, arg3)
}

// GetRentalComplaintsByRentalId mocks base method.
func (m *MockRepo) GetRentalComplaintsByRentalId(arg0 context.Context, arg1 int64, arg2, arg3 int32) ([]model.RentalComplaint, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRentalComplaintsByRentalId", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].([]model.RentalComplaint)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRentalComplaintsByRentalId indicates an expected call of GetRentalComplaintsByRentalId.
func (mr *MockRepoMockRecorder) GetRentalComplaintsByRentalId(arg0, arg1, arg2, arg3 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRentalComplaintsByRentalId", reflect.TypeOf((*MockRepo)(nil).GetRentalComplaintsByRentalId), arg0, arg1, arg2, arg3)
}

// GetRentalComplaintsOfUser mocks base method.
func (m *MockRepo) GetRentalComplaintsOfUser(arg0 context.Context, arg1 uuid.UUID, arg2 dto.GetRentalComplaintsOfUserQuery) ([]model.RentalComplaint, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRentalComplaintsOfUser", arg0, arg1, arg2)
	ret0, _ := ret[0].([]model.RentalComplaint)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRentalComplaintsOfUser indicates an expected call of GetRentalComplaintsOfUser.
func (mr *MockRepoMockRecorder) GetRentalComplaintsOfUser(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRentalComplaintsOfUser", reflect.TypeOf((*MockRepo)(nil).GetRentalComplaintsOfUser), arg0, arg1, arg2)
}

// GetRentalContractsOfUser mocks base method.
func (m *MockRepo) GetRentalContractsOfUser(arg0 context.Context, arg1 uuid.UUID, arg2 *dto.GetRentalContracts) ([]int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRentalContractsOfUser", arg0, arg1, arg2)
	ret0, _ := ret[0].([]int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRentalContractsOfUser indicates an expected call of GetRentalContractsOfUser.
func (mr *MockRepoMockRecorder) GetRentalContractsOfUser(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRentalContractsOfUser", reflect.TypeOf((*MockRepo)(nil).GetRentalContractsOfUser), arg0, arg1, arg2)
}

// GetRentalPayment mocks base method.
func (m *MockRepo) GetRentalPayment(arg0 context.Context, arg1 int64) (model.RentalPayment, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRentalPayment", arg0, arg1)
	ret0, _ := ret[0].(model.RentalPayment)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRentalPayment indicates an expected call of GetRentalPayment.
func (mr *MockRepoMockRecorder) GetRentalPayment(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRentalPayment", reflect.TypeOf((*MockRepo)(nil).GetRentalPayment), arg0, arg1)
}

// GetRentalSide mocks base method.
func (m *MockRepo) GetRentalSide(arg0 context.Context, arg1 int64, arg2 uuid.UUID) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRentalSide", arg0, arg1, arg2)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRentalSide indicates an expected call of GetRentalSide.
func (mr *MockRepoMockRecorder) GetRentalSide(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRentalSide", reflect.TypeOf((*MockRepo)(nil).GetRentalSide), arg0, arg1, arg2)
}

// GetRentalsByIds mocks base method.
func (m *MockRepo) GetRentalsByIds(arg0 context.Context, arg1 []int64, arg2 []string) ([]model.RentalModel, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRentalsByIds", arg0, arg1, arg2)
	ret0, _ := ret[0].([]model.RentalModel)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRentalsByIds indicates an expected call of GetRentalsByIds.
func (mr *MockRepoMockRecorder) GetRentalsByIds(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRentalsByIds", reflect.TypeOf((*MockRepo)(nil).GetRentalsByIds), arg0, arg1, arg2)
}

// MovePreRentalToRental mocks base method.
func (m *MockRepo) MovePreRentalToRental(arg0 context.Context, arg1 int64) (model.RentalModel, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MovePreRentalToRental", arg0, arg1)
	ret0, _ := ret[0].(model.RentalModel)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// MovePreRentalToRental indicates an expected call of MovePreRentalToRental.
func (mr *MockRepoMockRecorder) MovePreRentalToRental(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MovePreRentalToRental", reflect.TypeOf((*MockRepo)(nil).MovePreRentalToRental), arg0, arg1)
}

// PingRentalContract mocks base method.
func (m *MockRepo) PingRentalContract(arg0 context.Context, arg1 int64) (any, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PingRentalContract", arg0, arg1)
	ret0, _ := ret[0].(any)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PingRentalContract indicates an expected call of PingRentalContract.
func (mr *MockRepoMockRecorder) PingRentalContract(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PingRentalContract", reflect.TypeOf((*MockRepo)(nil).PingRentalContract), arg0, arg1)
}

// PlanRentalPayment mocks base method.
func (m *MockRepo) PlanRentalPayment(arg0 context.Context, arg1 int64) ([]int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PlanRentalPayment", arg0, arg1)
	ret0, _ := ret[0].([]int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PlanRentalPayment indicates an expected call of PlanRentalPayment.
func (mr *MockRepoMockRecorder) PlanRentalPayment(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PlanRentalPayment", reflect.TypeOf((*MockRepo)(nil).PlanRentalPayment), arg0, arg1)
}

// PlanRentalPayments mocks base method.
func (m *MockRepo) PlanRentalPayments(arg0 context.Context) ([]int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PlanRentalPayments", arg0)
	ret0, _ := ret[0].([]int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PlanRentalPayments indicates an expected call of PlanRentalPayments.
func (mr *MockRepoMockRecorder) PlanRentalPayments(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PlanRentalPayments", reflect.TypeOf((*MockRepo)(nil).PlanRentalPayments), arg0)
}

// RemovePreRental mocks base method.
func (m *MockRepo) RemovePreRental(arg0 context.Context, arg1 int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemovePreRental", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemovePreRental indicates an expected call of RemovePreRental.
func (mr *MockRepoMockRecorder) RemovePreRental(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemovePreRental", reflect.TypeOf((*MockRepo)(nil).RemovePreRental), arg0, arg1)
}

// UpdateContract mocks base method.
func (m *MockRepo) UpdateContract(arg0 context.Context, arg1 *dto.UpdateContract) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateContract", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateContract indicates an expected call of UpdateContract.
func (mr *MockRepoMockRecorder) UpdateContract(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateContract", reflect.TypeOf((*MockRepo)(nil).UpdateContract), arg0, arg1)
}

// UpdateContractContent mocks base method.
func (m *MockRepo) UpdateContractContent(arg0 context.Context, arg1 *dto.UpdateContractContent) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateContractContent", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateContractContent indicates an expected call of UpdateContractContent.
func (mr *MockRepoMockRecorder) UpdateContractContent(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateContractContent", reflect.TypeOf((*MockRepo)(nil).UpdateContractContent), arg0, arg1)
}

// UpdateFinePayments mocks base method.
func (m *MockRepo) UpdateFinePayments(arg0 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateFinePayments", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateFinePayments indicates an expected call of UpdateFinePayments.
func (mr *MockRepoMockRecorder) UpdateFinePayments(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateFinePayments", reflect.TypeOf((*MockRepo)(nil).UpdateFinePayments), arg0)
}

// UpdateFinePaymentsOfRental mocks base method.
func (m *MockRepo) UpdateFinePaymentsOfRental(arg0 context.Context, arg1 int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateFinePaymentsOfRental", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateFinePaymentsOfRental indicates an expected call of UpdateFinePaymentsOfRental.
func (mr *MockRepoMockRecorder) UpdateFinePaymentsOfRental(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateFinePaymentsOfRental", reflect.TypeOf((*MockRepo)(nil).UpdateFinePaymentsOfRental), arg0, arg1)
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

// UpdateRentalComplaint mocks base method.
func (m *MockRepo) UpdateRentalComplaint(arg0 context.Context, arg1 *dto.UpdateRentalComplaint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateRentalComplaint", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateRentalComplaint indicates an expected call of UpdateRentalComplaint.
func (mr *MockRepoMockRecorder) UpdateRentalComplaint(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateRentalComplaint", reflect.TypeOf((*MockRepo)(nil).UpdateRentalComplaint), arg0, arg1)
}

// UpdateRentalPayment mocks base method.
func (m *MockRepo) UpdateRentalPayment(arg0 context.Context, arg1 *dto.UpdateRentalPayment) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateRentalPayment", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateRentalPayment indicates an expected call of UpdateRentalPayment.
func (mr *MockRepoMockRecorder) UpdateRentalPayment(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateRentalPayment", reflect.TypeOf((*MockRepo)(nil).UpdateRentalPayment), arg0, arg1)
}

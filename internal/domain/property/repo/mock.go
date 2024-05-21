// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/user2410/rrms-backend/internal/domain/property/repo (interfaces: Repo)
//
// Generated by this command:
//
//	mockgen -package repo -destination internal/domain/property/repo/mock.go github.com/user2410/rrms-backend/internal/domain/property/repo Repo
//

// Package repo is a generated GoMock package.
package repo

import (
	context "context"
	reflect "reflect"

	uuid "github.com/google/uuid"
	dto "github.com/user2410/rrms-backend/internal/domain/application/dto"
	dto0 "github.com/user2410/rrms-backend/internal/domain/listing/dto"
	dto1 "github.com/user2410/rrms-backend/internal/domain/property/dto"
	model "github.com/user2410/rrms-backend/internal/domain/property/model"
	dto2 "github.com/user2410/rrms-backend/internal/domain/rental/dto"
	database "github.com/user2410/rrms-backend/internal/infrastructure/database"
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

// CreateProperty mocks base method.
func (m *MockRepo) CreateProperty(arg0 context.Context, arg1 *dto1.CreateProperty) (*model.PropertyModel, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateProperty", arg0, arg1)
	ret0, _ := ret[0].(*model.PropertyModel)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateProperty indicates an expected call of CreateProperty.
func (mr *MockRepoMockRecorder) CreateProperty(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateProperty", reflect.TypeOf((*MockRepo)(nil).CreateProperty), arg0, arg1)
}

// CreatePropertyManagerRequest mocks base method.
func (m *MockRepo) CreatePropertyManagerRequest(arg0 context.Context, arg1 *dto1.CreatePropertyManagerRequest) (model.NewPropertyManagerRequest, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreatePropertyManagerRequest", arg0, arg1)
	ret0, _ := ret[0].(model.NewPropertyManagerRequest)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreatePropertyManagerRequest indicates an expected call of CreatePropertyManagerRequest.
func (mr *MockRepoMockRecorder) CreatePropertyManagerRequest(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreatePropertyManagerRequest", reflect.TypeOf((*MockRepo)(nil).CreatePropertyManagerRequest), arg0, arg1)
}

// DeleteProperty mocks base method.
func (m *MockRepo) DeleteProperty(arg0 context.Context, arg1 uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteProperty", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteProperty indicates an expected call of DeleteProperty.
func (mr *MockRepoMockRecorder) DeleteProperty(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteProperty", reflect.TypeOf((*MockRepo)(nil).DeleteProperty), arg0, arg1)
}

// GetApplicationsOfProperty mocks base method.
func (m *MockRepo) GetApplicationsOfProperty(arg0 context.Context, arg1 uuid.UUID, arg2 *dto.GetApplicationsOfPropertyQuery) ([]int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetApplicationsOfProperty", arg0, arg1, arg2)
	ret0, _ := ret[0].([]int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetApplicationsOfProperty indicates an expected call of GetApplicationsOfProperty.
func (mr *MockRepoMockRecorder) GetApplicationsOfProperty(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetApplicationsOfProperty", reflect.TypeOf((*MockRepo)(nil).GetApplicationsOfProperty), arg0, arg1, arg2)
}

// GetListingsOfProperty mocks base method.
func (m *MockRepo) GetListingsOfProperty(arg0 context.Context, arg1 uuid.UUID, arg2 *dto0.GetListingsOfPropertyQuery) ([]uuid.UUID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetListingsOfProperty", arg0, arg1, arg2)
	ret0, _ := ret[0].([]uuid.UUID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetListingsOfProperty indicates an expected call of GetListingsOfProperty.
func (mr *MockRepoMockRecorder) GetListingsOfProperty(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetListingsOfProperty", reflect.TypeOf((*MockRepo)(nil).GetListingsOfProperty), arg0, arg1, arg2)
}

// GetManagedProperties mocks base method.
func (m *MockRepo) GetManagedProperties(arg0 context.Context, arg1 uuid.UUID) ([]database.GetManagedPropertiesRow, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetManagedProperties", arg0, arg1)
	ret0, _ := ret[0].([]database.GetManagedPropertiesRow)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetManagedProperties indicates an expected call of GetManagedProperties.
func (mr *MockRepoMockRecorder) GetManagedProperties(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetManagedProperties", reflect.TypeOf((*MockRepo)(nil).GetManagedProperties), arg0, arg1)
}

// GetNewPropertyManagerRequest mocks base method.
func (m *MockRepo) GetNewPropertyManagerRequest(arg0 context.Context, arg1 int64) (model.NewPropertyManagerRequest, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNewPropertyManagerRequest", arg0, arg1)
	ret0, _ := ret[0].(model.NewPropertyManagerRequest)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetNewPropertyManagerRequest indicates an expected call of GetNewPropertyManagerRequest.
func (mr *MockRepoMockRecorder) GetNewPropertyManagerRequest(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNewPropertyManagerRequest", reflect.TypeOf((*MockRepo)(nil).GetNewPropertyManagerRequest), arg0, arg1)
}

// GetPropertiesByIds mocks base method.
func (m *MockRepo) GetPropertiesByIds(arg0 context.Context, arg1, arg2 []string) ([]model.PropertyModel, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPropertiesByIds", arg0, arg1, arg2)
	ret0, _ := ret[0].([]model.PropertyModel)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPropertiesByIds indicates an expected call of GetPropertiesByIds.
func (mr *MockRepoMockRecorder) GetPropertiesByIds(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPropertiesByIds", reflect.TypeOf((*MockRepo)(nil).GetPropertiesByIds), arg0, arg1, arg2)
}

// GetPropertyById mocks base method.
func (m *MockRepo) GetPropertyById(arg0 context.Context, arg1 uuid.UUID) (*model.PropertyModel, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPropertyById", arg0, arg1)
	ret0, _ := ret[0].(*model.PropertyModel)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPropertyById indicates an expected call of GetPropertyById.
func (mr *MockRepoMockRecorder) GetPropertyById(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPropertyById", reflect.TypeOf((*MockRepo)(nil).GetPropertyById), arg0, arg1)
}

// GetPropertyManagers mocks base method.
func (m *MockRepo) GetPropertyManagers(arg0 context.Context, arg1 uuid.UUID) ([]model.PropertyManagerModel, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPropertyManagers", arg0, arg1)
	ret0, _ := ret[0].([]model.PropertyManagerModel)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPropertyManagers indicates an expected call of GetPropertyManagers.
func (mr *MockRepoMockRecorder) GetPropertyManagers(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPropertyManagers", reflect.TypeOf((*MockRepo)(nil).GetPropertyManagers), arg0, arg1)
}

// GetRentalsOfProperty mocks base method.
func (m *MockRepo) GetRentalsOfProperty(arg0 context.Context, arg1 uuid.UUID, arg2 *dto2.GetRentalsOfPropertyQuery) ([]int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRentalsOfProperty", arg0, arg1, arg2)
	ret0, _ := ret[0].([]int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRentalsOfProperty indicates an expected call of GetRentalsOfProperty.
func (mr *MockRepoMockRecorder) GetRentalsOfProperty(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRentalsOfProperty", reflect.TypeOf((*MockRepo)(nil).GetRentalsOfProperty), arg0, arg1, arg2)
}

// IsPublic mocks base method.
func (m *MockRepo) IsPublic(arg0 context.Context, arg1 uuid.UUID) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsPublic", arg0, arg1)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsPublic indicates an expected call of IsPublic.
func (mr *MockRepoMockRecorder) IsPublic(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsPublic", reflect.TypeOf((*MockRepo)(nil).IsPublic), arg0, arg1)
}

// SearchPropertyCombination mocks base method.
func (m *MockRepo) SearchPropertyCombination(arg0 context.Context, arg1 *dto1.SearchPropertyCombinationQuery) (*dto1.SearchPropertyCombinationResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SearchPropertyCombination", arg0, arg1)
	ret0, _ := ret[0].(*dto1.SearchPropertyCombinationResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SearchPropertyCombination indicates an expected call of SearchPropertyCombination.
func (mr *MockRepoMockRecorder) SearchPropertyCombination(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SearchPropertyCombination", reflect.TypeOf((*MockRepo)(nil).SearchPropertyCombination), arg0, arg1)
}

// UpdateProperty mocks base method.
func (m *MockRepo) UpdateProperty(arg0 context.Context, arg1 *dto1.UpdateProperty) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateProperty", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateProperty indicates an expected call of UpdateProperty.
func (mr *MockRepoMockRecorder) UpdateProperty(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateProperty", reflect.TypeOf((*MockRepo)(nil).UpdateProperty), arg0, arg1)
}

// UpdatePropertyManagerRequest mocks base method.
func (m *MockRepo) UpdatePropertyManagerRequest(arg0 context.Context, arg1 int64, arg2 uuid.UUID, arg3 bool) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdatePropertyManagerRequest", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdatePropertyManagerRequest indicates an expected call of UpdatePropertyManagerRequest.
func (mr *MockRepoMockRecorder) UpdatePropertyManagerRequest(arg0, arg1, arg2, arg3 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdatePropertyManagerRequest", reflect.TypeOf((*MockRepo)(nil).UpdatePropertyManagerRequest), arg0, arg1, arg2, arg3)
}

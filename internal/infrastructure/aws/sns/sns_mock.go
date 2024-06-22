// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/user2410/rrms-backend/internal/infrastructure/aws/sns (interfaces: SNSClient)
//
// Generated by this command:
//
//	mockgen -package sns -destination internal/infrastructure/aws/sns/sns_mock.go github.com/user2410/rrms-backend/internal/infrastructure/aws/sns SNSClient
//

// Package sns is a generated GoMock package.
package sns

import (
	context "context"
	reflect "reflect"

	sns "github.com/aws/aws-sdk-go-v2/service/sns"
	types "github.com/aws/aws-sdk-go-v2/service/sns/types"
	gomock "go.uber.org/mock/gomock"
)

// MockSNSClient is a mock of SNSClient interface.
type MockSNSClient struct {
	ctrl     *gomock.Controller
	recorder *MockSNSClientMockRecorder
}

// MockSNSClientMockRecorder is the mock recorder for MockSNSClient.
type MockSNSClientMockRecorder struct {
	mock *MockSNSClient
}

// NewMockSNSClient creates a new mock instance.
func NewMockSNSClient(ctrl *gomock.Controller) *MockSNSClient {
	mock := &MockSNSClient{ctrl: ctrl}
	mock.recorder = &MockSNSClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSNSClient) EXPECT() *MockSNSClientMockRecorder {
	return m.recorder
}

// CreateTopic mocks base method.
func (m *MockSNSClient) CreateTopic(arg0 context.Context, arg1 string, arg2, arg3 bool) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateTopic", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateTopic indicates an expected call of CreateTopic.
func (mr *MockSNSClientMockRecorder) CreateTopic(arg0, arg1, arg2, arg3 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTopic", reflect.TypeOf((*MockSNSClient)(nil).CreateTopic), arg0, arg1, arg2, arg3)
}

// DeleteTopic mocks base method.
func (m *MockSNSClient) DeleteTopic(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteTopic", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteTopic indicates an expected call of DeleteTopic.
func (mr *MockSNSClientMockRecorder) DeleteTopic(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteTopic", reflect.TypeOf((*MockSNSClient)(nil).DeleteTopic), arg0, arg1)
}

// GetClient mocks base method.
func (m *MockSNSClient) GetClient() *sns.Client {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetClient")
	ret0, _ := ret[0].(*sns.Client)
	return ret0
}

// GetClient indicates an expected call of GetClient.
func (mr *MockSNSClientMockRecorder) GetClient() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetClient", reflect.TypeOf((*MockSNSClient)(nil).GetClient))
}

// ListTopics mocks base method.
func (m *MockSNSClient) ListTopics(arg0 context.Context) ([]types.Topic, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListTopics", arg0)
	ret0, _ := ret[0].([]types.Topic)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListTopics indicates an expected call of ListTopics.
func (mr *MockSNSClientMockRecorder) ListTopics(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListTopics", reflect.TypeOf((*MockSNSClient)(nil).ListTopics), arg0)
}

// Publish mocks base method.
func (m *MockSNSClient) Publish(arg0 context.Context, arg1, arg2, arg3 string, arg4 map[string]any) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Publish", arg0, arg1, arg2, arg3, arg4)
	ret0, _ := ret[0].(error)
	return ret0
}

// Publish indicates an expected call of Publish.
func (mr *MockSNSClientMockRecorder) Publish(arg0, arg1, arg2, arg3, arg4 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Publish", reflect.TypeOf((*MockSNSClient)(nil).Publish), arg0, arg1, arg2, arg3, arg4)
}

// SubscribeQueue mocks base method.
func (m *MockSNSClient) SubscribeQueue(arg0, arg1 string, arg2 map[string][]string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SubscribeQueue", arg0, arg1, arg2)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SubscribeQueue indicates an expected call of SubscribeQueue.
func (mr *MockSNSClientMockRecorder) SubscribeQueue(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SubscribeQueue", reflect.TypeOf((*MockSNSClient)(nil).SubscribeQueue), arg0, arg1, arg2)
}
// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/sirjager/goth/worker (interfaces: TaskDistributor)
//
// Generated by this command:
//
//	mockgen -package mockTask -destination worker/mock/distributor.go github.com/sirjager/goth/worker TaskDistributor
//

// Package mockTask is a generated GoMock package.
package mockTask

import (
	context "context"
	reflect "reflect"

	asynq "github.com/hibiken/asynq"
	worker "github.com/sirjager/goth/worker"
	gomock "go.uber.org/mock/gomock"
)

// MockTaskDistributor is a mock of TaskDistributor interface.
type MockTaskDistributor struct {
	ctrl     *gomock.Controller
	recorder *MockTaskDistributorMockRecorder
}

// MockTaskDistributorMockRecorder is the mock recorder for MockTaskDistributor.
type MockTaskDistributorMockRecorder struct {
	mock *MockTaskDistributor
}

// NewMockTaskDistributor creates a new mock instance.
func NewMockTaskDistributor(ctrl *gomock.Controller) *MockTaskDistributor {
	mock := &MockTaskDistributor{ctrl: ctrl}
	mock.recorder = &MockTaskDistributorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTaskDistributor) EXPECT() *MockTaskDistributorMockRecorder {
	return m.recorder
}

// SendEmail mocks base method.
func (m *MockTaskDistributor) SendEmail(arg0 context.Context, arg1 worker.SendEmailParams, arg2 ...asynq.Option) error {
	m.ctrl.T.Helper()
	varargs := []any{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "SendEmail", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendEmail indicates an expected call of SendEmail.
func (mr *MockTaskDistributorMockRecorder) SendEmail(arg0, arg1 any, arg2 ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendEmail", reflect.TypeOf((*MockTaskDistributor)(nil).SendEmail), varargs...)
}

// Shutdown mocks base method.
func (m *MockTaskDistributor) Shutdown() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Shutdown")
}

// Shutdown indicates an expected call of Shutdown.
func (mr *MockTaskDistributorMockRecorder) Shutdown() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Shutdown", reflect.TypeOf((*MockTaskDistributor)(nil).Shutdown))
}
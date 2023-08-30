// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/openshift/assisted-service/pkg/generator (interfaces: ISOInstallConfigGenerator)

// Package generator is a generated GoMock package.
package generator

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	common "github.com/openshift/assisted-service/internal/common"
	events "github.com/openshift/assisted-service/internal/events"
)

// MockISOInstallConfigGenerator is a mock of ISOInstallConfigGenerator interface.
type MockISOInstallConfigGenerator struct {
	ctrl     *gomock.Controller
	recorder *MockISOInstallConfigGeneratorMockRecorder
}

// MockISOInstallConfigGeneratorMockRecorder is the mock recorder for MockISOInstallConfigGenerator.
type MockISOInstallConfigGeneratorMockRecorder struct {
	mock *MockISOInstallConfigGenerator
}

// NewMockISOInstallConfigGenerator creates a new mock instance.
func NewMockISOInstallConfigGenerator(ctrl *gomock.Controller) *MockISOInstallConfigGenerator {
	mock := &MockISOInstallConfigGenerator{ctrl: ctrl}
	mock.recorder = &MockISOInstallConfigGeneratorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockISOInstallConfigGenerator) EXPECT() *MockISOInstallConfigGeneratorMockRecorder {
	return m.recorder
}

// AbortInstallConfig mocks base method.
func (m *MockISOInstallConfigGenerator) AbortInstallConfig(arg0 context.Context, arg1 common.Cluster) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AbortInstallConfig", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// AbortInstallConfig indicates an expected call of AbortInstallConfig.
func (mr *MockISOInstallConfigGeneratorMockRecorder) AbortInstallConfig(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AbortInstallConfig", reflect.TypeOf((*MockISOInstallConfigGenerator)(nil).AbortInstallConfig), arg0, arg1)
}

// GenerateISO mocks base method.
func (m *MockISOInstallConfigGenerator) GenerateISO(arg0 context.Context, arg1 common.Cluster, arg2, arg3, arg4 string, arg5 events.Handler) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateISO", arg0, arg1, arg2, arg3, arg4, arg5)
	ret0, _ := ret[0].(error)
	return ret0
}

// GenerateISO indicates an expected call of GenerateISO.
func (mr *MockISOInstallConfigGeneratorMockRecorder) GenerateISO(arg0, arg1, arg2, arg3, arg4, arg5 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateISO", reflect.TypeOf((*MockISOInstallConfigGenerator)(nil).GenerateISO), arg0, arg1, arg2, arg3, arg4, arg5)
}

// GenerateInstallConfig mocks base method.
func (m *MockISOInstallConfigGenerator) GenerateInstallConfig(arg0 context.Context, arg1 common.Cluster, arg2 []byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateInstallConfig", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// GenerateInstallConfig indicates an expected call of GenerateInstallConfig.
func (mr *MockISOInstallConfigGeneratorMockRecorder) GenerateInstallConfig(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateInstallConfig", reflect.TypeOf((*MockISOInstallConfigGenerator)(nil).GenerateInstallConfig), arg0, arg1, arg2)
}

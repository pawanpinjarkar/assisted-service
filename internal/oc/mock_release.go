// Code generated by MockGen. DO NOT EDIT.
// Source: release.go

// Package oc is a generated GoMock package.
package oc

import (
	gomock "github.com/golang/mock/gomock"
	logrus "github.com/sirupsen/logrus"
	reflect "reflect"
)

// MockRelease is a mock of Release interface
type MockRelease struct {
	ctrl     *gomock.Controller
	recorder *MockReleaseMockRecorder
}

// MockReleaseMockRecorder is the mock recorder for MockRelease
type MockReleaseMockRecorder struct {
	mock *MockRelease
}

// NewMockRelease creates a new mock instance
func NewMockRelease(ctrl *gomock.Controller) *MockRelease {
	mock := &MockRelease{ctrl: ctrl}
	mock.recorder = &MockReleaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockRelease) EXPECT() *MockReleaseMockRecorder {
	return m.recorder
}

// GetMCOImage mocks base method
func (m *MockRelease) GetMCOImage(log logrus.FieldLogger, releaseImage, releaseImageMirror, pullSecret string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMCOImage", log, releaseImage, releaseImageMirror, pullSecret)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMCOImage indicates an expected call of GetMCOImage
func (mr *MockReleaseMockRecorder) GetMCOImage(log, releaseImage, releaseImageMirror, pullSecret interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMCOImage", reflect.TypeOf((*MockRelease)(nil).GetMCOImage), log, releaseImage, releaseImageMirror, pullSecret)
}

// GetMustGatherImage mocks base method
func (m *MockRelease) GetMustGatherImage(log logrus.FieldLogger, releaseImage, releaseImageMirror, pullSecret string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMustGatherImage", log, releaseImage, releaseImageMirror, pullSecret)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMustGatherImage indicates an expected call of GetMustGatherImage
func (mr *MockReleaseMockRecorder) GetMustGatherImage(log, releaseImage, releaseImageMirror, pullSecret interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMustGatherImage", reflect.TypeOf((*MockRelease)(nil).GetMustGatherImage), log, releaseImage, releaseImageMirror, pullSecret)
}

// GetOpenshiftVersion mocks base method
func (m *MockRelease) GetOpenshiftVersion(log logrus.FieldLogger, releaseImage, releaseImageMirror, pullSecret string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOpenshiftVersion", log, releaseImage, releaseImageMirror, pullSecret)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOpenshiftVersion indicates an expected call of GetOpenshiftVersion
func (mr *MockReleaseMockRecorder) GetOpenshiftVersion(log, releaseImage, releaseImageMirror, pullSecret interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOpenshiftVersion", reflect.TypeOf((*MockRelease)(nil).GetOpenshiftVersion), log, releaseImage, releaseImageMirror, pullSecret)
}

// GetMajorMinorVersion mocks base method
func (m *MockRelease) GetMajorMinorVersion(log logrus.FieldLogger, releaseImage, releaseImageMirror, pullSecret string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMajorMinorVersion", log, releaseImage, releaseImageMirror, pullSecret)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMajorMinorVersion indicates an expected call of GetMajorMinorVersion
func (mr *MockReleaseMockRecorder) GetMajorMinorVersion(log, releaseImage, releaseImageMirror, pullSecret interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMajorMinorVersion", reflect.TypeOf((*MockRelease)(nil).GetMajorMinorVersion), log, releaseImage, releaseImageMirror, pullSecret)
}

// Extract mocks base method
func (m *MockRelease) Extract(log logrus.FieldLogger, releaseImage, releaseImageMirror, cacheDir, pullSecret string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Extract", log, releaseImage, releaseImageMirror, cacheDir, pullSecret)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Extract indicates an expected call of Extract
func (mr *MockReleaseMockRecorder) Extract(log, releaseImage, releaseImageMirror, cacheDir, pullSecret interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Extract", reflect.TypeOf((*MockRelease)(nil).Extract), log, releaseImage, releaseImageMirror, cacheDir, pullSecret)
}

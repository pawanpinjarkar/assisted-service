// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/openshift/assisted-service/internal/controller/controllers (interfaces: SpokeClientCache)

// Package controllers is a generated GoMock package.
package controllers

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	spoke_k8s_client "github.com/openshift/assisted-service/internal/spoke_k8s_client"
	v1 "k8s.io/api/core/v1"
)

// MockSpokeClientCache is a mock of SpokeClientCache interface.
type MockSpokeClientCache struct {
	ctrl     *gomock.Controller
	recorder *MockSpokeClientCacheMockRecorder
}

// MockSpokeClientCacheMockRecorder is the mock recorder for MockSpokeClientCache.
type MockSpokeClientCacheMockRecorder struct {
	mock *MockSpokeClientCache
}

// NewMockSpokeClientCache creates a new mock instance.
func NewMockSpokeClientCache(ctrl *gomock.Controller) *MockSpokeClientCache {
	mock := &MockSpokeClientCache{ctrl: ctrl}
	mock.recorder = &MockSpokeClientCacheMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSpokeClientCache) EXPECT() *MockSpokeClientCacheMockRecorder {
	return m.recorder
}

// Get mocks base method.
func (m *MockSpokeClientCache) Get(arg0 *v1.Secret) (spoke_k8s_client.SpokeK8sClient, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0)
	ret0, _ := ret[0].(spoke_k8s_client.SpokeK8sClient)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockSpokeClientCacheMockRecorder) Get(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockSpokeClientCache)(nil).Get), arg0)
}

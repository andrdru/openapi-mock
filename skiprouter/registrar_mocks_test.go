// Code generated by MockGen. DO NOT EDIT.
// Source: registrar.go
//
// Generated by this command:
//
//	mockgen -source=registrar.go -destination=registrar_mocks_test.go -package=skiprouter_test
//

// Package skiprouter_test is a generated GoMock package.
package skiprouter_test

import (
	http "net/http"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// Mockrouter is a mock of router interface.
type Mockrouter struct {
	ctrl     *gomock.Controller
	recorder *MockrouterMockRecorder
	isgomock struct{}
}

// MockrouterMockRecorder is the mock recorder for Mockrouter.
type MockrouterMockRecorder struct {
	mock *Mockrouter
}

// NewMockrouter creates a new mock instance.
func NewMockrouter(ctrl *gomock.Controller) *Mockrouter {
	mock := &Mockrouter{ctrl: ctrl}
	mock.recorder = &MockrouterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Mockrouter) EXPECT() *MockrouterMockRecorder {
	return m.recorder
}

// HandleFunc mocks base method.
func (m *Mockrouter) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "HandleFunc", pattern, handler)
}

// HandleFunc indicates an expected call of HandleFunc.
func (mr *MockrouterMockRecorder) HandleFunc(pattern, handler any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HandleFunc", reflect.TypeOf((*Mockrouter)(nil).HandleFunc), pattern, handler)
}

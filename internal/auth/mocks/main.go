// Code generated by MockGen. DO NOT EDIT.
// Source: main.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
)

// MockclockI is a mock of clockI interface.
type MockclockI struct {
	ctrl     *gomock.Controller
	recorder *MockclockIMockRecorder
}

// MockclockIMockRecorder is the mock recorder for MockclockI.
type MockclockIMockRecorder struct {
	mock *MockclockI
}

// NewMockclockI creates a new mock instance.
func NewMockclockI(ctrl *gomock.Controller) *MockclockI {
	mock := &MockclockI{ctrl: ctrl}
	mock.recorder = &MockclockIMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockclockI) EXPECT() *MockclockIMockRecorder {
	return m.recorder
}

// Now mocks base method.
func (m *MockclockI) Now() time.Time {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Now")
	ret0, _ := ret[0].(time.Time)
	return ret0
}

// Now indicates an expected call of Now.
func (mr *MockclockIMockRecorder) Now() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Now", reflect.TypeOf((*MockclockI)(nil).Now))
}

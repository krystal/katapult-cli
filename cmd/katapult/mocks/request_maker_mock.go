// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/krystal/go-katapult/core (interfaces: RequestMaker)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	katapult "github.com/krystal/go-katapult"
)

// MockRequestMaker is a mock of RequestMaker interface.
type MockRequestMaker struct {
	ctrl     *gomock.Controller
	recorder *MockRequestMakerMockRecorder
}

// MockRequestMakerMockRecorder is the mock recorder for MockRequestMaker.
type MockRequestMakerMockRecorder struct {
	mock *MockRequestMaker
}

// NewMockRequestMaker creates a new mock instance.
func NewMockRequestMaker(ctrl *gomock.Controller) *MockRequestMaker {
	mock := &MockRequestMaker{ctrl: ctrl}
	mock.recorder = &MockRequestMakerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRequestMaker) EXPECT() *MockRequestMakerMockRecorder {
	return m.recorder
}

// Do mocks base method.
func (m *MockRequestMaker) Do(arg0 context.Context, arg1 *katapult.Request, arg2 interface{}) (*katapult.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Do", arg0, arg1, arg2)
	ret0, _ := ret[0].(*katapult.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Do indicates an expected call of Do.
func (mr *MockRequestMakerMockRecorder) Do(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Do", reflect.TypeOf((*MockRequestMaker)(nil).Do), arg0, arg1, arg2)
}

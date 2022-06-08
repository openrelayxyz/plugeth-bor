// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/ethereum/go-ethereum/consensus/bor (interfaces: IHeimdallClient)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	bor "github.com/ethereum/go-ethereum/consensus/bor"
	gomock "github.com/golang/mock/gomock"
)

// MockIHeimdallClient is a mock of IHeimdallClient interface.
type MockIHeimdallClient struct {
	ctrl     *gomock.Controller
	recorder *MockIHeimdallClientMockRecorder
}

// MockIHeimdallClientMockRecorder is the mock recorder for MockIHeimdallClient.
type MockIHeimdallClientMockRecorder struct {
	mock *MockIHeimdallClient
}

// NewMockIHeimdallClient creates a new mock instance.
func NewMockIHeimdallClient(ctrl *gomock.Controller) *MockIHeimdallClient {
	mock := &MockIHeimdallClient{ctrl: ctrl}
	mock.recorder = &MockIHeimdallClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIHeimdallClient) EXPECT() *MockIHeimdallClientMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockIHeimdallClient) Close() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Close")
}

// Close indicates an expected call of Close.
func (mr *MockIHeimdallClientMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockIHeimdallClient)(nil).Close))
}

// Fetch mocks base method.
func (m *MockIHeimdallClient) Fetch(arg0, arg1 string) (*bor.ResponseWithHeight, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Fetch", arg0, arg1)
	ret0, _ := ret[0].(*bor.ResponseWithHeight)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Fetch indicates an expected call of Fetch.
func (mr *MockIHeimdallClientMockRecorder) Fetch(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Fetch", reflect.TypeOf((*MockIHeimdallClient)(nil).Fetch), arg0, arg1)
}

// FetchLatestCheckpoint mocks base method.
func (m *MockIHeimdallClient) FetchLatestCheckpoint() (*bor.Checkpoint, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FetchLatestCheckpoint")
	ret0, _ := ret[0].(*bor.Checkpoint)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FetchLatestCheckpoint indicates an expected call of FetchLatestCheckpoint.
func (mr *MockIHeimdallClientMockRecorder) FetchLatestCheckpoint() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FetchLatestCheckpoint", reflect.TypeOf((*MockIHeimdallClient)(nil).FetchLatestCheckpoint))
}

// FetchStateSyncEvents mocks base method.
func (m *MockIHeimdallClient) FetchStateSyncEvents(arg0 uint64, arg1 int64) ([]*bor.EventRecordWithTime, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FetchStateSyncEvents", arg0, arg1)
	ret0, _ := ret[0].([]*bor.EventRecordWithTime)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FetchStateSyncEvents indicates an expected call of FetchStateSyncEvents.
func (mr *MockIHeimdallClientMockRecorder) FetchStateSyncEvents(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FetchStateSyncEvents", reflect.TypeOf((*MockIHeimdallClient)(nil).FetchStateSyncEvents), arg0, arg1)
}

// FetchWithRetry mocks base method.
func (m *MockIHeimdallClient) FetchWithRetry(arg0, arg1 string) (*bor.ResponseWithHeight, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FetchWithRetry", arg0, arg1)
	ret0, _ := ret[0].(*bor.ResponseWithHeight)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FetchWithRetry indicates an expected call of FetchWithRetry.
func (mr *MockIHeimdallClientMockRecorder) FetchWithRetry(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FetchWithRetry", reflect.TypeOf((*MockIHeimdallClient)(nil).FetchWithRetry), arg0, arg1)
}

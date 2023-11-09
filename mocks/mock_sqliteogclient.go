// Code generated by MockGen. DO NOT EDIT.
// Source: gitlab.com/aous.omr/sqlite-og/gen/proto (interfaces: SqliteOGClient)

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	sqlite_og "gitlab.com/aous.omr/sqlite-og/gen/proto"
	gomock "go.uber.org/mock/gomock"
	grpc "google.golang.org/grpc"
)

// MockSqliteOGClient is a mock of SqliteOGClient interface.
type MockSqliteOGClient struct {
	ctrl     *gomock.Controller
	recorder *MockSqliteOGClientMockRecorder
}

// MockSqliteOGClientMockRecorder is the mock recorder for MockSqliteOGClient.
type MockSqliteOGClientMockRecorder struct {
	mock *MockSqliteOGClient
}

// NewMockSqliteOGClient creates a new mock instance.
func NewMockSqliteOGClient(ctrl *gomock.Controller) *MockSqliteOGClient {
	mock := &MockSqliteOGClient{ctrl: ctrl}
	mock.recorder = &MockSqliteOGClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSqliteOGClient) EXPECT() *MockSqliteOGClientMockRecorder {
	return m.recorder
}

// Callback mocks base method.
func (m *MockSqliteOGClient) Callback(arg0 context.Context, arg1 ...grpc.CallOption) (sqlite_og.SqliteOG_CallbackClient, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Callback", varargs...)
	ret0, _ := ret[0].(sqlite_og.SqliteOG_CallbackClient)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Callback indicates an expected call of Callback.
func (mr *MockSqliteOGClientMockRecorder) Callback(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Callback", reflect.TypeOf((*MockSqliteOGClient)(nil).Callback), varargs...)
}

// Close mocks base method.
func (m *MockSqliteOGClient) Close(arg0 context.Context, arg1 *sqlite_og.ConnectionId, arg2 ...grpc.CallOption) (*sqlite_og.Empty, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Close", varargs...)
	ret0, _ := ret[0].(*sqlite_og.Empty)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Close indicates an expected call of Close.
func (mr *MockSqliteOGClientMockRecorder) Close(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockSqliteOGClient)(nil).Close), varargs...)
}

// Connection mocks base method.
func (m *MockSqliteOGClient) Connection(arg0 context.Context, arg1 *sqlite_og.ConnectionRequest, arg2 ...grpc.CallOption) (*sqlite_og.ConnectionId, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Connection", varargs...)
	ret0, _ := ret[0].(*sqlite_og.ConnectionId)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Connection indicates an expected call of Connection.
func (mr *MockSqliteOGClientMockRecorder) Connection(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Connection", reflect.TypeOf((*MockSqliteOGClient)(nil).Connection), varargs...)
}

// Execute mocks base method.
func (m *MockSqliteOGClient) Execute(arg0 context.Context, arg1 *sqlite_og.Statement, arg2 ...grpc.CallOption) (*sqlite_og.ExecuteResult, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Execute", varargs...)
	ret0, _ := ret[0].(*sqlite_og.ExecuteResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Execute indicates an expected call of Execute.
func (mr *MockSqliteOGClientMockRecorder) Execute(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Execute", reflect.TypeOf((*MockSqliteOGClient)(nil).Execute), varargs...)
}

// ExecuteOrQuery mocks base method.
func (m *MockSqliteOGClient) ExecuteOrQuery(arg0 context.Context, arg1 *sqlite_og.Statement, arg2 ...grpc.CallOption) (*sqlite_og.ExecuteOrQueryResult, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ExecuteOrQuery", varargs...)
	ret0, _ := ret[0].(*sqlite_og.ExecuteOrQueryResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ExecuteOrQuery indicates an expected call of ExecuteOrQuery.
func (mr *MockSqliteOGClientMockRecorder) ExecuteOrQuery(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ExecuteOrQuery", reflect.TypeOf((*MockSqliteOGClient)(nil).ExecuteOrQuery), varargs...)
}

// IsValid mocks base method.
func (m *MockSqliteOGClient) IsValid(arg0 context.Context, arg1 *sqlite_og.ConnectionId, arg2 ...grpc.CallOption) (*sqlite_og.Empty, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "IsValid", varargs...)
	ret0, _ := ret[0].(*sqlite_og.Empty)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsValid indicates an expected call of IsValid.
func (mr *MockSqliteOGClientMockRecorder) IsValid(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsValid", reflect.TypeOf((*MockSqliteOGClient)(nil).IsValid), varargs...)
}

// Ping mocks base method.
func (m *MockSqliteOGClient) Ping(arg0 context.Context, arg1 *sqlite_og.Empty, arg2 ...grpc.CallOption) (*sqlite_og.Empty, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Ping", varargs...)
	ret0, _ := ret[0].(*sqlite_og.Empty)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Ping indicates an expected call of Ping.
func (mr *MockSqliteOGClientMockRecorder) Ping(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Ping", reflect.TypeOf((*MockSqliteOGClient)(nil).Ping), varargs...)
}

// Query mocks base method.
func (m *MockSqliteOGClient) Query(arg0 context.Context, arg1 *sqlite_og.Statement, arg2 ...grpc.CallOption) (*sqlite_og.QueryResult, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Query", varargs...)
	ret0, _ := ret[0].(*sqlite_og.QueryResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Query indicates an expected call of Query.
func (mr *MockSqliteOGClientMockRecorder) Query(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Query", reflect.TypeOf((*MockSqliteOGClient)(nil).Query), varargs...)
}

// ResetSession mocks base method.
func (m *MockSqliteOGClient) ResetSession(arg0 context.Context, arg1 *sqlite_og.ConnectionId, arg2 ...grpc.CallOption) (*sqlite_og.ConnectionId, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ResetSession", varargs...)
	ret0, _ := ret[0].(*sqlite_og.ConnectionId)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ResetSession indicates an expected call of ResetSession.
func (mr *MockSqliteOGClientMockRecorder) ResetSession(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ResetSession", reflect.TypeOf((*MockSqliteOGClient)(nil).ResetSession), varargs...)
}

// Code generated by MockGen. DO NOT EDIT.
// Source: storage.go

// Package mock_storage is a generated GoMock package.
package storage

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	api "github.com/gsk148/urlShorteningService/internal/app/api"
)

// MockStorage is a mock of Storage interface.
type MockStorage struct {
	ctrl     *gomock.Controller
	recorder *MockStorageMockRecorder
}

// MockStorageMockRecorder is the mock recorder for MockStorage.
type MockStorageMockRecorder struct {
	mock *MockStorage
}

// NewMockStorage creates a new mock instance.
func NewMockStorage(ctrl *gomock.Controller) *MockStorage {
	mock := &MockStorage{ctrl: ctrl}
	mock.recorder = &MockStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStorage) EXPECT() *MockStorageMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockStorage) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockStorageMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockStorage)(nil).Close))
}

// DeleteByUserIDAndShort mocks base method.
func (m *MockStorage) DeleteByUserIDAndShort(userID, shortURL string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteByUserIDAndShort", userID, shortURL)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteByUserIDAndShort indicates an expected call of DeleteByUserIDAndShort.
func (mr *MockStorageMockRecorder) DeleteByUserIDAndShort(userID, shortURL interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteByUserIDAndShort", reflect.TypeOf((*MockStorage)(nil).DeleteByUserIDAndShort), userID, shortURL)
}

// Get mocks base method.
func (m *MockStorage) Get(key string) (api.ShortenedData, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", key)
	ret0, _ := ret[0].(api.ShortenedData)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockStorageMockRecorder) Get(key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockStorage)(nil).Get), key)
}

// GetBatchByUserID mocks base method.
func (m *MockStorage) GetBatchByUserID(userID string) ([]api.ShortenedData, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBatchByUserID", userID)
	ret0, _ := ret[0].([]api.ShortenedData)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBatchByUserID indicates an expected call of GetBatchByUserID.
func (mr *MockStorageMockRecorder) GetBatchByUserID(userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBatchByUserID", reflect.TypeOf((*MockStorage)(nil).GetBatchByUserID), userID)
}

// GetStatistic mocks base method.
func (m *MockStorage) GetStatistic() *api.Statistic {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStatistic")
	ret0, _ := ret[0].(*api.Statistic)
	return ret0
}

// GetStatistic indicates an expected call of GetStatistic.
func (mr *MockStorageMockRecorder) GetStatistic() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStatistic", reflect.TypeOf((*MockStorage)(nil).GetStatistic))
}

// Ping mocks base method.
func (m *MockStorage) Ping() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Ping")
	ret0, _ := ret[0].(error)
	return ret0
}

// Ping indicates an expected call of Ping.
func (mr *MockStorageMockRecorder) Ping() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Ping", reflect.TypeOf((*MockStorage)(nil).Ping))
}

// Store mocks base method.
func (m *MockStorage) Store(data api.ShortenedData) (api.ShortenedData, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Store", data)
	ret0, _ := ret[0].(api.ShortenedData)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Store indicates an expected call of Store.
func (mr *MockStorageMockRecorder) Store(data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Store", reflect.TypeOf((*MockStorage)(nil).Store), data)
}

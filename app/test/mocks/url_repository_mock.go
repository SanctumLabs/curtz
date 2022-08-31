// Code generated by MockGen. DO NOT EDIT.
// Source: app/internal/core/contracts/url_repository.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	entities "github.com/sanctumlabs/curtz/app/internal/core/entities"
)

// MockUrlRepository is a mock of UrlRepository interface.
type MockUrlRepository struct {
	ctrl     *gomock.Controller
	recorder *MockUrlRepositoryMockRecorder
}

// MockUrlRepositoryMockRecorder is the mock recorder for MockUrlRepository.
type MockUrlRepositoryMockRecorder struct {
	mock *MockUrlRepository
}

// NewMockUrlRepository creates a new mock instance.
func NewMockUrlRepository(ctrl *gomock.Controller) *MockUrlRepository {
	mock := &MockUrlRepository{ctrl: ctrl}
	mock.recorder = &MockUrlRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUrlRepository) EXPECT() *MockUrlRepositoryMockRecorder {
	return m.recorder
}

// Delete mocks base method.
func (m *MockUrlRepository) Delete(id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockUrlRepositoryMockRecorder) Delete(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockUrlRepository)(nil).Delete), id)
}

// GetById mocks base method.
func (m *MockUrlRepository) GetById(id string) (entities.URL, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetById", id)
	ret0, _ := ret[0].(entities.URL)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetById indicates an expected call of GetById.
func (mr *MockUrlRepositoryMockRecorder) GetById(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetById", reflect.TypeOf((*MockUrlRepository)(nil).GetById), id)
}

// GetByKeyword mocks base method.
func (m *MockUrlRepository) GetByKeyword(keyword string) ([]entities.URL, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByKeyword", keyword)
	ret0, _ := ret[0].([]entities.URL)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByKeyword indicates an expected call of GetByKeyword.
func (mr *MockUrlRepositoryMockRecorder) GetByKeyword(keyword interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByKeyword", reflect.TypeOf((*MockUrlRepository)(nil).GetByKeyword), keyword)
}

// GetByKeywords mocks base method.
func (m *MockUrlRepository) GetByKeywords(keywords []string) ([]entities.URL, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByKeywords", keywords)
	ret0, _ := ret[0].([]entities.URL)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByKeywords indicates an expected call of GetByKeywords.
func (mr *MockUrlRepositoryMockRecorder) GetByKeywords(keywords interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByKeywords", reflect.TypeOf((*MockUrlRepository)(nil).GetByKeywords), keywords)
}

// GetByOriginalUrl mocks base method.
func (m *MockUrlRepository) GetByOriginalUrl(originalUrl string) (entities.URL, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByOriginalUrl", originalUrl)
	ret0, _ := ret[0].(entities.URL)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByOriginalUrl indicates an expected call of GetByOriginalUrl.
func (mr *MockUrlRepositoryMockRecorder) GetByOriginalUrl(originalUrl interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByOriginalUrl", reflect.TypeOf((*MockUrlRepository)(nil).GetByOriginalUrl), originalUrl)
}

// GetByOwner mocks base method.
func (m *MockUrlRepository) GetByOwner(owner string) ([]entities.URL, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByOwner", owner)
	ret0, _ := ret[0].([]entities.URL)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByOwner indicates an expected call of GetByOwner.
func (mr *MockUrlRepositoryMockRecorder) GetByOwner(owner interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByOwner", reflect.TypeOf((*MockUrlRepository)(nil).GetByOwner), owner)
}

// GetByShortCode mocks base method.
func (m *MockUrlRepository) GetByShortCode(shortCode string) (entities.URL, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByShortCode", shortCode)
	ret0, _ := ret[0].(entities.URL)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByShortCode indicates an expected call of GetByShortCode.
func (mr *MockUrlRepositoryMockRecorder) GetByShortCode(shortCode interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByShortCode", reflect.TypeOf((*MockUrlRepository)(nil).GetByShortCode), shortCode)
}

// IncrementHits mocks base method.
func (m *MockUrlRepository) IncrementHits(shortCode string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IncrementHits", shortCode)
	ret0, _ := ret[0].(error)
	return ret0
}

// IncrementHits indicates an expected call of IncrementHits.
func (mr *MockUrlRepositoryMockRecorder) IncrementHits(shortCode interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IncrementHits", reflect.TypeOf((*MockUrlRepository)(nil).IncrementHits), shortCode)
}

// Save mocks base method.
func (m *MockUrlRepository) Save(arg0 entities.URL) (entities.URL, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", arg0)
	ret0, _ := ret[0].(entities.URL)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Save indicates an expected call of Save.
func (mr *MockUrlRepositoryMockRecorder) Save(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockUrlRepository)(nil).Save), arg0)
}

// MockUrlWriteRepository is a mock of UrlWriteRepository interface.
type MockUrlWriteRepository struct {
	ctrl     *gomock.Controller
	recorder *MockUrlWriteRepositoryMockRecorder
}

// MockUrlWriteRepositoryMockRecorder is the mock recorder for MockUrlWriteRepository.
type MockUrlWriteRepositoryMockRecorder struct {
	mock *MockUrlWriteRepository
}

// NewMockUrlWriteRepository creates a new mock instance.
func NewMockUrlWriteRepository(ctrl *gomock.Controller) *MockUrlWriteRepository {
	mock := &MockUrlWriteRepository{ctrl: ctrl}
	mock.recorder = &MockUrlWriteRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUrlWriteRepository) EXPECT() *MockUrlWriteRepositoryMockRecorder {
	return m.recorder
}

// Delete mocks base method.
func (m *MockUrlWriteRepository) Delete(id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockUrlWriteRepositoryMockRecorder) Delete(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockUrlWriteRepository)(nil).Delete), id)
}

// IncrementHits mocks base method.
func (m *MockUrlWriteRepository) IncrementHits(shortCode string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IncrementHits", shortCode)
	ret0, _ := ret[0].(error)
	return ret0
}

// IncrementHits indicates an expected call of IncrementHits.
func (mr *MockUrlWriteRepositoryMockRecorder) IncrementHits(shortCode interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IncrementHits", reflect.TypeOf((*MockUrlWriteRepository)(nil).IncrementHits), shortCode)
}

// Save mocks base method.
func (m *MockUrlWriteRepository) Save(arg0 entities.URL) (entities.URL, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", arg0)
	ret0, _ := ret[0].(entities.URL)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Save indicates an expected call of Save.
func (mr *MockUrlWriteRepositoryMockRecorder) Save(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockUrlWriteRepository)(nil).Save), arg0)
}

// MockUrlReadRepository is a mock of UrlReadRepository interface.
type MockUrlReadRepository struct {
	ctrl     *gomock.Controller
	recorder *MockUrlReadRepositoryMockRecorder
}

// MockUrlReadRepositoryMockRecorder is the mock recorder for MockUrlReadRepository.
type MockUrlReadRepositoryMockRecorder struct {
	mock *MockUrlReadRepository
}

// NewMockUrlReadRepository creates a new mock instance.
func NewMockUrlReadRepository(ctrl *gomock.Controller) *MockUrlReadRepository {
	mock := &MockUrlReadRepository{ctrl: ctrl}
	mock.recorder = &MockUrlReadRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUrlReadRepository) EXPECT() *MockUrlReadRepositoryMockRecorder {
	return m.recorder
}

// GetById mocks base method.
func (m *MockUrlReadRepository) GetById(id string) (entities.URL, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetById", id)
	ret0, _ := ret[0].(entities.URL)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetById indicates an expected call of GetById.
func (mr *MockUrlReadRepositoryMockRecorder) GetById(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetById", reflect.TypeOf((*MockUrlReadRepository)(nil).GetById), id)
}

// GetByKeyword mocks base method.
func (m *MockUrlReadRepository) GetByKeyword(keyword string) ([]entities.URL, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByKeyword", keyword)
	ret0, _ := ret[0].([]entities.URL)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByKeyword indicates an expected call of GetByKeyword.
func (mr *MockUrlReadRepositoryMockRecorder) GetByKeyword(keyword interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByKeyword", reflect.TypeOf((*MockUrlReadRepository)(nil).GetByKeyword), keyword)
}

// GetByKeywords mocks base method.
func (m *MockUrlReadRepository) GetByKeywords(keywords []string) ([]entities.URL, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByKeywords", keywords)
	ret0, _ := ret[0].([]entities.URL)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByKeywords indicates an expected call of GetByKeywords.
func (mr *MockUrlReadRepositoryMockRecorder) GetByKeywords(keywords interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByKeywords", reflect.TypeOf((*MockUrlReadRepository)(nil).GetByKeywords), keywords)
}

// GetByOriginalUrl mocks base method.
func (m *MockUrlReadRepository) GetByOriginalUrl(originalUrl string) (entities.URL, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByOriginalUrl", originalUrl)
	ret0, _ := ret[0].(entities.URL)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByOriginalUrl indicates an expected call of GetByOriginalUrl.
func (mr *MockUrlReadRepositoryMockRecorder) GetByOriginalUrl(originalUrl interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByOriginalUrl", reflect.TypeOf((*MockUrlReadRepository)(nil).GetByOriginalUrl), originalUrl)
}

// GetByOwner mocks base method.
func (m *MockUrlReadRepository) GetByOwner(owner string) ([]entities.URL, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByOwner", owner)
	ret0, _ := ret[0].([]entities.URL)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByOwner indicates an expected call of GetByOwner.
func (mr *MockUrlReadRepositoryMockRecorder) GetByOwner(owner interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByOwner", reflect.TypeOf((*MockUrlReadRepository)(nil).GetByOwner), owner)
}

// GetByShortCode mocks base method.
func (m *MockUrlReadRepository) GetByShortCode(shortCode string) (entities.URL, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByShortCode", shortCode)
	ret0, _ := ret[0].(entities.URL)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByShortCode indicates an expected call of GetByShortCode.
func (mr *MockUrlReadRepositoryMockRecorder) GetByShortCode(shortCode interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByShortCode", reflect.TypeOf((*MockUrlReadRepository)(nil).GetByShortCode), shortCode)
}

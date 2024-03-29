// Code generated by MockGen. DO NOT EDIT.
// Source: app/internal/core/contracts/services.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
)

// MockAuthService is a mock of AuthService interface.
type MockAuthService struct {
	ctrl     *gomock.Controller
	recorder *MockAuthServiceMockRecorder
}

// MockAuthServiceMockRecorder is the mock recorder for MockAuthService.
type MockAuthServiceMockRecorder struct {
	mock *MockAuthService
}

// NewMockAuthService creates a new mock instance.
func NewMockAuthService(ctrl *gomock.Controller) *MockAuthService {
	mock := &MockAuthService{ctrl: ctrl}
	mock.recorder = &MockAuthServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuthService) EXPECT() *MockAuthServiceMockRecorder {
	return m.recorder
}

// Authenticate mocks base method.
func (m *MockAuthService) Authenticate(token string) (string, time.Time, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Authenticate", token)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(time.Time)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Authenticate indicates an expected call of Authenticate.
func (mr *MockAuthServiceMockRecorder) Authenticate(token interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Authenticate", reflect.TypeOf((*MockAuthService)(nil).Authenticate), token)
}

// GenerateRefreshToken mocks base method.
func (m *MockAuthService) GenerateRefreshToken(userID string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateRefreshToken", userID)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GenerateRefreshToken indicates an expected call of GenerateRefreshToken.
func (mr *MockAuthServiceMockRecorder) GenerateRefreshToken(userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateRefreshToken", reflect.TypeOf((*MockAuthService)(nil).GenerateRefreshToken), userID)
}

// GenerateToken mocks base method.
func (m *MockAuthService) GenerateToken(userID string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateToken", userID)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GenerateToken indicates an expected call of GenerateToken.
func (mr *MockAuthServiceMockRecorder) GenerateToken(userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateToken", reflect.TypeOf((*MockAuthService)(nil).GenerateToken), userID)
}

// MockNotificationService is a mock of NotificationService interface.
type MockNotificationService struct {
	ctrl     *gomock.Controller
	recorder *MockNotificationServiceMockRecorder
}

// MockNotificationServiceMockRecorder is the mock recorder for MockNotificationService.
type MockNotificationServiceMockRecorder struct {
	mock *MockNotificationService
}

// NewMockNotificationService creates a new mock instance.
func NewMockNotificationService(ctrl *gomock.Controller) *MockNotificationService {
	mock := &MockNotificationService{ctrl: ctrl}
	mock.recorder = &MockNotificationServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockNotificationService) EXPECT() *MockNotificationServiceMockRecorder {
	return m.recorder
}

// SendEmailNotification mocks base method.
func (m *MockNotificationService) SendEmailNotification(recipient, subject, message string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendEmailNotification", recipient, subject, message)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendEmailNotification indicates an expected call of SendEmailNotification.
func (mr *MockNotificationServiceMockRecorder) SendEmailNotification(recipient, subject, message interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendEmailNotification", reflect.TypeOf((*MockNotificationService)(nil).SendEmailNotification), recipient, subject, message)
}

// SendEmailVerificationNotification mocks base method.
func (m *MockNotificationService) SendEmailVerificationNotification(recipient, token string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendEmailVerificationNotification", recipient, token)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendEmailVerificationNotification indicates an expected call of SendEmailVerificationNotification.
func (mr *MockNotificationServiceMockRecorder) SendEmailVerificationNotification(recipient, token interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendEmailVerificationNotification", reflect.TypeOf((*MockNotificationService)(nil).SendEmailVerificationNotification), recipient, token)
}

// MockEmailService is a mock of EmailService interface.
type MockEmailService struct {
	ctrl     *gomock.Controller
	recorder *MockEmailServiceMockRecorder
}

// MockEmailServiceMockRecorder is the mock recorder for MockEmailService.
type MockEmailServiceMockRecorder struct {
	mock *MockEmailService
}

// NewMockEmailService creates a new mock instance.
func NewMockEmailService(ctrl *gomock.Controller) *MockEmailService {
	mock := &MockEmailService{ctrl: ctrl}
	mock.recorder = &MockEmailServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockEmailService) EXPECT() *MockEmailServiceMockRecorder {
	return m.recorder
}

// SendEmail mocks base method.
func (m *MockEmailService) SendEmail(recipient, subject, body string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendEmail", recipient, subject, body)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendEmail indicates an expected call of SendEmail.
func (mr *MockEmailServiceMockRecorder) SendEmail(recipient, subject, body interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendEmail", reflect.TypeOf((*MockEmailService)(nil).SendEmail), recipient, subject, body)
}

// MockSmsService is a mock of SmsService interface.
type MockSmsService struct {
	ctrl     *gomock.Controller
	recorder *MockSmsServiceMockRecorder
}

// MockSmsServiceMockRecorder is the mock recorder for MockSmsService.
type MockSmsServiceMockRecorder struct {
	mock *MockSmsService
}

// NewMockSmsService creates a new mock instance.
func NewMockSmsService(ctrl *gomock.Controller) *MockSmsService {
	mock := &MockSmsService{ctrl: ctrl}
	mock.recorder = &MockSmsServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSmsService) EXPECT() *MockSmsServiceMockRecorder {
	return m.recorder
}

// SendSms mocks base method.
func (m *MockSmsService) SendSms(recipient, message string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendSms", recipient, message)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendSms indicates an expected call of SendSms.
func (mr *MockSmsServiceMockRecorder) SendSms(recipient, message interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendSms", reflect.TypeOf((*MockSmsService)(nil).SendSms), recipient, message)
}

// MockCacheService is a mock of CacheService interface.
type MockCacheService struct {
	ctrl     *gomock.Controller
	recorder *MockCacheServiceMockRecorder
}

// MockCacheServiceMockRecorder is the mock recorder for MockCacheService.
type MockCacheServiceMockRecorder struct {
	mock *MockCacheService
}

// NewMockCacheService creates a new mock instance.
func NewMockCacheService(ctrl *gomock.Controller) *MockCacheService {
	mock := &MockCacheService{ctrl: ctrl}
	mock.recorder = &MockCacheServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCacheService) EXPECT() *MockCacheServiceMockRecorder {
	return m.recorder
}

// LookupUrl mocks base method.
func (m *MockCacheService) LookupUrl(shortCode string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LookupUrl", shortCode)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LookupUrl indicates an expected call of LookupUrl.
func (mr *MockCacheServiceMockRecorder) LookupUrl(shortCode interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LookupUrl", reflect.TypeOf((*MockCacheService)(nil).LookupUrl), shortCode)
}

// SaveURL mocks base method.
func (m *MockCacheService) SaveURL(shortCode, originalUrl string, expiryTime time.Duration) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveURL", shortCode, originalUrl, expiryTime)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SaveURL indicates an expected call of SaveURL.
func (mr *MockCacheServiceMockRecorder) SaveURL(shortCode, originalUrl, expiryTime interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveURL", reflect.TypeOf((*MockCacheService)(nil).SaveURL), shortCode, originalUrl, expiryTime)
}

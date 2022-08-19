// Code generated by MockGen. DO NOT EDIT.
// Source: app/internal/core/contracts/services.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
	entities "github.com/sanctumlabs/curtz/app/internal/core/entities"
	identifier "github.com/sanctumlabs/curtz/app/pkg/identifier"
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
func (m *MockAuthService) GenerateRefreshToken(userId string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateRefreshToken", userId)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GenerateRefreshToken indicates an expected call of GenerateRefreshToken.
func (mr *MockAuthServiceMockRecorder) GenerateRefreshToken(userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateRefreshToken", reflect.TypeOf((*MockAuthService)(nil).GenerateRefreshToken), userId)
}

// GenerateToken mocks base method.
func (m *MockAuthService) GenerateToken(userId string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateToken", userId)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GenerateToken indicates an expected call of GenerateToken.
func (mr *MockAuthServiceMockRecorder) GenerateToken(userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateToken", reflect.TypeOf((*MockAuthService)(nil).GenerateToken), userId)
}

// MockUrlService is a mock of UrlService interface.
type MockUrlService struct {
	ctrl     *gomock.Controller
	recorder *MockUrlServiceMockRecorder
}

// MockUrlServiceMockRecorder is the mock recorder for MockUrlService.
type MockUrlServiceMockRecorder struct {
	mock *MockUrlService
}

// NewMockUrlService creates a new mock instance.
func NewMockUrlService(ctrl *gomock.Controller) *MockUrlService {
	mock := &MockUrlService{ctrl: ctrl}
	mock.recorder = &MockUrlServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUrlService) EXPECT() *MockUrlServiceMockRecorder {
	return m.recorder
}

// CreateUrl mocks base method.
func (m *MockUrlService) CreateUrl(userId, originalUrl, customAlias, expiresOn string, keywords []string) (entities.URL, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUrl", userId, originalUrl, customAlias, expiresOn, keywords)
	ret0, _ := ret[0].(entities.URL)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUrl indicates an expected call of CreateUrl.
func (mr *MockUrlServiceMockRecorder) CreateUrl(userId, originalUrl, customAlias, expiresOn, keywords interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUrl", reflect.TypeOf((*MockUrlService)(nil).CreateUrl), userId, originalUrl, customAlias, expiresOn, keywords)
}

// GetById mocks base method.
func (m *MockUrlService) GetById(id string) (entities.URL, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetById", id)
	ret0, _ := ret[0].(entities.URL)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetById indicates an expected call of GetById.
func (mr *MockUrlServiceMockRecorder) GetById(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetById", reflect.TypeOf((*MockUrlService)(nil).GetById), id)
}

// GetByKeyword mocks base method.
func (m *MockUrlService) GetByKeyword(keyword string) ([]entities.URL, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByKeyword", keyword)
	ret0, _ := ret[0].([]entities.URL)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByKeyword indicates an expected call of GetByKeyword.
func (mr *MockUrlServiceMockRecorder) GetByKeyword(keyword interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByKeyword", reflect.TypeOf((*MockUrlService)(nil).GetByKeyword), keyword)
}

// GetByKeywords mocks base method.
func (m *MockUrlService) GetByKeywords(keywords []string) ([]entities.URL, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByKeywords", keywords)
	ret0, _ := ret[0].([]entities.URL)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByKeywords indicates an expected call of GetByKeywords.
func (mr *MockUrlServiceMockRecorder) GetByKeywords(keywords interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByKeywords", reflect.TypeOf((*MockUrlService)(nil).GetByKeywords), keywords)
}

// GetByOriginalUrl mocks base method.
func (m *MockUrlService) GetByOriginalUrl(originalUrl string) (entities.URL, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByOriginalUrl", originalUrl)
	ret0, _ := ret[0].(entities.URL)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByOriginalUrl indicates an expected call of GetByOriginalUrl.
func (mr *MockUrlServiceMockRecorder) GetByOriginalUrl(originalUrl interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByOriginalUrl", reflect.TypeOf((*MockUrlService)(nil).GetByOriginalUrl), originalUrl)
}

// GetByShortCode mocks base method.
func (m *MockUrlService) GetByShortCode(shortCode string) (entities.URL, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByShortCode", shortCode)
	ret0, _ := ret[0].(entities.URL)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByShortCode indicates an expected call of GetByShortCode.
func (mr *MockUrlServiceMockRecorder) GetByShortCode(shortCode interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByShortCode", reflect.TypeOf((*MockUrlService)(nil).GetByShortCode), shortCode)
}

// GetByUserId mocks base method.
func (m *MockUrlService) GetByUserId(userId string) ([]entities.URL, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByUserId", userId)
	ret0, _ := ret[0].([]entities.URL)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByUserId indicates an expected call of GetByUserId.
func (mr *MockUrlServiceMockRecorder) GetByUserId(userId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByUserId", reflect.TypeOf((*MockUrlService)(nil).GetByUserId), userId)
}

// LookupUrl mocks base method.
func (m *MockUrlService) LookupUrl(shortCode string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LookupUrl", shortCode)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LookupUrl indicates an expected call of LookupUrl.
func (mr *MockUrlServiceMockRecorder) LookupUrl(shortCode interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LookupUrl", reflect.TypeOf((*MockUrlService)(nil).LookupUrl), shortCode)
}

// Remove mocks base method.
func (m *MockUrlService) Remove(id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Remove", id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Remove indicates an expected call of Remove.
func (mr *MockUrlServiceMockRecorder) Remove(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Remove", reflect.TypeOf((*MockUrlService)(nil).Remove), id)
}

// MockUserService is a mock of UserService interface.
type MockUserService struct {
	ctrl     *gomock.Controller
	recorder *MockUserServiceMockRecorder
}

// MockUserServiceMockRecorder is the mock recorder for MockUserService.
type MockUserServiceMockRecorder struct {
	mock *MockUserService
}

// NewMockUserService creates a new mock instance.
func NewMockUserService(ctrl *gomock.Controller) *MockUserService {
	mock := &MockUserService{ctrl: ctrl}
	mock.recorder = &MockUserServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserService) EXPECT() *MockUserServiceMockRecorder {
	return m.recorder
}

// CreateUser mocks base method.
func (m *MockUserService) CreateUser(email, password string) (entities.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", email, password)
	ret0, _ := ret[0].(entities.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockUserServiceMockRecorder) CreateUser(email, password interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockUserService)(nil).CreateUser), email, password)
}

// GetByVerificationToken mocks base method.
func (m *MockUserService) GetByVerificationToken(verificationToken string) (entities.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByVerificationToken", verificationToken)
	ret0, _ := ret[0].(entities.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByVerificationToken indicates an expected call of GetByVerificationToken.
func (mr *MockUserServiceMockRecorder) GetByVerificationToken(verificationToken interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByVerificationToken", reflect.TypeOf((*MockUserService)(nil).GetByVerificationToken), verificationToken)
}

// GetUserByEmail mocks base method.
func (m *MockUserService) GetUserByEmail(email string) (entities.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByEmail", email)
	ret0, _ := ret[0].(entities.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByEmail indicates an expected call of GetUserByEmail.
func (mr *MockUserServiceMockRecorder) GetUserByEmail(email interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByEmail", reflect.TypeOf((*MockUserService)(nil).GetUserByEmail), email)
}

// GetUserByID mocks base method.
func (m *MockUserService) GetUserByID(id string) (entities.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByID", id)
	ret0, _ := ret[0].(entities.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByID indicates an expected call of GetUserByID.
func (mr *MockUserServiceMockRecorder) GetUserByID(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByID", reflect.TypeOf((*MockUserService)(nil).GetUserByID), id)
}

// RemoveUser mocks base method.
func (m *MockUserService) RemoveUser(id string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveUser", id)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveUser indicates an expected call of RemoveUser.
func (mr *MockUserServiceMockRecorder) RemoveUser(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveUser", reflect.TypeOf((*MockUserService)(nil).RemoveUser), id)
}

// SetVerified mocks base method.
func (m *MockUserService) SetVerified(id identifier.ID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetVerified", id)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetVerified indicates an expected call of SetVerified.
func (mr *MockUserServiceMockRecorder) SetVerified(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetVerified", reflect.TypeOf((*MockUserService)(nil).SetVerified), id)
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

// SaveUrl mocks base method.
func (m *MockCacheService) SaveUrl(shortCode, originalUrl string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveUrl", shortCode, originalUrl)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SaveUrl indicates an expected call of SaveUrl.
func (mr *MockCacheServiceMockRecorder) SaveUrl(shortCode, originalUrl interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveUrl", reflect.TypeOf((*MockCacheService)(nil).SaveUrl), shortCode, originalUrl)
}
package contracts

import (
	"time"

	"github.com/sanctumlabs/curtz/app/internal/core/entities"
	"github.com/sanctumlabs/curtz/app/pkg/identifier"
)

type AuthService interface {
	// Authenticate a user given the token. Returns user id if authenticated, error otherwise.
	Authenticate(token string) (string, time.Time, error)
	// GenerateToken generates a user token provided the user id
	GenerateToken(userID string) (string, error)
	// GenerateRefreshToken generates a new refresh token
	GenerateRefreshToken(userID string) (string, error)
}

type UrlService interface {
	CreateUrl(userID string, originalUrl string, customAlias string, expiresOn time.Time, keywords []string) (entities.URL, error)
	GetByShortCode(shortCode string) (entities.URL, error)
	GetByUserId(userID string) ([]entities.URL, error)
	GetByKeyword(keyword string) ([]entities.URL, error)
	GetByKeywords(keywords []string) ([]entities.URL, error)
	GetByOriginalUrl(originalUrl string) (entities.URL, error)
	GetById(id string) (entities.URL, error)
	Remove(id string) error
	LookupUrl(shortCode string) (string, error)
}

type UserService interface {
	CreateUser(email, password string) (entities.User, error)
	GetUserByEmail(email string) (entities.User, error)
	GetUserByID(id string) (entities.User, error)
	GetByVerificationToken(verificationToken string) (entities.User, error)
	SetVerified(id identifier.ID) error
	RemoveUser(id string) error
}

type NotificationService interface {
	SendEmailNotification(recipient, subject, message string) error
	SendEmailVerificationNotification(recipient, token string) error
}

type EmailService interface {
	SendEmail(recipient, subject, body string) error
}

type SmsService interface {
	SendSms(recipient, message string) error
}

// CacheService interface to be used by services that implement cache like functionality
type CacheService interface {
	LookupUrl(shortCode string) (string, error)
	SaveUrl(shortCode, originalUrl string) (string, error)
}

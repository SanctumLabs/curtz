package contracts

import (
	"time"

	"github.com/sanctumlabs/curtz/app/internal/core/entities"
)

type AuthService interface {
	// Authenticate a user given the token. Returns user id if authenticated, error otherwise.
	Authenticate(token string) (string, time.Time, error)
	// GenerateToken generates a user token provided the user id
	GenerateToken(userId string) (string, error)
}

type UrlService interface {
	CreateUrl(userId, originalUrl, customAlias, expiresOn string, keywords []string) (entities.URL, error)
	GetByShortCode(shortCode string) (entities.URL, error)
	GetByUserId(userId string) ([]entities.URL, error)
	GetByKeyword(keyword string) ([]entities.URL, error)
	GetByKeywords(keywords []string) ([]entities.URL, error)
	GetByOriginalUrl(originalUrl string) ([]entities.URL, error)
	GetById(id string) (entities.URL, error)
	Remove(id string) error
}

type UserService interface {
	CreateUser(email, password string) (entities.User, error)
	GetUserByEmail(email string) (entities.User, error)
	GetUserByID(id string) (entities.User, error)
	GetUserByToken(token string) (entities.User, error)
	GetUserByUsername(username string) (entities.User, error)
	RemoveUser(id string) error
}

type NotificationService interface {
	SendNotification(recipient, message, notifyType string) error
}

type EmailService interface {
	SendEmail(recipient, subject, body string) error
}

type SmsService interface {
	SendSms(recipient, message string) error
}

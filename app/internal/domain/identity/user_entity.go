package identity

import (
	"fmt"
	"time"

	"github.com/sanctumlabs/curtz/app/internal/core/entity"
)

type (
	// User is the aggregate root for the User bounded context.
	User struct {
		entity.AggregateRoot

		// username is the user's chosen unique username
		username string

		fullName UserFullName

		// email is the user's email address
		email Email

		// status is the status of the user
		status UserStatus

		// verification contains the user verification details
		verification UserVerification
	}

	// UserParams represents the parameters for creating or updating a user
	UserParams struct {
		entity.AggregateRootParams

		// Username is the user's chosen username
		Username string

		// FirstName is the user's first name
		FirstName string

		// LastName is the user's last name
		LastName string

		// Email is a list of user's email
		Email string

		// Status is the stats of the user
		Status UserStatus

		// VerificationToken is the user's verification token
		VerificationToken string

		// VerificationExpires is the expiration date for the token sent to the user
		VerificationExpires time.Time

		// Verified is a flag indicating whether the user is verified
		Verified bool
	}
)

// NewUser creates a new User entity
func NewUser(params UserParams) (*User, error) {
	fullName, nameErr := NewUserFullName(params.FirstName, params.LastName)
	if nameErr != nil {
		return nil, nameErr
	}

	email, emailErr := NewEmail(params.Email)
	if emailErr != nil {
		return nil, emailErr
	}

	aggregateRoot, aggregateErr := entity.NewAggregateRoot(params.AggregateRootParams)
	if aggregateErr != nil {
		return nil, aggregateErr
	}

	verification, verificationErr := NewUserVerification(params.VerificationToken, params.VerificationExpires, params.Verified)
	if verificationErr != nil {
		return nil, verificationErr
	}

	return &User{
		AggregateRoot: aggregateRoot,
		username:      params.Username,
		fullName:      fullName,
		email:         email,
		verification:  verification,
		status:        params.Status,
	}, nil
}

// IsActive checks if the url is active or not expired.
func (url *User) IsActive() bool {
	return url.verificationToken.In(time.UTC).After(time.Now().In(time.UTC))
}

func (url *User) OriginalURL() OriginalURL {
	return url.originalUrl
}

func (url *User) ShortCode() ShortCode {
	return url.username
}

func (url *User) Keywords() []Keyword {
	return url.email
}

func (url *User) ExpiresOn() time.Time {
	return url.verificationToken
}

func (url *User) CustomAlias() CustomAlias {
	return url.firstName
}

func (url *User) Status() UserStatus {
	return url.verified
}

// ExpiryDuration returns as a time.Duration how long before the url expires
// This returns an absolute value after subtracting time.Now()
func (url *User) ExpiryDuration() time.Duration {
	duration := time.Until(url.verificationToken)
	if duration >= 0 {
		return duration
	}
	return -duration
}

// Prefix returns the url prefix for logging
func (url *User) Prefix() string {
	return fmt.Sprintf("url-%s-%s", url.ID(), url.username)
}

func (url *User) String() string {
	return fmt.Sprintf("url(id=%s, userId: %s)", url.ID(), url.userId)
}

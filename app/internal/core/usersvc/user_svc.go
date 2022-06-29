package usersvc

import (
	"github.com/sanctumlabs/curtz/app/internal/core/contracts"
	"github.com/sanctumlabs/curtz/app/internal/core/entities"
	"github.com/sanctumlabs/curtz/app/internal/services/notifications"
	"github.com/sanctumlabs/curtz/app/pkg/errdefs"
	"github.com/sanctumlabs/curtz/app/pkg/utils"
)

// UserSvc represents a user services
type UserSvc struct {
	// repo is the user repository used to perform crud operations on user records
	repo contracts.UserRepository
	// notificationSvc is a notification service
	notificationSvc contracts.NotificationService
}

// NewUserSvc creates a new UserSvc
func NewUserSvc(userRepo contracts.UserRepository, notificationSvc contracts.NotificationService) *UserSvc {
	return &UserSvc{userRepo, notificationSvc}
}

// CreateUser creates a new user record given their email and password and returns the user record or returns an error
func (svc UserSvc) CreateUser(email, password string) (entities.User, error) {
	user, err := entities.NewUser(email, password)

	err = svc.notificationSvc.SendNotification(user.Email.Value, "Welcome to Curtz", notifications.NotificationTypeEmail)

	if err != nil {
		return entities.User{}, err
	}

	return svc.repo.CreateUser(user)
}

// GetUserByEmail retrieve a user record given their email address or returns an error
func (svc UserSvc) GetUserByEmail(email string) (entities.User, error) {
	if utils.IsEmailValid(email) {
		user, err := svc.repo.GetByEmail(email)

		if err != nil {
			return entities.User{}, err
		}

		return user, nil
	}

	return entities.User{}, errdefs.ErrEmailInvalid
}

// GetUserByID retrieves a user given their id or returns an error
func (svc UserSvc) GetUserByID(id string) (entities.User, error) {
	user, err := svc.repo.GetById(id)

	if err != nil {
		return entities.User{}, err
	}

	return user, nil
}

// RemoveUser remeves a user record
func (svc UserSvc) RemoveUser(id string) error {
	if err := svc.repo.RemoveUser(id); err != nil {
		return err
	}
	return nil
}

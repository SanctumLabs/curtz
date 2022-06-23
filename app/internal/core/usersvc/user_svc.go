package usersvc

import (
	"github.com/sanctumlabs/curtz/app/internal/core/contracts"
	"github.com/sanctumlabs/curtz/app/internal/core/entities"
	"github.com/sanctumlabs/curtz/app/pkg/errdefs"
	"github.com/sanctumlabs/curtz/app/pkg/utils"
)

type UserSvc struct {
	repo contracts.UserRepository
}

func NewUserSvc(userRepo contracts.UserRepository) *UserSvc {
	return &UserSvc{userRepo}
}

func (svc UserSvc) CreateUser(email, password string) (entities.User, error) {
	user, err := entities.NewUser(email, password)

	if err != nil {
		return entities.User{}, err
	}

	return svc.repo.CreateUser(user)
}

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

func (svc UserSvc) GetUserByID(id string) (entities.User, error) {
	panic("implement me")
}

func (svc UserSvc) RemoveUser(id string) error {
	panic("implement me")
}

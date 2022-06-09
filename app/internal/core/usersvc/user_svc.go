package usersvc

import (
	"github.com/sanctumlabs/curtz/app/internal/core/contracts"
	"github.com/sanctumlabs/curtz/app/internal/core/entities"
)

type UserSvc struct {
	repo contracts.UserRepository
}

func NewUserSvc(userRepo contracts.UserRepository) *UserSvc {
	return &UserSvc{userRepo}
}

func (svc UserSvc) CreateUser(email, password string) (entities.User, error) {
	return svc.repo.CreateUser(email, password)
}

func (svc UserSvc) GetUserByEmail(email string) (entities.User, error) {
	panic("implement me")
}

func (svc UserSvc) GetUserByID(id string) (entities.User, error) {
	panic("implement me")
}

func (svc UserSvc) GetUserByToken(token string) (entities.User, error) {
	panic("implement me")
}

func (svc UserSvc) GetUserByUsername(username string) (entities.User, error) {
	panic("implement me")
}

func (svc UserSvc) RemoveUser(id string) error {
	panic("implement me")
}

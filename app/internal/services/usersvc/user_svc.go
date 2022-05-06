package usersvc

import (
	"github.com/sanctumlabs/curtz/app/internal/core/domain"
	"github.com/sanctumlabs/curtz/app/internal/core/domain/entity"
)

type service struct {
	usecase *domain.UserInteractor
}

func NewService(userInteractor *domain.UserInteractor) *service {
	return &service{usecase: userInteractor}
}

func (svc *service) CreateUser(email, password string) (entity.User, error) {
	return svc.usecase.CreateUser(email, password)
}

func (svc *service) GetUserByEmail(email string) (entity.User, error) {
	panic("implement me")
}

func (svc *service) GetUserByID(id string) (entity.User, error) {
	panic("implement me")
}

func (svc *service) GetUserByToken(token string) (entity.User, error) {
	panic("implement me")
}

func (svc *service) GetUserByUsername(username string) (entity.User, error) {
	panic("implement me")
}

func (svc *service) RemoveUser(id string) error {
	panic("implement me")
}

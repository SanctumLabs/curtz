package domain

import (
	"github.com/sanctumlabs/curtz/app/internal/core/contracts"
	"github.com/sanctumlabs/curtz/app/internal/core/domain/entity"
)

type UserInteractor struct {
	repo contracts.UserRepository
}

func NewUserInteractor(userRepo contracts.UserRepository) *UserInteractor {
	return &UserInteractor{userRepo}
}

func (useCase UserInteractor) CreateUser(email, password string) (entity.User, error) {
	return useCase.repo.CreateUser(email, password)
}

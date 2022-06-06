package userepo

import (
	entity "github.com/sanctumlabs/curtz/app/internal/core/entities"
	"gorm.io/gorm"
)

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

func (u *UserRepo) CreateUser(email, password string) (entity.User, error) {
	panic("implement me")
}

func (u *UserRepo) GetByEmail(email string) (entity.User, error) {
	panic("implement me")
}

func (u *UserRepo) GetById(id string) (entity.User, error) {
	panic("implement me")
}

func (u *UserRepo) GetByUsername(username string) (entity.User, error) {
	panic("implement me")
}

func (u *UserRepo) RemoveUser(id string) error {
	panic("implement me")
}

package userepo

import (
	"github.com/google/uuid"
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

func (u UserRepo) CreateUser(email, password string) (uuid.UUID, error) {

	panic("implement me")
}

func (u UserRepo) GetByEmail(email string) (uuid.UUID, error) {

	panic("implement me")
}

func (u UserRepo) GetById(id uuid.UUID) (uuid.UUID, error) {

	panic("implement me")
}

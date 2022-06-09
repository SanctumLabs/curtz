package userepo

import (
	"context"

	entity "github.com/sanctumlabs/curtz/app/internal/core/entities"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepo struct {
	dbClient *mongo.Collection
	context  context.Context
}

func NewUserRepo(dbClient *mongo.Collection, ctx context.Context) *UserRepo {
	return &UserRepo{
		dbClient: dbClient,
		context:  ctx,
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

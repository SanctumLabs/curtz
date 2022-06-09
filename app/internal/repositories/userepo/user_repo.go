package userepo

import (
	"context"

	"github.com/sanctumlabs/curtz/app/internal/core/entities"
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

func (u *UserRepo) CreateUser(user entities.User) (entities.User, error) {
	panic("implement me")
}

func (u *UserRepo) GetByEmail(email string) (entities.User, error) {
	panic("implement me")
}

func (u *UserRepo) GetById(id string) (entities.User, error) {
	panic("implement me")
}

func (u *UserRepo) GetByUsername(username string) (entities.User, error) {
	panic("implement me")
}

func (u *UserRepo) RemoveUser(id string) error {
	panic("implement me")
}

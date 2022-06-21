package userepo

import (
	"context"

	"github.com/sanctumlabs/curtz/app/internal/core/entities"
	"github.com/sanctumlabs/curtz/app/internal/repositories/models"
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

	userModel := models.User{
		BaseModel: models.BaseModel{
			Id:        user.ID.String(),
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
		Email:               user.Email.Value,
		Password:            user.Password.Value,
		VerificationToken:   user.VerificationToken,
		VerificationExpires: user.VerificationExpires,
	}

	_, err := u.dbClient.InsertOne(u.context, userModel)

	if err != nil {
		return user, err
	}

	return user, nil
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

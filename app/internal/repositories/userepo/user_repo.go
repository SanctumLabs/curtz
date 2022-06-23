package userepo

import (
	"context"

	"github.com/sanctumlabs/curtz/app/internal/core/entities"
	"github.com/sanctumlabs/curtz/app/internal/repositories/models"
	"github.com/sanctumlabs/curtz/app/pkg/errdefs"
	"github.com/sanctumlabs/curtz/app/pkg/identifier"
	"go.mongodb.org/mongo-driver/bson"
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
	if _, err := u.GetByEmail(user.Email.Value); err == nil {
		return entities.User{}, errdefs.ErrUserExists
	}

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
	filter := bson.D{{Key: "email", Value: email}}

	var result bson.D
	if err := u.dbClient.FindOne(u.context, filter).Decode(&result); err != nil {
		return entities.User{}, err
	}

	document, err := bson.Marshal(result)
	if err != nil {
		return entities.User{}, err
	}

	var user models.User
	err = bson.Unmarshal(document, &user)
	if err != nil {
		return entities.User{}, err
	}

	return entities.User{
		ID:       identifier.New().FromString(user.BaseModel.Id),
		Email:    entities.Email{Value: user.Email},
		Password: entities.Password{Value: user.Password},
		BaseEntity: entities.BaseEntity{
			CreatedAt: user.BaseModel.CreatedAt,
			UpdatedAt: user.BaseModel.UpdatedAt,
			DeletedAt: user.BaseModel.DeletedAt,
		},
		Verified: user.Verified,
	}, nil
}

func (u *UserRepo) GetById(id string) (entities.User, error) {
	panic("implement me")
}

func (u *UserRepo) RemoveUser(id string) error {
	panic("implement me")
}

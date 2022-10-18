package userepo

import (
	"github.com/sanctumlabs/curtz/app/internal/core/entities"
	"github.com/sanctumlabs/curtz/app/internal/repositories/models"
	"github.com/sanctumlabs/curtz/app/pkg/identifier"
)

func mapEntityToModel(user entities.User) models.User {
	return models.User{
		BaseModel: models.BaseModel{
			Id:        user.GetId(),
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
		Email:               user.GetEmail(),
		Password:            user.GetPassword(),
		VerificationToken:   user.VerificationToken.String(),
		VerificationExpires: user.VerificationExpires,
	}
}

func mapModelToEntity(user models.User) (entities.User, error) {
	email, emailErr := entities.NewEmail(user.Email)
	if emailErr != nil {
		return entities.User{}, emailErr
	}

	password, passwordErr := entities.NewPassword(user.Password)
	if passwordErr != nil {
		return entities.User{}, passwordErr
	}

	id, idErr := identifier.New().FromString(user.BaseModel.Id)
	if idErr != nil {
		return entities.User{}, idErr
	}

	u, err := entities.NewUser(email.GetValue(), password.GetValue())
	if err != nil {
		return entities.User{}, err
	}

	u.SetId(id.String())

	u.CreatedAt = user.BaseModel.CreatedAt
	u.UpdatedAt = user.BaseModel.UpdatedAt
	u.Verified = user.Verified

	return *u, nil
}

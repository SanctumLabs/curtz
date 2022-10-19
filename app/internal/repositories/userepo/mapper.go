package userepo

import (
	"github.com/sanctumlabs/curtz/app/internal/core/entities"
	"github.com/sanctumlabs/curtz/app/internal/repositories/models"
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
	u := entities.User{
		BaseEntity: entities.BaseEntity{
			CreatedAt: user.BaseModel.CreatedAt,
			UpdatedAt: user.BaseModel.UpdatedAt,
		},
		Verified: user.Verified,
	}

	idErr := u.SetId(user.BaseModel.Id)
	if idErr != nil {
		return entities.User{}, idErr
	}

	emailErr := u.SetEmail(user.Email)
	if emailErr != nil {
		return entities.User{}, emailErr
	}

	passwordErr := u.SetPassword(user.Password)
	if passwordErr != nil {
		return entities.User{}, passwordErr
	}

	return u, nil
}

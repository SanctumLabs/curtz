package userepo

import (
	"github.com/sanctumlabs/curtz/app/internal/core/entities"
	"github.com/sanctumlabs/curtz/app/internal/repositories/models"
	"github.com/sanctumlabs/curtz/app/pkg/identifier"
)

func mapEntityToModel(user entities.User) models.User {
	return models.User{
		BaseModel: models.BaseModel{
			Id:        user.ID.String(),
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
		Email:               user.Email.Value,
		Password:            user.Password.Value,
		VerificationToken:   user.VerificationToken.String(),
		VerificationExpires: user.VerificationExpires,
	}
}

func mapModelToEntity(user models.User) entities.User {
	return entities.User{
		ID:       identifier.New().FromString(user.BaseModel.Id),
		Email:    entities.Email{Value: user.Email},
		Password: entities.Password{Value: user.Password},
		BaseEntity: entities.BaseEntity{
			CreatedAt: user.BaseModel.CreatedAt,
			UpdatedAt: user.BaseModel.UpdatedAt,
		},
		Verified: user.Verified,
	}
}

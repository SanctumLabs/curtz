//go:build legacy

package contracts

import (
	"github.com/sanctumlabs/fupi/app/internal/core/entities"
	"github.com/sanctumlabs/fupi/app/pkg/identifier"
)

type UserRepository interface {
	CreateUser(entities.User) (entities.User, error)
	GetByEmail(email string) (entities.User, error)
	GetById(id string) (entities.User, error)
	RemoveUser(id string) error
	GetByVerificationToken(verificationToken string) (entities.User, error)
	SetVerified(id identifier.ID) error
}

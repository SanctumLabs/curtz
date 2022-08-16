package data

import "github.com/sanctumlabs/curtz/app/internal/core/entities"

func MockUser(email, password string) (entities.User, error) {
	user, err := entities.NewUser(email, password)
	return user, err
}

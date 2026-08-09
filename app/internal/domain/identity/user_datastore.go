package identity

import (
	"context"

	"github.com/sanctumlabs/curtz/app/internal/core/ports/repository"
)

type (

	// UserDatastore defines the interface for interacting with User entities in an underlying storage implementation
	UserDatastore interface {
		UserReadDatastore
		UserWriteDatastore
	}

	// UserWriteDatastore defines the interface for writing User entities to an underlying storage implementation
	UserWriteDatastore interface {
		repository.WriteRepositoryPort[User]

		// Create creates a new User entity based on the provided request and returns the created User entity
		Create(ctx context.Context, request CreateUserRequest) (User, error)

		// UpdateVerification updates the verification details of a User entity based on the provided request and returns the updated User entity
		UpdateVerification(ctx context.Context, request UpdateUserVerificationRequest) (User, error)
	}

	// UserReadDatastore defines the interface for reading User entities from an underlying storage implementation
	UserReadDatastore interface {
		repository.ReadRepositoryPort[User]

		// FetchByUsername retrieves a User entity by its username
		FetchByUsername(ctx context.Context, username string) (User, error)

		// FetchByEmail retrieves a User entity by its email
		FetchByEmail(ctx context.Context, email string) (User, error)

		// FetchByStatus retrieves a list of User entities by their status
		FetchByStatus(ctx context.Context, status UserStatus) (repository.FetchRecordsResponse[User], error)
	}
)

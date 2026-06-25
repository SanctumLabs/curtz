package identity

import (
	"context"

	"github.com/sanctumlabs/curtz/app/internal/core/ports/repository"
)

type (

	// UserRepository defines the interface for interacting with User entities in an underlying storage implementation
	UserRepository interface {
		UserReadRepository
		UserWriteRepository
	}

	// UserWriteRepository defines the interface for writing User entities to an underlying storage implementation
	UserWriteRepository interface {
		repository.WriteRepositoryPort[User]
	}

	// UserReadRepository defines the interface for reading User entities from an underlying storage implementation
	UserReadRepository interface {
		repository.ReadRepositoryPort[User]

		// FetchByShortCode retrieves a URL entity by its short code
		FetchByShortCode(ctx context.Context, shortCode string) (User, error)

		// FetchByCustomAlias retrieves a URL entity by its custom alias
		FetchByCustomAlias(ctx context.Context, customAlias string) (User, error)

		// FetchByUserId retrieves a list of URL entities by the user id of their owner
		FetchByUserId(ctx context.Context, userId string) (repository.FetchRecordsResponse[User], error)

		// FetchByStatus retrieves a list of URL entities by their status
		FetchByStatus(ctx context.Context, status UserStatus) (repository.FetchRecordsResponse[User], error)

		// FetchByOriginalUrl retrieves a URL entity by its original url
		FetchByOriginalUrl(ctx context.Context, originalUrl string) (User, error)
	}
)

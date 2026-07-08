package identityrepo

import (
	"context"

	"github.com/sanctumlabs/curtz/app/internal/core/ports/repository"
	"github.com/sanctumlabs/curtz/app/internal/domain/identity"
	"github.com/sanctumlabs/curtz/app/internal/pkg/common"
	"github.com/sanctumlabs/curtz/app/pkg/infra/database"
)

func NewUserReadRepoAdapter(dbClient database.PostgresDatabaseClient) identity.UserReadRepository {
	repo := &userReadRepositoryAdapter{
		dbClient:  dbClient,
		logPrefix: "UserReadRepoAdapter"}

	return repo
}

func (repo *userReadRepositoryAdapter) FetchById(ctx context.Context, id string) (identity.User, error) {
	panic("not implemented")
}

func (repo *userReadRepositoryAdapter) FetchAll(ctx context.Context, params common.RequestParams) (repository.FetchRecordsResponse[identity.User], error) {
	panic("not implemented")
}

func (repo *userReadRepositoryAdapter) FetchByUsername(ctx context.Context, username string) (identity.User, error) {
	panic("not implemented")
}

func (repo *userReadRepositoryAdapter) FetchByEmail(ctx context.Context, email string) (identity.User, error) {
	panic("not implemented")
}

func (repo *userReadRepositoryAdapter) FetchByStatus(ctx context.Context, status identity.UserStatus) (repository.FetchRecordsResponse[identity.User], error) {
	panic("not implemented")
}

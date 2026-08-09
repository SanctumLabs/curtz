package identitydatastore

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	postgresrepo "github.com/sanctumlabs/curtz/app/internal/adapters/postgres"
	postgresql "github.com/sanctumlabs/curtz/app/internal/adapters/postgres/sql"
	"github.com/sanctumlabs/curtz/app/internal/domain/identity"
	"github.com/sanctumlabs/curtz/app/pkg/errdefs"
	"github.com/sanctumlabs/curtz/app/pkg/infra/database"
	"github.com/sanctumlabs/curtz/app/pkg/infra/database/postgres"
	"github.com/sanctumlabs/curtz/app/pkg/utils"
	recoveryutils "github.com/sanctumlabs/curtz/app/pkg/utils/recover"
)

func NewUserWriteDatastoreAdapter(dbClient database.PostgresDatabaseClient, config database.Config) identity.UserWriteDatastore {
	repo := &userWriteDatastoreAdapter{
		dbClient:  dbClient,
		config:    config,
		logPrefix: "UserWriteRepoAdapter",
	}

	// Wire up the real transaction executor. This delegates to postgres.WithTransaction,
	// which handles the pgxpool.Pool lifecycle. Tests override this field directly.
	repo.withTx = func(ctx context.Context, fn func(q postgresrepo.UserWriteQuerier) (identity.User, error)) (identity.User, error) {
		return postgres.WithTransaction(ctx, dbClient, func(qtx *postgresql.Queries) (identity.User, error) {
			// *postgresql.Queries satisfies userWriteQuerier, so we can pass it straight through.
			return fn(qtx)
		})
	}

	return repo
}

func (writeDatastore *userWriteDatastoreAdapter) Save(ctx context.Context, userEntity identity.User) (identity.User, error) {
	handlerLogPrefix := fmt.Sprintf("%s<Save>", writeDatastore.logPrefix)
	slog.InfoContext(ctx, fmt.Sprintf("%s Saving User", handlerLogPrefix), "user", userEntity)

	operationCtx, operationCancel := context.WithTimeout(ctx, writeDatastore.config.OperationTimeout)
	defer operationCancel()

	return recoveryutils.ExecuteWithRetry(
		operationCtx,
		func(retryCtx context.Context) (identity.User, error) {
			// Use writeDatastore.withTx instead of postgres.WithTransaction directly.
			// This is the only change to the business logic — everything inside
			// the closure is identical to the original implementation.
			return writeDatastore.withTx(retryCtx, func(qtx postgresrepo.UserWriteQuerier) (identity.User, error) {
				// Check context before proceeding
				select {
				case <-retryCtx.Done():
					slog.ErrorContext(retryCtx, "Operation cancelled before validation with error", "error", retryCtx.Err())
					return identity.User{}, fmt.Errorf("operation cancelled before validation: %w", retryCtx.Err())
				default:
				}

				// query the status ID
				status, statusErr := qtx.QueryUserStatusByName(retryCtx, string(userEntity.Status()))
				if statusErr != nil {
					slog.ErrorContext(
						retryCtx,
						fmt.Sprintf("%s Failed to retrieve user status", handlerLogPrefix),
						"user_status", userEntity.Status(),
						"error", statusErr,
					)
					if errors.Is(statusErr, pgx.ErrNoRows) {
						return identity.User{}, errdefs.NotFound(statusErr)
					}

					return identity.User{}, fmt.Errorf("failed to query user status: %w", statusErr)
				}

				metadata, metadataErr := userEntity.MetadataToBytes()
				if metadataErr != nil {
					slog.WarnContext(ctx, fmt.Sprintf("%s Failed to convert user metadata to bytes", handlerLogPrefix),
						"user", userEntity,
						"error", metadataErr)
				}

				email := userEntity.Email()
				createdUser, createdUserErr := qtx.QueryCreateUser(
					retryCtx,
					postgresql.QueryCreateUserParams{
						Username:     userEntity.Username(),
						FirstName:    pgtype.Text{String: userEntity.FirstName(), Valid: true},
						LastName:     pgtype.Text{String: userEntity.LastName(), Valid: true},
						Email:        email.Value(),
						PasswordHash: userEntity.PasswordHash(),
						StatusID:     status.ID,
						Metadata:     metadata,
					},
				)
				if createdUserErr != nil {
					slog.ErrorContext(
						retryCtx,
						fmt.Sprintf("%s Failed to create user", handlerLogPrefix),
						"username", userEntity.Username(),
						"error", createdUserErr,
					)
					return identity.User{}, fmt.Errorf("failed to create user: %w", createdUserErr)
				}

				// Map the created URL model back to an entity to return
				mappedUser, mapErr := MapUserModelToEntity(UserMapperParams{
					UserModel: createdUser,
					Status:    status.Name,
				})
				if mapErr != nil {
					slog.ErrorContext(
						retryCtx,
						fmt.Sprintf("%s Failed to map created user model to entity", handlerLogPrefix),
						"user_model", createdUser,
						"error", mapErr,
					)
					return identity.User{}, fmt.Errorf("failed to map created user model to entity: %w", mapErr)
				}
				slog.InfoContext(retryCtx, "created model", "user", mappedUser)

				return mappedUser, nil
			})
		},
		writeDatastore.config.RetryConfig,
		fmt.Sprintf("%s.Create", writeDatastore.logPrefix),
	)
}

func (writeDatastore *userWriteDatastoreAdapter) Create(ctx context.Context, request identity.CreateUserRequest) (identity.User, error) {
	handlerLogPrefix := fmt.Sprintf("%s<Create>", writeDatastore.logPrefix)
	slog.InfoContext(ctx, fmt.Sprintf("%s Creating User", handlerLogPrefix), "user", request)

	operationCtx, operationCancel := context.WithTimeout(ctx, writeDatastore.config.OperationTimeout)
	defer operationCancel()

	return recoveryutils.ExecuteWithRetry(
		operationCtx,
		func(retryCtx context.Context) (identity.User, error) {
			// Use writeDatastore.withTx instead of postgres.WithTransaction directly.
			// This is the only change to the business logic — everything inside
			// the closure is identical to the original implementation.
			return writeDatastore.withTx(retryCtx, func(qtx postgresrepo.UserWriteQuerier) (identity.User, error) {
				// Check context before proceeding
				select {
				case <-retryCtx.Done():
					slog.ErrorContext(retryCtx, "Operation cancelled before validation with error", "error", retryCtx.Err())
					return identity.User{}, fmt.Errorf("operation cancelled before validation: %w", retryCtx.Err())
				default:
				}

				// query the status ID
				status, statusErr := qtx.QueryUserStatusByName(retryCtx, string(identity.UserStatusInactive))
				if statusErr != nil {
					slog.ErrorContext(
						retryCtx,
						fmt.Sprintf("%s Failed to retrieve user status", handlerLogPrefix),
						"user_status", identity.UserStatusInactive,
						"error", statusErr,
					)
					if errors.Is(statusErr, pgx.ErrNoRows) {
						return identity.User{}, errdefs.NotFound(statusErr)
					}

					return identity.User{}, fmt.Errorf("failed to query user status: %w", statusErr)
				}

				metadata, metadataErr := utils.MapToBytes(request.Metadata)
				if metadataErr != nil {
					slog.WarnContext(ctx, fmt.Sprintf("%s Failed to convert user metadata to bytes", handlerLogPrefix),
						"user", request,
						"error", metadataErr)
				}

				createdUser, createdUserErr := qtx.QueryCreateUser(
					retryCtx,
					postgresql.QueryCreateUserParams{
						Username:     request.Username,
						FirstName:    pgtype.Text{String: request.FullName.FirstName(), Valid: true},
						LastName:     pgtype.Text{String: request.FullName.LastName(), Valid: true},
						Email:        request.Email.Value(),
						PasswordHash: request.PasswordHash,
						StatusID:     status.ID,
						Metadata:     metadata,
					},
				)
				if createdUserErr != nil {
					slog.ErrorContext(
						retryCtx,
						fmt.Sprintf("%s Failed to create user", handlerLogPrefix),
						"username", request.Username,
						"error", createdUserErr,
					)
					return identity.User{}, fmt.Errorf("failed to create user: %w", createdUserErr)
				}

				// Map the created URL model back to an entity to return
				mappedUser, mapErr := MapUserModelToEntity(UserMapperParams{
					UserModel: createdUser,
					Status:    status.Name,
				})
				if mapErr != nil {
					slog.ErrorContext(
						retryCtx,
						fmt.Sprintf("%s Failed to map created user model to entity", handlerLogPrefix),
						"user_model", createdUser,
						"error", mapErr,
					)
					return identity.User{}, fmt.Errorf("failed to map created user model to entity: %w", mapErr)
				}
				slog.InfoContext(retryCtx, "created model", "user", mappedUser)

				return mappedUser, nil
			})
		},
		writeDatastore.config.RetryConfig,
		fmt.Sprintf("%s.Create", writeDatastore.logPrefix),
	)
}

func (writeDatastore *userWriteDatastoreAdapter) Update(ctx context.Context, userEntity identity.User) (identity.User, error) {
	handlerLogPrefix := fmt.Sprintf("%s<Update>", writeDatastore.logPrefix)
	slog.InfoContext(ctx, fmt.Sprintf("%s Updating User", handlerLogPrefix), "userId", userEntity.ID())

	operationCtx, operationCancel := context.WithTimeout(ctx, writeDatastore.config.OperationTimeout)
	defer operationCancel()

	return recoveryutils.ExecuteWithRetry(
		operationCtx,
		func(retryCtx context.Context) (identity.User, error) {
			// Use writeDatastore.withTx instead of postgres.WithTransaction directly.
			// This is the only change to the business logic — everything inside
			// the closure is identical to the original implementation.
			return writeDatastore.withTx(retryCtx, func(qtx postgresrepo.UserWriteQuerier) (identity.User, error) {
				// Check context before proceeding
				select {
				case <-retryCtx.Done():
					slog.ErrorContext(retryCtx, "Operation cancelled before validation with error", "error", retryCtx.Err())
					return identity.User{}, fmt.Errorf("operation cancelled before validation: %w", retryCtx.Err())
				default:
				}

				// query the status ID
				status, statusErr := qtx.QueryUserStatusByName(retryCtx, string(userEntity.Status()))
				if statusErr != nil {
					slog.ErrorContext(
						retryCtx,
						fmt.Sprintf("%s Failed to retrieve user status", handlerLogPrefix),
						"user_status", userEntity.Status(),
						"error", statusErr,
					)
					if errors.Is(statusErr, pgx.ErrNoRows) {
						return identity.User{}, errdefs.NotFound(statusErr)
					}

					return identity.User{}, fmt.Errorf("failed to query user status: %w", statusErr)
				}

				email := userEntity.Email()
				updatedUser, updatedUserErr := qtx.QueryUpdateUserDetails(
					retryCtx,
					postgresql.QueryUpdateUserDetailsParams{
						ID:        pgtype.UUID{Bytes: userEntity.ID(), Valid: true},
						Username:  userEntity.Username(),
						FirstName: pgtype.Text{String: userEntity.FirstName(), Valid: true},
						LastName:  pgtype.Text{String: userEntity.LastName(), Valid: true},
						Email:     email.Value(),
					},
				)
				if updatedUserErr != nil {
					slog.ErrorContext(
						retryCtx,
						fmt.Sprintf("%s Failed to update user", handlerLogPrefix),
						"username", userEntity.ID(),
						"error", updatedUserErr,
					)
					return identity.User{}, fmt.Errorf("failed to update user: %w", updatedUserErr)
				}

				// Map the created URL model back to an entity to return
				mappedUser, mapErr := MapUserModelToEntity(UserMapperParams{
					UserModel: updatedUser,
					Status:    status.Name,
				})
				if mapErr != nil {
					slog.ErrorContext(
						retryCtx,
						fmt.Sprintf("%s Failed to map updated user model to entity", handlerLogPrefix),
						"userId", updatedUser.ID,
						"error", mapErr,
					)
					return identity.User{}, fmt.Errorf("failed to map updated user model to entity: %w", mapErr)
				}
				slog.InfoContext(retryCtx, "updated model", "user", mappedUser)

				return mappedUser, nil
			})
		},
		writeDatastore.config.RetryConfig,
		fmt.Sprintf("%s.Create", writeDatastore.logPrefix),
	)
}

func (writeDatastore *userWriteDatastoreAdapter) UpdateVerification(ctx context.Context, request identity.UpdateUserVerificationRequest) (identity.User, error) {
	return identity.User{}, nil
}

// UpdateMetadata updates the metadata of a User entity based on the provided request and returns the updated User entity
func (writeDatastore *userWriteDatastoreAdapter) UpdateMetadata(ctx context.Context, request identity.UpdateUserMetadataVerificationRequest) (identity.User, error) {
	panic("not implemented")
}

// UpdatePassword updates the password of a User entity based on the provided request and returns the updated User entity
func (writeDatastore *userWriteDatastoreAdapter) UpdatePassword(ctx context.Context, request identity.UpdateUserPasswordRequest) (identity.User, error) {
	panic("not implemented")
}

// UpdateStatus updates the status of a User entity based on the provided request and returns the updated User entity
func (writeDatastore *userWriteDatastoreAdapter) UpdateStatus(ctx context.Context, request identity.UpdateUserStatusRequest) (identity.User, error) {
	panic("not implemented")
}

func (writeDatastore *userWriteDatastoreAdapter) SoftDelete(ctx context.Context, id string) error {
	panic("not implemented")
}

// Delete deletes a given entity by its ID
func (writeDatastore *userWriteDatastoreAdapter) Delete(ctx context.Context, id string) error {
	panic("not implemented")
}

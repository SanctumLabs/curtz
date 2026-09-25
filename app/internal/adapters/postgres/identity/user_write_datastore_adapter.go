package identitydatastore

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	postgresrepo "github.com/sanctumlabs/curtz/app/internal/adapters/postgres"
	postgresql "github.com/sanctumlabs/curtz/app/internal/adapters/postgres/sql"
	"github.com/sanctumlabs/curtz/app/internal/core/entity"
	"github.com/sanctumlabs/curtz/app/internal/domain/identity"
	"github.com/sanctumlabs/curtz/app/pkg/errdefs"
	"github.com/sanctumlabs/curtz/app/pkg/infra/database"
	"github.com/sanctumlabs/curtz/app/pkg/infra/database/postgres"
	"github.com/sanctumlabs/curtz/app/pkg/utils"
)

func NewUserWriteDatastoreAdapter(dbClient database.PostgresDatabaseClient, config database.Config) identity.UserWriteDatastore {
	repo := &userWriteDatastoreAdapter{
		dbClient:  dbClient,
		config:    config,
		logPrefix: "UserWriteRepoAdapter",
	}

	// Wire up the real transaction executor. This delegates to postgres.WithTransactionVoid,
	// which handles the pgxpool.Pool lifecycle. Tests override this field directly.
	repo.withTx = func(ctx context.Context, fn func(q postgresrepo.UserWriteQuerier) error) error {
		return postgres.WithTransactionVoid(ctx, dbClient, func(qtx *postgresql.Queries) error {
			// *postgresql.Queries satisfies UserWriteQuerier, so we can pass it straight through.
			return fn(qtx)
		})
	}

	return repo
}

// run executes fn inside a retried transaction, returning the resulting User.
func (writeDatastore *userWriteDatastoreAdapter) run(ctx context.Context, operation string, fn func(ctx context.Context, qtx postgresrepo.UserWriteQuerier) (identity.User, error)) (identity.User, error) {
	return execute(ctx, writeDatastore.config, fmt.Sprintf("%s.%s", writeDatastore.logPrefix, operation), writeDatastore.withTx, fn)
}

// mapUser maps a user model to an entity, logging and wrapping any mapping failure.
func mapUser(ctx context.Context, handlerLogPrefix string, userModel postgresql.User, status string) (identity.User, error) {
	mappedUser, mapErr := MapUserModelToEntity(UserMapperParams{
		UserModel: userModel,
		Status:    status,
	})
	if mapErr != nil {
		slog.ErrorContext(
			ctx,
			fmt.Sprintf("%s Failed to map user model to entity", handlerLogPrefix),
			"userId", userModel.ID,
			"error", mapErr,
		)
		return identity.User{}, fmt.Errorf("failed to map user model to entity: %w", mapErr)
	}
	return mappedUser, nil
}

func (writeDatastore *userWriteDatastoreAdapter) Save(ctx context.Context, userEntity identity.User) (identity.User, error) {
	handlerLogPrefix := fmt.Sprintf("%s<Save>", writeDatastore.logPrefix)
	slog.InfoContext(ctx, fmt.Sprintf("%s Saving User", handlerLogPrefix), "user", userEntity)

	return writeDatastore.run(ctx, "Save", func(ctx context.Context, qtx postgresrepo.UserWriteQuerier) (identity.User, error) {
		status, statusErr := queryUserStatusByName(ctx, qtx, string(userEntity.Status()))
		if statusErr != nil {
			return identity.User{}, statusErr
		}

		metadata, metadataErr := userEntity.MetadataToBytes()
		if metadataErr != nil {
			slog.WarnContext(ctx, fmt.Sprintf("%s Failed to convert user metadata to bytes", handlerLogPrefix),
				"user", userEntity,
				"error", metadataErr)
		}

		email := userEntity.Email()
		verification := userEntity.Verification()
		createdUser, createdUserErr := qtx.QueryCreateUser(
			ctx,
			postgresql.QueryCreateUserParams{
				ID:                  pgtype.UUID{Bytes: userEntity.ID(), Valid: true},
				Username:            userEntity.Username(),
				FirstName:           pgtype.Text{String: userEntity.FirstName(), Valid: true},
				LastName:            pgtype.Text{String: userEntity.LastName(), Valid: true},
				Email:               email.Value(),
				PasswordHash:        userEntity.PasswordHash(),
				StatusID:            status.ID,
				Metadata:            metadata,
				Verified:            verification.Verified(),
				VerificationToken:   pgtype.Text{String: verification.Token(), Valid: verification.Token() != ""},
				VerificationExpires: pgtype.Timestamptz{Time: verification.Expires(), Valid: !verification.Expires().IsZero()},
			},
		)
		if createdUserErr != nil {
			slog.ErrorContext(
				ctx,
				fmt.Sprintf("%s Failed to create user", handlerLogPrefix),
				"username", userEntity.Username(),
				"error", createdUserErr,
			)
			return identity.User{}, asConflict(fmt.Errorf("failed to create user: %w", createdUserErr))
		}

		return mapUser(ctx, handlerLogPrefix, createdUser, status.Name)
	})
}

func (writeDatastore *userWriteDatastoreAdapter) Create(ctx context.Context, request identity.CreateUserRequest) (identity.User, error) {
	handlerLogPrefix := fmt.Sprintf("%s<Create>", writeDatastore.logPrefix)
	slog.InfoContext(ctx, fmt.Sprintf("%s Creating User", handlerLogPrefix), "user", request)

	return writeDatastore.run(ctx, "Create", func(ctx context.Context, qtx postgresrepo.UserWriteQuerier) (identity.User, error) {
		status, statusErr := queryUserStatusByName(ctx, qtx, string(identity.UserStatusInactive))
		if statusErr != nil {
			return identity.User{}, statusErr
		}

		metadata, metadataErr := utils.MapToBytes(request.Metadata)
		if metadataErr != nil {
			slog.WarnContext(ctx, fmt.Sprintf("%s Failed to convert user metadata to bytes", handlerLogPrefix),
				"user", request,
				"error", metadataErr)
		}

		createdUser, createdUserErr := qtx.QueryCreateUser(
			ctx,
			postgresql.QueryCreateUserParams{
				ID:           pgtype.UUID{Bytes: entity.NewID(), Valid: true},
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
				ctx,
				fmt.Sprintf("%s Failed to create user", handlerLogPrefix),
				"username", request.Username,
				"error", createdUserErr,
			)
			return identity.User{}, asConflict(fmt.Errorf("failed to create user: %w", createdUserErr))
		}

		return mapUser(ctx, handlerLogPrefix, createdUser, status.Name)
	})
}

func (writeDatastore *userWriteDatastoreAdapter) Update(ctx context.Context, userEntity identity.User) (identity.User, error) {
	handlerLogPrefix := fmt.Sprintf("%s<Update>", writeDatastore.logPrefix)
	slog.InfoContext(ctx, fmt.Sprintf("%s Updating User", handlerLogPrefix), "userId", userEntity.ID())

	return writeDatastore.run(ctx, "Update", func(ctx context.Context, qtx postgresrepo.UserWriteQuerier) (identity.User, error) {
		status, statusErr := queryUserStatusByName(ctx, qtx, string(userEntity.Status()))
		if statusErr != nil {
			return identity.User{}, statusErr
		}

		email := userEntity.Email()
		updatedUser, updatedUserErr := qtx.QueryUpdateUserDetails(
			ctx,
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
				ctx,
				fmt.Sprintf("%s Failed to update user", handlerLogPrefix),
				"userId", userEntity.ID(),
				"error", updatedUserErr,
			)
			return identity.User{}, fmt.Errorf("failed to update user: %w", updatedUserErr)
		}

		return mapUser(ctx, handlerLogPrefix, updatedUser, status.Name)
	})
}

func (writeDatastore *userWriteDatastoreAdapter) UpdateVerification(ctx context.Context, request identity.UpdateUserVerificationRequest) (identity.User, error) {
	handlerLogPrefix := fmt.Sprintf("%s<UpdateVerification>", writeDatastore.logPrefix)
	slog.InfoContext(ctx, fmt.Sprintf("%s Updating User Verification", handlerLogPrefix), "userId", request.ID)

	return writeDatastore.run(ctx, "UpdateVerification", func(ctx context.Context, qtx postgresrepo.UserWriteQuerier) (identity.User, error) {
		existingUser, existingUserErr := queryUserById(ctx, qtx, request.ID)
		if existingUserErr != nil {
			return identity.User{}, existingUserErr
		}

		updatedUser, updatedUserErr := qtx.QueryUpdateUserVerification(
			ctx,
			postgresql.QueryUpdateUserVerificationParams{
				ID:                  existingUser.User.ID,
				Verified:            request.Verified,
				VerificationToken:   pgtype.Text{String: request.VerificationToken, Valid: true},
				VerificationExpires: pgtype.Timestamptz{Time: request.VerificationExpires, Valid: true},
			},
		)
		if updatedUserErr != nil {
			slog.ErrorContext(
				ctx,
				fmt.Sprintf("%s Failed to update user verification", handlerLogPrefix),
				"id", request.ID,
				"error", updatedUserErr,
			)
			return identity.User{}, fmt.Errorf("failed to update user %s verification: %w", request.ID, updatedUserErr)
		}

		return mapUser(ctx, handlerLogPrefix, updatedUser, existingUser.UserStatus.Name)
	})
}

// MarkVerified atomically flags the user verified and applies the status the aggregate
// transitioned to. Both writes share one transaction so a user can never end up verified but
// still inactive.
func (writeDatastore *userWriteDatastoreAdapter) MarkVerified(ctx context.Context, request identity.MarkUserVerifiedRequest) (identity.User, error) {
	handlerLogPrefix := fmt.Sprintf("%s<MarkVerified>", writeDatastore.logPrefix)
	slog.InfoContext(ctx, fmt.Sprintf("%s Marking User verified", handlerLogPrefix), "userId", request.ID, "status", request.Status)

	return writeDatastore.run(ctx, "MarkVerified", func(ctx context.Context, qtx postgresrepo.UserWriteQuerier) (identity.User, error) {
		existingUser, existingUserErr := queryUserById(ctx, qtx, request.ID)
		if existingUserErr != nil {
			return identity.User{}, existingUserErr
		}

		status, statusErr := queryUserStatusByName(ctx, qtx, string(request.Status))
		if statusErr != nil {
			return identity.User{}, statusErr
		}

		verifiedUser, verifiedErr := qtx.QueryUpdateUserVerification(
			ctx,
			postgresql.QueryUpdateUserVerificationParams{
				ID:                  existingUser.User.ID,
				Verified:            true,
				VerificationToken:   existingUser.User.VerificationToken,
				VerificationExpires: existingUser.User.VerificationExpires,
			},
		)
		if verifiedErr != nil {
			slog.ErrorContext(ctx, fmt.Sprintf("%s Failed to flag user verified", handlerLogPrefix), "id", request.ID, "error", verifiedErr)
			return identity.User{}, fmt.Errorf("failed to mark user %s verified: %w", request.ID, verifiedErr)
		}

		activatedUser, activatedErr := qtx.QueryUpdateUserStatusId(
			ctx,
			postgresql.QueryUpdateUserStatusIdParams{
				ID:       verifiedUser.ID,
				StatusID: status.ID,
			},
		)
		if activatedErr != nil {
			slog.ErrorContext(ctx, fmt.Sprintf("%s Failed to apply user status", handlerLogPrefix), "id", request.ID, "error", activatedErr)
			return identity.User{}, fmt.Errorf("failed to apply status to user %s: %w", request.ID, activatedErr)
		}

		return mapUser(ctx, handlerLogPrefix, activatedUser, status.Name)
	})
}

// UpdateMetadata updates the metadata of a User entity based on the provided request and returns the updated User entity
func (writeDatastore *userWriteDatastoreAdapter) UpdateMetadata(ctx context.Context, request identity.UpdateUserMetadataVerificationRequest) (identity.User, error) {
	handlerLogPrefix := fmt.Sprintf("%s<UpdateMetadata>", writeDatastore.logPrefix)
	slog.InfoContext(ctx, fmt.Sprintf("%s Updating User metadata", handlerLogPrefix), "userId", request.ID)

	return writeDatastore.run(ctx, "UpdateMetadata", func(ctx context.Context, qtx postgresrepo.UserWriteQuerier) (identity.User, error) {
		existingUser, existingUserErr := queryUserById(ctx, qtx, request.ID)
		if existingUserErr != nil {
			return identity.User{}, existingUserErr
		}

		metadata, metadataErr := utils.MapToBytes(request.Metadata)
		if metadataErr != nil {
			slog.WarnContext(ctx, fmt.Sprintf("%s Failed to convert user metadata to bytes", handlerLogPrefix),
				"user", request,
				"error", metadataErr)
		}

		updatedUser, updatedUserErr := qtx.QueryUpdateUserMetadata(
			ctx,
			postgresql.QueryUpdateUserMetadataParams{
				ID:       existingUser.User.ID,
				Metadata: metadata,
			},
		)
		if updatedUserErr != nil {
			slog.ErrorContext(
				ctx,
				fmt.Sprintf("%s Failed to update user metadata", handlerLogPrefix),
				"id", request.ID,
				"error", updatedUserErr,
			)
			return identity.User{}, fmt.Errorf("failed to update user %s metadata: %w", request.ID, updatedUserErr)
		}

		return mapUser(ctx, handlerLogPrefix, updatedUser, existingUser.UserStatus.Name)
	})
}

// UpdatePassword updates the password of a User entity based on the provided request and returns the updated User entity
func (writeDatastore *userWriteDatastoreAdapter) UpdatePassword(ctx context.Context, request identity.UpdateUserPasswordRequest) (identity.User, error) {
	handlerLogPrefix := fmt.Sprintf("%s<UpdatePassword>", writeDatastore.logPrefix)
	slog.InfoContext(ctx, fmt.Sprintf("%s Updating User password", handlerLogPrefix), "userId", request.ID)

	return writeDatastore.run(ctx, "UpdatePassword", func(ctx context.Context, qtx postgresrepo.UserWriteQuerier) (identity.User, error) {
		existingUser, existingUserErr := queryUserById(ctx, qtx, request.ID)
		if existingUserErr != nil {
			return identity.User{}, existingUserErr
		}

		updatedUser, updatedUserErr := qtx.QueryUpdateUserPassword(
			ctx,
			postgresql.QueryUpdateUserPasswordParams{
				ID:           existingUser.User.ID,
				PasswordHash: request.PasswordHash,
			},
		)
		if updatedUserErr != nil {
			slog.ErrorContext(
				ctx,
				fmt.Sprintf("%s Failed to update user password", handlerLogPrefix),
				"id", request.ID,
				"error", updatedUserErr,
			)
			return identity.User{}, fmt.Errorf("failed to update user %s password: %w", request.ID, updatedUserErr)
		}

		return mapUser(ctx, handlerLogPrefix, updatedUser, existingUser.UserStatus.Name)
	})
}

// UpdateStatus updates the status of a User entity based on the provided request and returns the updated User entity
func (writeDatastore *userWriteDatastoreAdapter) UpdateStatus(ctx context.Context, request identity.UpdateUserStatusRequest) (identity.User, error) {
	handlerLogPrefix := fmt.Sprintf("%s<UpdateStatus>", writeDatastore.logPrefix)
	slog.InfoContext(ctx, fmt.Sprintf("%s Updating User status", handlerLogPrefix), "userId", request.ID, "status", request.Status)

	return writeDatastore.run(ctx, "UpdateStatus", func(ctx context.Context, qtx postgresrepo.UserWriteQuerier) (identity.User, error) {
		status, statusErr := queryUserStatusByName(ctx, qtx, string(request.Status))
		if statusErr != nil {
			return identity.User{}, statusErr
		}

		userUUID, userUUIDErr := postgres.StringToUUID(request.ID)
		if userUUIDErr != nil {
			return identity.User{}, fmt.Errorf("failed to convert user ID to UUID: %w", userUUIDErr)
		}

		updatedUser, updatedUserErr := qtx.QueryUpdateUserStatusId(
			ctx,
			postgresql.QueryUpdateUserStatusIdParams{
				ID:       userUUID,
				StatusID: status.ID,
			},
		)
		if updatedUserErr != nil {
			slog.ErrorContext(
				ctx,
				fmt.Sprintf("%s Failed to update user status", handlerLogPrefix),
				"id", request.ID,
				"error", updatedUserErr,
			)
			if errors.Is(updatedUserErr, pgx.ErrNoRows) {
				return identity.User{}, errdefs.NotFound(updatedUserErr)
			}
			return identity.User{}, fmt.Errorf("failed to update user %s status: %w", request.ID, updatedUserErr)
		}

		return mapUser(ctx, handlerLogPrefix, updatedUser, status.Name)
	})
}

// SoftDelete marks a User as deleted without removing the record
func (writeDatastore *userWriteDatastoreAdapter) SoftDelete(ctx context.Context, id string) error {
	handlerLogPrefix := fmt.Sprintf("%s<SoftDelete>", writeDatastore.logPrefix)
	slog.InfoContext(ctx, fmt.Sprintf("%s Soft deleting User", handlerLogPrefix), "userId", id)

	_, err := writeDatastore.run(ctx, "SoftDelete", func(ctx context.Context, qtx postgresrepo.UserWriteQuerier) (identity.User, error) {
		userUUID, userUUIDErr := postgres.StringToUUID(id)
		if userUUIDErr != nil {
			return identity.User{}, fmt.Errorf("failed to convert user ID to UUID: %w", userUUIDErr)
		}

		_, deleteErr := qtx.QuerySoftDeleteUser(
			ctx,
			postgresql.QuerySoftDeleteUserParams{
				ID:        userUUID,
				DeletedAt: pgtype.Timestamptz{Time: time.Now().UTC(), Valid: true},
			},
		)
		if deleteErr != nil {
			slog.ErrorContext(
				ctx,
				fmt.Sprintf("%s Failed to soft delete user", handlerLogPrefix),
				"id", id,
				"error", deleteErr,
			)
			if errors.Is(deleteErr, pgx.ErrNoRows) {
				return identity.User{}, errdefs.NotFound(deleteErr)
			}
			return identity.User{}, fmt.Errorf("failed to soft delete user %s: %w", id, deleteErr)
		}

		return identity.User{}, nil
	})
	return err
}

// Delete permanently deletes a User by its ID
func (writeDatastore *userWriteDatastoreAdapter) Delete(ctx context.Context, id string) error {
	handlerLogPrefix := fmt.Sprintf("%s<Delete>", writeDatastore.logPrefix)
	slog.InfoContext(ctx, fmt.Sprintf("%s Deleting User", handlerLogPrefix), "userId", id)

	_, err := writeDatastore.run(ctx, "Delete", func(ctx context.Context, qtx postgresrepo.UserWriteQuerier) (identity.User, error) {
		userUUID, userUUIDErr := postgres.StringToUUID(id)
		if userUUIDErr != nil {
			return identity.User{}, fmt.Errorf("failed to convert user ID to UUID: %w", userUUIDErr)
		}

		if _, deleteErr := qtx.QueryDeleteUserWithId(ctx, userUUID); deleteErr != nil {
			slog.ErrorContext(
				ctx,
				fmt.Sprintf("%s Failed to delete user", handlerLogPrefix),
				"id", id,
				"error", deleteErr,
			)
			if errors.Is(deleteErr, pgx.ErrNoRows) {
				return identity.User{}, errdefs.NotFound(deleteErr)
			}
			return identity.User{}, fmt.Errorf("failed to delete user %s: %w", id, deleteErr)
		}

		return identity.User{}, nil
	})
	return err
}

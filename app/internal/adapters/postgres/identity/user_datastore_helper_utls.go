package identitydatastore

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
	postgresrepo "github.com/sanctumlabs/curtz/app/internal/adapters/postgres"
	postgresql "github.com/sanctumlabs/curtz/app/internal/adapters/postgres/sql"
	"github.com/sanctumlabs/curtz/app/pkg/errdefs"
	"github.com/sanctumlabs/curtz/app/pkg/infra/database/postgres"
)

// queryUserStatusByName queries the database for a user status by its name using the provided UserWriteQuerier.
func queryUserStatusByName(ctx context.Context, qtx postgresrepo.UserWriteQuerier, statusName string) (*postgresql.UserStatus, error) {
	existingStatus, statusErr := qtx.QueryUserStatusByName(ctx, statusName)
	if statusErr != nil {
		slog.ErrorContext(
			ctx,
			"Failed to retrieve user status",
			"user_status", statusName,
			"error", statusErr,
		)
		if errors.Is(statusErr, pgx.ErrNoRows) {
			return nil, errdefs.NotFound(statusErr)
		}

		return nil, fmt.Errorf("failed to query user status: %w", statusErr)
	}
	return &existingStatus, nil
}

// queryUserById is a helper function to query a user by ID
func queryUserById(ctx context.Context, qtx postgresrepo.UserReadQuerier, userId string) (*postgresql.QueryUserByIdRow, error) {
	userUUID, userUUIDErr := postgres.StringToUUID(userId)
	if userUUIDErr != nil {
		return nil, fmt.Errorf("failed to convert user ID to UUID: %w", userUUIDErr)
	}
	existingUser, existingUserErr := qtx.QueryUserById(ctx, userUUID)
	if existingUserErr != nil {
		slog.ErrorContext(
			ctx,
			"Failed to retrieve user",
			"id", userId,
			"error", existingUserErr,
		)
		if errors.Is(existingUserErr, pgx.ErrNoRows) {
			return nil, errdefs.NotFound(existingUserErr)
		}

		return nil, fmt.Errorf("failed to query user: %s %w", userId, existingUserErr)
	}
	return &existingUser, nil
}

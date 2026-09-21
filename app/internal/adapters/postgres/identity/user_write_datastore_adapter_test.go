package identitydatastore

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/jackc/pgx/v5"
	"github.com/onsi/ginkgo/v2"
	mockpostgresrepo "github.com/sanctumlabs/curtz/app/internal/adapters/postgres/mocks"
	postgresql "github.com/sanctumlabs/curtz/app/internal/adapters/postgres/sql"
	mockpostgresql "github.com/sanctumlabs/curtz/app/internal/adapters/postgres/sql/mocks"
	"github.com/sanctumlabs/curtz/app/internal/core/entity"
	"github.com/sanctumlabs/curtz/app/internal/domain/identity"
	mockidentity "github.com/sanctumlabs/curtz/app/internal/domain/identity/mocks"
	"github.com/sanctumlabs/curtz/app/pkg/errdefs"
	"github.com/sanctumlabs/curtz/app/pkg/infra/database"
	mockdatabase "github.com/sanctumlabs/curtz/app/pkg/infra/database/mocks"
	recoveryutils "github.com/sanctumlabs/curtz/app/pkg/utils/recover"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

var _ = ginkgo.Describe("User Write Datastore Adapter Unit Test Suite", ginkgo.Ordered, func() {
	ctx := context.Background()
	var (
		mockCtrl             *gomock.Controller
		mockDbClient         *mockdatabase.MockPostgresDatabaseClient
		mockUserWriteQuerier *mockpostgresrepo.MockUserWriteQuerier
		userWriteDatastore   *userWriteDatastoreAdapter
	)

	ginkgo.BeforeAll(func() {
		config := database.Config{
			OperationTimeout: 30 * time.Second,
			RetryConfig:      recoveryutils.DefaultRetryConfig,
		}
		mockCtrl = gomock.NewController(ginkgo.GinkgoT())
		mockDbClient = mockdatabase.NewMockPostgresDatabaseClient(mockCtrl)
		mockUserWriteQuerier = mockpostgresrepo.NewMockUserWriteQuerier(mockCtrl)
		userWriteDatastore = &userWriteDatastoreAdapter{
			logPrefix: "UserWriteRepoAdapter",
			dbClient:  mockDbClient,
			config:    config,
		}

		injectMockUserWriteTx(userWriteDatastore, mockUserWriteQuerier)
	})

	ginkgo.AfterEach(func() {
		mockCtrl.Finish()
	})

	mockUserId := entity.NewID()
	mockUser, mockUserErr := mockidentity.MockUser(
		mockidentity.WithId(mockUserId),
	)
	assert.NoError(ginkgo.GinkgoT(), mockUserErr)

	mockUserRecord := mockpostgresql.MockUser(mockpostgresql.WithUser(*mockUser))
	mockUserStatus := mockpostgresql.MockUserStatus(identity.UserStatusActive)
	mockWUserQueryByIdRow := postgresql.QueryUserByIdRow{
		User:       mockUserRecord,
		UserStatus: mockUserStatus,
	}

	ginkgo.Describe("Save", func() {
		ginkgo.It("saves a new user successfully", func() {
			mockUserWriteQuerier.
				EXPECT().
				QueryUserStatusByName(gomock.Any(), gomock.Any()).
				Return(mockUserStatus, nil).
				Times(1)

			mockUserWriteQuerier.
				EXPECT().
				QueryCreateUser(gomock.Any(), gomock.Any()).
				Return(mockUserRecord, nil).
				Times(1)

			actual, actualErr := userWriteDatastore.Save(ctx, *mockUser)
			assert.Nil(ginkgo.GinkgoT(), actualErr)
			assert.Equal(ginkgo.GinkgoT(), mockUser.ID(), actual.ID())
			assert.Equal(ginkgo.GinkgoT(), mockUser.Username(), actual.Username())
			assert.Equal(ginkgo.GinkgoT(), mockUser.Email(), actual.Email())
			assert.Equal(ginkgo.GinkgoT(), mockUser.FirstName(), actual.FirstName())
			assert.Equal(ginkgo.GinkgoT(), mockUser.LastName(), actual.LastName())
			assert.Equal(ginkgo.GinkgoT(), mockUser.CreatedAt(), actual.CreatedAt())
			assert.Equal(ginkgo.GinkgoT(), mockUser.UpdatedAt(), actual.UpdatedAt())
		})
	})

	ginkgo.Describe("Create", func() {
		ginkgo.It("creates a new user successfully and returns the user", func() {
			username := mockUser.Username()
			fullName := mockUser.FullName()
			createNewUserRequest := identity.CreateUserRequest{
				Username:     username,
				FullName:     fullName,
				Email:        mockUser.Email(),
				PasswordHash: mockUser.PasswordHash(),
				Metadata:     mockUser.Metadata(),
			}

			mockUserWriteQuerier.
				EXPECT().
				QueryUserStatusByName(gomock.Any(), gomock.Any()).
				Return(mockUserStatus, nil).
				Times(1)

			mockUserWriteQuerier.
				EXPECT().
				QueryCreateUser(gomock.Any(), gomock.Any()).
				Return(mockUserRecord, nil).
				Times(1)

			actual, actualErr := userWriteDatastore.Create(ctx, createNewUserRequest)
			assert.Nil(ginkgo.GinkgoT(), actualErr)
			assert.Equal(ginkgo.GinkgoT(), mockUser.ID(), actual.ID())
			assert.Equal(ginkgo.GinkgoT(), mockUser.Username(), actual.Username())
			assert.Equal(ginkgo.GinkgoT(), mockUser.Email(), actual.Email())
			assert.Equal(ginkgo.GinkgoT(), mockUser.FirstName(), actual.FirstName())
			assert.Equal(ginkgo.GinkgoT(), mockUser.LastName(), actual.LastName())
			assert.Equal(ginkgo.GinkgoT(), mockUser.CreatedAt(), actual.CreatedAt())
			assert.Equal(ginkgo.GinkgoT(), mockUser.UpdatedAt(), actual.UpdatedAt())
		})

		ginkgo.It("returns error if the status does not exist", func() {
			username := mockUser.Username()
			fullName := mockUser.FullName()
			createNewUserRequest := identity.CreateUserRequest{
				Username:     username,
				FullName:     fullName,
				Email:        mockUser.Email(),
				PasswordHash: mockUser.PasswordHash(),
				Metadata:     mockUser.Metadata(),
			}

			statusErr := fmt.Errorf("failed to retrieve user status")
			mockUserWriteQuerier.
				EXPECT().
				QueryUserStatusByName(gomock.Any(), gomock.Any()).
				Return(postgresql.UserStatus{}, statusErr).
				Times(1)

			mockUserWriteQuerier.
				EXPECT().
				QueryCreateUser(gomock.Any(), gomock.Any()).
				Return(mockUserRecord, nil).
				Times(0)

			actual, actualErr := userWriteDatastore.Create(ctx, createNewUserRequest)
			assert.NotNil(ginkgo.GinkgoT(), actualErr)
			assert.Empty(ginkgo.GinkgoT(), actual)
		})
	})

	ginkgo.Describe("Update", func() {
		ginkgo.It("successfully updates a user and returns the updated user", func() {
			mockUserStatusActive := mockpostgresql.MockUserStatus(identity.UserStatusActive)

			mockUserWriteQuerier.
				EXPECT().
				QueryUserStatusByName(gomock.Any(), gomock.Any()).
				Return(mockUserStatusActive, nil).
				Times(1)

			mockUserWriteQuerier.
				EXPECT().
				QueryUpdateUserDetails(gomock.Any(), gomock.Any()).
				Return(mockUserRecord, nil).
				Times(1)

			actual, actualErr := userWriteDatastore.Update(ctx, *mockUser)
			assert.Nil(ginkgo.GinkgoT(), actualErr)
			assert.Equal(ginkgo.GinkgoT(), mockUser.ID(), actual.ID())
			assert.Equal(ginkgo.GinkgoT(), mockUser.Username(), actual.Username())
			assert.Equal(ginkgo.GinkgoT(), mockUser.Email(), actual.Email())
			assert.Equal(ginkgo.GinkgoT(), mockUser.FirstName(), actual.FirstName())
			assert.Equal(ginkgo.GinkgoT(), mockUser.LastName(), actual.LastName())
			assert.Equal(ginkgo.GinkgoT(), mockUser.CreatedAt(), actual.CreatedAt())
			assert.Equal(ginkgo.GinkgoT(), mockUser.UpdatedAt(), actual.UpdatedAt())
		})

		ginkgo.It("returns error when there is a failure to update a user", func() {
			mockUserStatusActive := mockpostgresql.MockUserStatus(identity.UserStatusActive)

			mockUserWriteQuerier.
				EXPECT().
				QueryUserStatusByName(gomock.Any(), gomock.Any()).
				Return(mockUserStatusActive, nil).
				Times(1)

			updateUserDetailsErr := errors.New("Failed to update user")
			mockUserWriteQuerier.
				EXPECT().
				QueryUpdateUserDetails(gomock.Any(), gomock.Any()).
				Return(postgresql.User{}, updateUserDetailsErr).
				Times(1)

			actual, actualErr := userWriteDatastore.Update(ctx, *mockUser)
			assert.NotNil(ginkgo.GinkgoT(), actualErr)
			assert.Empty(ginkgo.GinkgoT(), actual)
		})
	})

	ginkgo.Describe("UpdateVerification", func() {
		verifiedVerificationRequest := identity.UpdateUserVerificationRequest{
			ID:                  mockUserId.String(),
			Verified:            true,
			VerificationToken:   faker.UUIDHyphenated(),
			VerificationExpires: time.Now(),
		}

		ginkgo.It("successfully updates user verification returns the updated user", func() {
			mockUserWriteQuerier.
				EXPECT().
				QueryUserById(gomock.Any(), gomock.Any()).
				Return(mockWUserQueryByIdRow, nil).
				Times(1)

			mockUserWriteQuerier.
				EXPECT().
				QueryUpdateUserVerification(gomock.Any(), gomock.Any()).
				Return(mockUserRecord, nil).
				Times(1)

			actual, actualErr := userWriteDatastore.UpdateVerification(ctx, verifiedVerificationRequest)
			assert.Nil(ginkgo.GinkgoT(), actualErr)
			assert.Equal(ginkgo.GinkgoT(), mockUser.ID(), actual.ID())
			assert.Equal(ginkgo.GinkgoT(), mockUser.Username(), actual.Username())
			assert.Equal(ginkgo.GinkgoT(), mockUser.Email(), actual.Email())
			assert.Equal(ginkgo.GinkgoT(), mockUser.FirstName(), actual.FirstName())
			assert.Equal(ginkgo.GinkgoT(), mockUser.LastName(), actual.LastName())
			assert.Equal(ginkgo.GinkgoT(), mockUser.CreatedAt(), actual.CreatedAt())
			assert.Equal(ginkgo.GinkgoT(), mockUser.UpdatedAt(), actual.UpdatedAt())
			assert.Equal(ginkgo.GinkgoT(), mockUser.Verification(), actual.Verification())
		})

		ginkgo.It("returns error when there is a failure to update user verification", func() {
			mockUserWriteQuerier.
				EXPECT().
				QueryUserById(gomock.Any(), gomock.Any()).
				Return(mockWUserQueryByIdRow, nil).
				Times(1)

			updateErr := errors.New("failed to update user verification")
			mockUserWriteQuerier.
				EXPECT().
				QueryUpdateUserVerification(gomock.Any(), gomock.Any()).
				Return(postgresql.User{}, updateErr).
				Times(1)

			actual, actualErr := userWriteDatastore.UpdateVerification(ctx, verifiedVerificationRequest)
			assert.NotNil(ginkgo.GinkgoT(), actualErr)
			assert.Empty(ginkgo.GinkgoT(), actual)
		})
	})

	ginkgo.Describe("UpdateMetadata", func() {
		updateMetadataRequest := identity.UpdateUserMetadataVerificationRequest{
			ID: mockUserId.String(),
			Metadata: map[string]interface{}{
				"key1": "value1",
				"key2": 42,
				"key3": true,
			},
		}

		ginkgo.It("successfully updates user metadata returns the updated user", func() {
			mockUserWriteQuerier.
				EXPECT().
				QueryUserById(gomock.Any(), gomock.Any()).
				Return(mockWUserQueryByIdRow, nil).
				Times(1)

			mockUserWriteQuerier.
				EXPECT().
				QueryUpdateUserMetadata(gomock.Any(), gomock.Any()).
				Return(mockUserRecord, nil).
				Times(1)

			actual, actualErr := userWriteDatastore.UpdateMetadata(ctx, updateMetadataRequest)
			assert.Nil(ginkgo.GinkgoT(), actualErr)
			assert.Equal(ginkgo.GinkgoT(), mockUser.ID(), actual.ID())
			assert.Equal(ginkgo.GinkgoT(), mockUser.Username(), actual.Username())
			assert.Equal(ginkgo.GinkgoT(), mockUser.Email(), actual.Email())
			assert.Equal(ginkgo.GinkgoT(), mockUser.FirstName(), actual.FirstName())
			assert.Equal(ginkgo.GinkgoT(), mockUser.LastName(), actual.LastName())
			assert.Equal(ginkgo.GinkgoT(), mockUser.CreatedAt(), actual.CreatedAt())
			assert.Equal(ginkgo.GinkgoT(), mockUser.UpdatedAt(), actual.UpdatedAt())
			assert.Equal(ginkgo.GinkgoT(), mockUser.Verification(), actual.Verification())
			assert.Equal(ginkgo.GinkgoT(), mockUser.Metadata(), actual.Metadata())
		})

		ginkgo.It("returns error when there is a failure to update user metadata", func() {
			mockUserWriteQuerier.
				EXPECT().
				QueryUserById(gomock.Any(), gomock.Any()).
				Return(mockWUserQueryByIdRow, nil).
				Times(1)

			updateErr := errors.New("failed to update user metadata")
			mockUserWriteQuerier.
				EXPECT().
				QueryUpdateUserMetadata(gomock.Any(), gomock.Any()).
				Return(postgresql.User{}, updateErr).
				Times(1)

			actual, actualErr := userWriteDatastore.UpdateMetadata(ctx, updateMetadataRequest)
			assert.NotNil(ginkgo.GinkgoT(), actualErr)
			assert.Empty(ginkgo.GinkgoT(), actual)
		})
	})
	ginkgo.Describe("UpdatePassword", func() {
		updatePasswordRequest := identity.UpdateUserPasswordRequest{
			ID:           mockUserId.String(),
			PasswordHash: "new-hash",
		}

		ginkgo.It("successfully updates the password and returns the updated user", func() {
			mockUserWriteQuerier.
				EXPECT().
				QueryUserById(gomock.Any(), gomock.Any()).
				Return(mockWUserQueryByIdRow, nil).
				Times(1)

			updatedRecord := mockpostgresql.MockUser(
				mockpostgresql.WithUser(*mockUser),
				mockpostgresql.UserWithPasswordHash("new-hash"),
			)
			mockUserWriteQuerier.
				EXPECT().
				QueryUpdateUserPassword(gomock.Any(), postgresql.QueryUpdateUserPasswordParams{
					ID:           mockUserRecord.ID,
					PasswordHash: "new-hash",
				}).
				Return(updatedRecord, nil).
				Times(1)

			actual, actualErr := userWriteDatastore.UpdatePassword(ctx, updatePasswordRequest)
			assert.Nil(ginkgo.GinkgoT(), actualErr)
			assert.Equal(ginkgo.GinkgoT(), mockUser.ID(), actual.ID())
			assert.Equal(ginkgo.GinkgoT(), "new-hash", actual.PasswordHash())
		})

		ginkgo.It("returns not found when the user does not exist", func() {
			mockUserWriteQuerier.
				EXPECT().
				QueryUserById(gomock.Any(), gomock.Any()).
				Return(postgresql.QueryUserByIdRow{}, pgx.ErrNoRows).
				Times(1)

			actual, actualErr := userWriteDatastore.UpdatePassword(ctx, updatePasswordRequest)
			assert.True(ginkgo.GinkgoT(), errdefs.IsNotFound(actualErr), "expected NotFound, got %v", actualErr)
			assert.Empty(ginkgo.GinkgoT(), actual)
		})

		ginkgo.It("returns error when there is a failure to update the password", func() {
			mockUserWriteQuerier.
				EXPECT().
				QueryUserById(gomock.Any(), gomock.Any()).
				Return(mockWUserQueryByIdRow, nil).
				Times(1)

			mockUserWriteQuerier.
				EXPECT().
				QueryUpdateUserPassword(gomock.Any(), gomock.Any()).
				Return(postgresql.User{}, errors.New("failed to update password")).
				Times(1)

			actual, actualErr := userWriteDatastore.UpdatePassword(ctx, updatePasswordRequest)
			assert.NotNil(ginkgo.GinkgoT(), actualErr)
			assert.Empty(ginkgo.GinkgoT(), actual)
		})
	})

	ginkgo.Describe("UpdateStatus", func() {
		updateStatusRequest := identity.UpdateUserStatusRequest{
			ID:     mockUserId.String(),
			Status: identity.UserStatusSuspended,
		}
		suspendedStatus := mockpostgresql.MockUserStatus(identity.UserStatusSuspended)

		ginkgo.It("successfully updates the status and returns the updated user", func() {
			mockUserWriteQuerier.
				EXPECT().
				QueryUserStatusByName(gomock.Any(), string(identity.UserStatusSuspended)).
				Return(suspendedStatus, nil).
				Times(1)

			mockUserWriteQuerier.
				EXPECT().
				QueryUpdateUserStatusId(gomock.Any(), postgresql.QueryUpdateUserStatusIdParams{
					ID:       mockUserRecord.ID,
					StatusID: suspendedStatus.ID,
				}).
				Return(mockUserRecord, nil).
				Times(1)

			actual, actualErr := userWriteDatastore.UpdateStatus(ctx, updateStatusRequest)
			assert.Nil(ginkgo.GinkgoT(), actualErr)
			assert.Equal(ginkgo.GinkgoT(), mockUser.ID(), actual.ID())
			assert.Equal(ginkgo.GinkgoT(), identity.UserStatusSuspended, actual.Status())
		})

		ginkgo.It("returns not found when the status does not exist", func() {
			mockUserWriteQuerier.
				EXPECT().
				QueryUserStatusByName(gomock.Any(), gomock.Any()).
				Return(postgresql.UserStatus{}, pgx.ErrNoRows).
				Times(1)

			actual, actualErr := userWriteDatastore.UpdateStatus(ctx, updateStatusRequest)
			assert.True(ginkgo.GinkgoT(), errdefs.IsNotFound(actualErr), "expected NotFound, got %v", actualErr)
			assert.Empty(ginkgo.GinkgoT(), actual)
		})

		ginkgo.It("returns not found when the user does not exist", func() {
			mockUserWriteQuerier.
				EXPECT().
				QueryUserStatusByName(gomock.Any(), gomock.Any()).
				Return(suspendedStatus, nil).
				Times(1)

			mockUserWriteQuerier.
				EXPECT().
				QueryUpdateUserStatusId(gomock.Any(), gomock.Any()).
				Return(postgresql.User{}, pgx.ErrNoRows).
				Times(1)

			actual, actualErr := userWriteDatastore.UpdateStatus(ctx, updateStatusRequest)
			assert.True(ginkgo.GinkgoT(), errdefs.IsNotFound(actualErr), "expected NotFound, got %v", actualErr)
			assert.Empty(ginkgo.GinkgoT(), actual)
		})
	})

	ginkgo.Describe("SoftDelete", func() {
		ginkgo.It("marks the user as deleted", func() {
			mockUserWriteQuerier.
				EXPECT().
				QuerySoftDeleteUser(gomock.Any(), gomock.Cond(func(p postgresql.QuerySoftDeleteUserParams) bool {
					return p.ID == mockUserRecord.ID && p.DeletedAt.Valid && !p.DeletedAt.Time.IsZero()
				})).
				Return(mockUserRecord, nil).
				Times(1)

			actualErr := userWriteDatastore.SoftDelete(ctx, mockUserId.String())
			assert.Nil(ginkgo.GinkgoT(), actualErr)
		})

		ginkgo.It("returns not found when the user does not exist", func() {
			mockUserWriteQuerier.
				EXPECT().
				QuerySoftDeleteUser(gomock.Any(), gomock.Any()).
				Return(postgresql.User{}, pgx.ErrNoRows).
				Times(1)

			actualErr := userWriteDatastore.SoftDelete(ctx, mockUserId.String())
			assert.True(ginkgo.GinkgoT(), errdefs.IsNotFound(actualErr), "expected NotFound, got %v", actualErr)
		})

		ginkgo.It("returns error for an invalid id without hitting the database", func() {
			actualErr := userWriteDatastore.SoftDelete(ctx, "not-a-uuid")
			assert.NotNil(ginkgo.GinkgoT(), actualErr)
		})
	})

	ginkgo.Describe("Delete", func() {
		ginkgo.It("permanently deletes the user", func() {
			mockUserWriteQuerier.
				EXPECT().
				QueryDeleteUserWithId(gomock.Any(), mockUserRecord.ID).
				Return(mockUserRecord, nil).
				Times(1)

			actualErr := userWriteDatastore.Delete(ctx, mockUserId.String())
			assert.Nil(ginkgo.GinkgoT(), actualErr)
		})

		ginkgo.It("returns not found when the user does not exist", func() {
			mockUserWriteQuerier.
				EXPECT().
				QueryDeleteUserWithId(gomock.Any(), gomock.Any()).
				Return(postgresql.User{}, pgx.ErrNoRows).
				Times(1)

			actualErr := userWriteDatastore.Delete(ctx, mockUserId.String())
			assert.True(ginkgo.GinkgoT(), errdefs.IsNotFound(actualErr), "expected NotFound, got %v", actualErr)
		})
	})
})

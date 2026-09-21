//go:build integration

package identitydatastore

import (
	"context"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/onsi/ginkgo/v2"
	"github.com/sanctumlabs/curtz/app/internal/core/entity"
	"github.com/sanctumlabs/curtz/app/internal/domain/identity"
	mockidentity "github.com/sanctumlabs/curtz/app/internal/domain/identity/mocks"
	"github.com/sanctumlabs/curtz/app/pkg/errdefs"
	"github.com/sanctumlabs/curtz/app/pkg/infra/database"
	recoveryutils "github.com/sanctumlabs/curtz/app/pkg/utils/recover"
	"github.com/sanctumlabs/curtz/app/test"
	"github.com/stretchr/testify/assert"
)

var _ = ginkgo.Describe("User Write Datastore Adapter Integration Test Suite", ginkgo.Ordered, func() {
	ctx := context.Background()

	var (
		testPostgresDatabaseClient database.PostgresDatabaseClient
		userWriteDatastoreAdapter  identity.UserWriteDatastore
	)

	ginkgo.BeforeAll(func() {
		var err error
		testPostgresDatabaseClient, err = test.TestPostgresDatabaseClient(ctx)
		if err != nil {
			assert.FailNow(ginkgo.GinkgoT(), "failed to create database client: %s", err.Error())
		}

		config := database.Config{
			OperationTimeout: 5 * time.Minute,
			RetryConfig:      recoveryutils.DefaultRetryConfig,
		}

		userWriteDatastoreAdapter = NewUserWriteDatastoreAdapter(
			testPostgresDatabaseClient,
			config,
		)
	})

	ginkgo.AfterAll(func() {
		testPostgresDatabaseClient.Close()
	})

	ginkgo.Describe("Save", func() {
		ginkgo.It("saves a new user successfully", func() {
			mockUser, mockUserErr := mockidentity.MockUser()
			assert.NoError(ginkgo.GinkgoT(), mockUserErr)

			// Require stops the test immediately on failure, preventing a nil
			// dereference on *mockUser in the Create call below.
			assert.NoError(ginkgo.GinkgoT(), mockUserErr)

			actual, actualErr := userWriteDatastoreAdapter.Save(ctx, *mockUser)
			assert.NoError(ginkgo.GinkgoT(), actualErr)
			assert.NotNil(ginkgo.GinkgoT(), actual)
		})
	})

	ginkgo.Describe("Create", func() {
		ginkgo.It("creates a new user successfully", func() {
			mockUser, mockUserErr := mockidentity.MockUser()
			assert.NoError(ginkgo.GinkgoT(), mockUserErr)

			username := mockUser.Username()
			fullName := mockUser.FullName()
			createNewUserRequest := identity.CreateUserRequest{
				Username:     username,
				FullName:     fullName,
				Email:        mockUser.Email(),
				PasswordHash: mockUser.PasswordHash(),
				Metadata:     mockUser.Metadata(),
			}

			actual, actualErr := userWriteDatastoreAdapter.Create(ctx, createNewUserRequest)
			assert.NoError(ginkgo.GinkgoT(), actualErr)
			assert.NotNil(ginkgo.GinkgoT(), actual)
		})
	})

	ginkgo.Describe("Update", func() {
		mockUser, mockUserErr := mockidentity.MockUser()
		assert.NoError(ginkgo.GinkgoT(), mockUserErr)
		username := mockUser.Username()
		fullName := mockUser.FullName()
		createNewUserRequest := identity.CreateUserRequest{
			Username:     username,
			FullName:     fullName,
			Email:        mockUser.Email(),
			PasswordHash: mockUser.PasswordHash(),
			Metadata:     mockUser.Metadata(),
		}

		ginkgo.It("updates an existing user successfully", func() {
			actualCreatedUser, createError := userWriteDatastoreAdapter.Create(ctx, createNewUserRequest)
			assert.NoError(ginkgo.GinkgoT(), createError)
			assert.NotNil(ginkgo.GinkgoT(), actualCreatedUser)

			updatedEmail, updatedEmailErr := identity.NewEmail(faker.Email())
			assert.NoError(ginkgo.GinkgoT(), updatedEmailErr)

			newUser := actualCreatedUser.WithEmail(updatedEmail)

			actual, actualErr := userWriteDatastoreAdapter.Update(ctx, newUser)
			assert.NoError(ginkgo.GinkgoT(), actualErr)
			assert.NotNil(ginkgo.GinkgoT(), actual)

			actualEmail := actual.Email()
			assert.Equal(ginkgo.GinkgoT(), updatedEmail.Value(), actualEmail.Value())
		})

		ginkgo.It("returns error if there is a failure to update a user", func() {
			mockUser, mockUserErr := mockidentity.MockUser()
			assert.NoError(ginkgo.GinkgoT(), mockUserErr)

			actual, actualErr := userWriteDatastoreAdapter.Update(ctx, *mockUser)
			assert.Error(ginkgo.GinkgoT(), actualErr)
			assert.Empty(ginkgo.GinkgoT(), actual)
		})

		ginkgo.Describe("Metadata", func() {
			anotherMockUser, anotherMockUserErr := mockidentity.MockUser(
				mockidentity.WithUsername(faker.Username()),
				mockidentity.WithEmail(faker.Email()),
			)
			assert.NoError(ginkgo.GinkgoT(), anotherMockUserErr)

			anotherNewUserRequest := identity.CreateUserRequest{
				Username:     anotherMockUser.Username(),
				FullName:     anotherMockUser.FullName(),
				Email:        anotherMockUser.Email(),
				PasswordHash: anotherMockUser.PasswordHash(),
				Metadata:     anotherMockUser.Metadata(),
			}

			ginkgo.It("updates an existing user metadata successfully", func() {
				actualAnotherCreatedUser, anotherCreatedError := userWriteDatastoreAdapter.Create(ctx, anotherNewUserRequest)
				assert.NoError(ginkgo.GinkgoT(), anotherCreatedError)
				assert.NotNil(ginkgo.GinkgoT(), actualAnotherCreatedUser)

				updateMetadataRequest := identity.UpdateUserMetadataVerificationRequest{
					ID: actualAnotherCreatedUser.ID().String(),
					Metadata: map[string]any{
						"key1": "value1",
						"key2": 42,
						"key3": true,
					},
				}

				actual, actualErr := userWriteDatastoreAdapter.UpdateMetadata(ctx, updateMetadataRequest)
				assert.NoError(ginkgo.GinkgoT(), actualErr)
				assert.NotNil(ginkgo.GinkgoT(), actual)

				actualMetadata := actual.Metadata()
				// metadata round-trips through JSONB, so numbers come back as float64
				assert.EqualValues(ginkgo.GinkgoT(), updateMetadataRequest.Metadata["key1"], actualMetadata["key1"])
				assert.EqualValues(ginkgo.GinkgoT(), updateMetadataRequest.Metadata["key2"], actualMetadata["key2"])
				assert.EqualValues(ginkgo.GinkgoT(), updateMetadataRequest.Metadata["key3"], actualMetadata["key3"])
			})

			ginkgo.It("returns not found when the user does not exist", func() {
				actual, actualErr := userWriteDatastoreAdapter.UpdateMetadata(ctx, identity.UpdateUserMetadataVerificationRequest{
					ID:       entity.NewID().String(),
					Metadata: map[string]any{"key": "value"},
				})
				assert.True(ginkgo.GinkgoT(), errdefs.IsNotFound(actualErr), "expected NotFound, got %v", actualErr)
				assert.Empty(ginkgo.GinkgoT(), actual)
			})
		})
	})

	// createUser persists a fresh user with a unique username/email and fails the spec if that fails.
	createUser := func() identity.User {
		mockUser, mockUserErr := mockidentity.MockUser(
			mockidentity.WithUsername(faker.Username()),
			mockidentity.WithEmail(faker.Email()),
		)
		assert.NoError(ginkgo.GinkgoT(), mockUserErr)

		created, createErr := userWriteDatastoreAdapter.Create(ctx, identity.CreateUserRequest{
			Username:     mockUser.Username(),
			FullName:     mockUser.FullName(),
			Email:        mockUser.Email(),
			PasswordHash: mockUser.PasswordHash(),
			Metadata:     mockUser.Metadata(),
		})
		if createErr != nil {
			assert.FailNow(ginkgo.GinkgoT(), "failed to create user: %s", createErr.Error())
		}
		return created
	}

	ginkgo.Describe("UpdatePassword", func() {
		ginkgo.It("updates the password hash of an existing user", func() {
			created := createUser()

			actual, actualErr := userWriteDatastoreAdapter.UpdatePassword(ctx, identity.UpdateUserPasswordRequest{
				ID:           created.ID().String(),
				PasswordHash: "rotated-hash",
			})
			assert.NoError(ginkgo.GinkgoT(), actualErr)
			assert.Equal(ginkgo.GinkgoT(), created.ID(), actual.ID())
			assert.Equal(ginkgo.GinkgoT(), "rotated-hash", actual.PasswordHash())
		})

		ginkgo.It("returns not found when the user does not exist", func() {
			_, actualErr := userWriteDatastoreAdapter.UpdatePassword(ctx, identity.UpdateUserPasswordRequest{
				ID:           entity.NewID().String(),
				PasswordHash: "rotated-hash",
			})
			assert.True(ginkgo.GinkgoT(), errdefs.IsNotFound(actualErr), "expected NotFound, got %v", actualErr)
		})
	})

	ginkgo.Describe("UpdateStatus", func() {
		ginkgo.It("transitions an existing user to a new status", func() {
			created := createUser()
			assert.Equal(ginkgo.GinkgoT(), identity.UserStatusInactive, created.Status())

			actual, actualErr := userWriteDatastoreAdapter.UpdateStatus(ctx, identity.UpdateUserStatusRequest{
				ID:     created.ID().String(),
				Status: identity.UserStatusActive,
			})
			assert.NoError(ginkgo.GinkgoT(), actualErr)
			assert.Equal(ginkgo.GinkgoT(), created.ID(), actual.ID())
			assert.Equal(ginkgo.GinkgoT(), identity.UserStatusActive, actual.Status())
		})

		ginkgo.It("returns not found when the user does not exist", func() {
			_, actualErr := userWriteDatastoreAdapter.UpdateStatus(ctx, identity.UpdateUserStatusRequest{
				ID:     entity.NewID().String(),
				Status: identity.UserStatusActive,
			})
			assert.True(ginkgo.GinkgoT(), errdefs.IsNotFound(actualErr), "expected NotFound, got %v", actualErr)
		})
	})

	ginkgo.Describe("SoftDelete", func() {
		ginkgo.It("marks an existing user as deleted without removing it", func() {
			created := createUser()

			assert.NoError(ginkgo.GinkgoT(), userWriteDatastoreAdapter.SoftDelete(ctx, created.ID().String()))

			readDatastore := NewUserReadRepoAdapter(testPostgresDatabaseClient, database.Config{
				OperationTimeout: 5 * time.Minute,
				RetryConfig:      recoveryutils.DefaultRetryConfig,
			})
			actual, actualErr := readDatastore.FetchById(ctx, created.ID().String())
			assert.NoError(ginkgo.GinkgoT(), actualErr)
			assert.True(ginkgo.GinkgoT(), actual.IsDeleted())
		})

		ginkgo.It("returns not found when the user does not exist", func() {
			actualErr := userWriteDatastoreAdapter.SoftDelete(ctx, entity.NewID().String())
			assert.True(ginkgo.GinkgoT(), errdefs.IsNotFound(actualErr), "expected NotFound, got %v", actualErr)
		})
	})

	ginkgo.Describe("Delete", func() {
		ginkgo.It("permanently removes an existing user", func() {
			created := createUser()

			assert.NoError(ginkgo.GinkgoT(), userWriteDatastoreAdapter.Delete(ctx, created.ID().String()))

			readDatastore := NewUserReadRepoAdapter(testPostgresDatabaseClient, database.Config{
				OperationTimeout: 5 * time.Minute,
				RetryConfig:      recoveryutils.DefaultRetryConfig,
			})
			_, actualErr := readDatastore.FetchById(ctx, created.ID().String())
			assert.True(ginkgo.GinkgoT(), errdefs.IsNotFound(actualErr), "expected NotFound, got %v", actualErr)
		})

		ginkgo.It("returns not found when the user does not exist", func() {
			actualErr := userWriteDatastoreAdapter.Delete(ctx, entity.NewID().String())
			assert.True(ginkgo.GinkgoT(), errdefs.IsNotFound(actualErr), "expected NotFound, got %v", actualErr)
		})
	})
})

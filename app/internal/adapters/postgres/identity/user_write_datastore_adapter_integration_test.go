//go:build integration
// +build integration

package identitydatastore

import (
	"context"
	"testing"
	"time"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"github.com/sanctumlabs/curtz/app/internal/domain/identity"
	mockidentity "github.com/sanctumlabs/curtz/app/internal/domain/identity/mocks"
	"github.com/sanctumlabs/curtz/app/pkg/infra/database"
	recoveryutils "github.com/sanctumlabs/curtz/app/pkg/utils/recover"
	"github.com/sanctumlabs/curtz/app/test"
	"github.com/stretchr/testify/assert"
)

func TestUserWriteDatastoreAdapterIntegrationTest(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "User Write Datastore Adapter Integration Test Suite")
}

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
})

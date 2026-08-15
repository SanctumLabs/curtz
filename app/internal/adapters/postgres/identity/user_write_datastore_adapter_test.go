package identitydatastore

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	mockpostgresrepo "github.com/sanctumlabs/curtz/app/internal/adapters/postgres/mocks"
	postgresql "github.com/sanctumlabs/curtz/app/internal/adapters/postgres/sql"
	mockpostgresql "github.com/sanctumlabs/curtz/app/internal/adapters/postgres/sql/mocks"
	"github.com/sanctumlabs/curtz/app/internal/core/entity"
	"github.com/sanctumlabs/curtz/app/internal/domain/identity"
	mockidentity "github.com/sanctumlabs/curtz/app/internal/domain/identity/mocks"
	"github.com/sanctumlabs/curtz/app/pkg/infra/database"
	mockdatabase "github.com/sanctumlabs/curtz/app/pkg/infra/database/mocks"
	recoveryutils "github.com/sanctumlabs/curtz/app/pkg/utils/recover"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestUserWriteDatastoreAdapterUnitTest(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "User Write Datastore Adapter Unit Test Suite")
}

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
})

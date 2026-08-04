package identitydatastore

import (
	"context"
	"testing"
	"time"

	mockpostgresrepo "github.com/sanctumlabs/curtz/app/internal/adapters/postgres/mocks"
	postgresql "github.com/sanctumlabs/curtz/app/internal/adapters/postgres/sql"
	mockpostgresql "github.com/sanctumlabs/curtz/app/internal/adapters/postgres/sql/mocks"
	"github.com/sanctumlabs/curtz/app/internal/core/entity"
	"github.com/sanctumlabs/curtz/app/internal/domain/identity"
	mockidentity "github.com/sanctumlabs/curtz/app/internal/domain/identity/mocks"
	"github.com/sanctumlabs/curtz/app/pkg/infra/database"
	mockdatabase "github.com/sanctumlabs/curtz/app/pkg/infra/database/mocks"
	recoveryutils "github.com/sanctumlabs/curtz/app/pkg/utils/recover"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type UserReadDatastoreAdapterTestSuite struct {
	suite.Suite
	mockCtrl                 *gomock.Controller
	mockDbClient             *mockdatabase.MockPostgresDatabaseClient
	mockUserReadQuerier      *mockpostgresrepo.MockUserReadQuerier
	mockUserReadDatastore    *mockidentity.MockUserReadRepository
	userReadDatastoreAdapter *userReadDatastoreAdapter
	config                   database.Config
}

func (suite *UserReadDatastoreAdapterTestSuite) SetupTest() {
	config := database.Config{
		OperationTimeout: 30 * time.Second,
		RetryConfig:      recoveryutils.DefaultRetryConfig,
	}
	mockCtrl := gomock.NewController(suite.T())
	suite.mockCtrl = mockCtrl
	suite.mockDbClient = mockdatabase.NewMockPostgresDatabaseClient(mockCtrl)
	suite.mockUserReadQuerier = mockpostgresrepo.NewMockUserReadQuerier(mockCtrl)
	suite.mockUserReadDatastore = mockidentity.NewMockUserReadRepository(mockCtrl)
	suite.userReadDatastoreAdapter = &userReadDatastoreAdapter{
		logPrefix: "UserReadRepoAdapter",
		dbClient:  suite.mockDbClient,
		config:    config,
	}
	suite.config = config

	injectMockUserReadTx(suite.userReadDatastoreAdapter, suite.mockUserReadQuerier)
}

func TestUserReadDatastoreAdapterTestSuite(t *testing.T) {
	suite.Run(t, new(UserReadDatastoreAdapterTestSuite))
}

func (suite *UserReadDatastoreAdapterTestSuite) AfterTest(_, _ string) {
	suite.mockCtrl.Finish()
}

// TestFetchById_Success tests the FetchById method of the UserReadRepositoryAdapter
func (suite *UserReadDatastoreAdapterTestSuite) TestFetchById_Success() {
	bcgCtx := context.Background()
	ctx, cancel := context.WithTimeout(bcgCtx, suite.config.OperationTimeout)
	defer cancel()

	userId := entity.NewID()

	mockUser, mockUserErr := mockidentity.MockUser(
		mockidentity.WithId(userId),
	)
	suite.NoError(mockUserErr)

	mockUserRecord := mockpostgresql.MockUser(mockpostgresql.WithUser(*mockUser))
	mockUserStatus := mockpostgresql.MockUserStatus(identity.UserStatusActive)

	existingUser := postgresql.QueryUserByIdRow{
		User:       mockUserRecord,
		UserStatus: mockUserStatus,
	}

	suite.mockUserReadQuerier.
		EXPECT().
		QueryUserById(gomock.Any(), gomock.Any()).
		Return(existingUser, nil).
		Times(1)

	actual, actualErr := suite.userReadDatastoreAdapter.FetchById(ctx, userId.String())
	suite.Nil(actualErr)
	suite.Equal(mockUser.ID(), actual.ID())
	suite.Equal(mockUser.Username(), actual.Username())
	suite.Equal(mockUser.Email(), actual.Email())
	suite.Equal(mockUser.FirstName(), actual.FirstName())
	suite.Equal(mockUser.LastName(), actual.LastName())
	suite.Equal(mockUser.CreatedAt(), actual.CreatedAt())
	suite.Equal(mockUser.UpdatedAt(), actual.UpdatedAt())
}

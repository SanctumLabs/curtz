package identitydatastore

import (
	"context"
	"testing"
	"time"

	mockpostgresrepo "github.com/sanctumlabs/curtz/app/internal/adapters/postgres/mocks"
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

type UserWriteDatastoreAdapterTestSuite struct {
	suite.Suite
	mockCtrl                  *gomock.Controller
	mockDbClient              *mockdatabase.MockPostgresDatabaseClient
	mockUserWriteQuerier      *mockpostgresrepo.MockUserWriteQuerier
	userWriteDatastoreAdapter *userWriteDatastoreAdapter
	config                    database.Config
}

func (suite *UserWriteDatastoreAdapterTestSuite) SetupTest() {
	config := database.Config{
		OperationTimeout: 30 * time.Second,
		RetryConfig:      recoveryutils.DefaultRetryConfig,
	}
	mockCtrl := gomock.NewController(suite.T())
	suite.mockCtrl = mockCtrl
	suite.mockDbClient = mockdatabase.NewMockPostgresDatabaseClient(mockCtrl)
	suite.mockUserWriteQuerier = mockpostgresrepo.NewMockUserWriteQuerier(mockCtrl)
	suite.userWriteDatastoreAdapter = &userWriteDatastoreAdapter{
		logPrefix: "UserWriteRepoAdapter",
		dbClient:  suite.mockDbClient,
		config:    config,
	}
	suite.config = config

	injectMockUserWriteTx(suite.userWriteDatastoreAdapter, suite.mockUserWriteQuerier)
}

func TestUserWriteDatastoreAdapterTestSuite(t *testing.T) {
	suite.Run(t, new(UserWriteDatastoreAdapterTestSuite))
}

func (suite *UserWriteDatastoreAdapterTestSuite) AfterTest(_, _ string) {
	suite.mockCtrl.Finish()
}

// TestCreateUser_Success tests the CreateUser method of the UserWriteDatastoreAdapter
func (suite *UserWriteDatastoreAdapterTestSuite) TestCreateUser_Success() {
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

	suite.mockUserWriteQuerier.
		EXPECT().
		QueryUserStatusByName(gomock.Any(), gomock.Any()).
		Return(mockUserStatus, nil).
		Times(1)

	suite.mockUserWriteQuerier.
		EXPECT().
		QueryCreateUser(gomock.Any(), gomock.Any()).
		Return(mockUserRecord, nil).
		Times(1)

	actual, actualErr := suite.userWriteDatastoreAdapter.Create(ctx, *mockUser)
	suite.Nil(actualErr)
	suite.Equal(mockUser.ID(), actual.ID())
	suite.Equal(mockUser.Username(), actual.Username())
	suite.Equal(mockUser.Email(), actual.Email())
	suite.Equal(mockUser.FirstName(), actual.FirstName())
	suite.Equal(mockUser.LastName(), actual.LastName())
	suite.Equal(mockUser.CreatedAt(), actual.CreatedAt())
	suite.Equal(mockUser.UpdatedAt(), actual.UpdatedAt())
}

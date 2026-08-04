package identitydatastore

import (
	"context"
	"testing"
	"time"

	"github.com/sanctumlabs/curtz/app/internal/domain/identity"
	mockidentity "github.com/sanctumlabs/curtz/app/internal/domain/identity/mocks"
	"github.com/sanctumlabs/curtz/app/pkg/infra/database"
	recoveryutils "github.com/sanctumlabs/curtz/app/pkg/utils/recover"
	"github.com/sanctumlabs/curtz/app/test"
	"github.com/stretchr/testify/suite"
)

type UserWriteDatastoreAdapterIntegrationTestSuite struct {
	suite.Suite
	userWriteRepoAdapter       identity.UserWriteDatastore
	config                     database.Config
	testPostgresDatabaseClient database.PostgresDatabaseClient
}

func (suite *UserWriteDatastoreAdapterIntegrationTestSuite) SetupTest() {
	ctx := context.Background()
	testPostgresDatabaseClient := test.TestPostgresDatabaseClient(suite.T(), ctx)

	config := database.Config{
		OperationTimeout: 5 * time.Minute,
		RetryConfig:      recoveryutils.DefaultRetryConfig,
	}
	suite.testPostgresDatabaseClient = testPostgresDatabaseClient
	urlWriteRepoAdapter := NewUserWriteDatastoreAdapter(
		testPostgresDatabaseClient, config,
	)
	suite.userWriteRepoAdapter = urlWriteRepoAdapter
	suite.config = config
}

func TestUserWriteDatastoreAdapterIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(UserWriteDatastoreAdapterIntegrationTestSuite))
}

func (suite *UserWriteDatastoreAdapterIntegrationTestSuite) AfterTest(_, _ string) {
	suite.testPostgresDatabaseClient.Close()
}

// TestCreate_CreatesNewUrlSuccessfully tests the Create method of the UrlWriteRepositoryAdapter
func (suite *UserWriteDatastoreAdapterIntegrationTestSuite) TestCreate_CreatesNewUrlSuccessfully() {
	bcgCtx := context.Background()
	ctx, cancel := context.WithTimeout(bcgCtx, suite.config.OperationTimeout)
	defer cancel()

	mockUser, mockUserErr := mockidentity.MockUser()
	suite.NoError(mockUserErr)

	// Require stops the test immediately on failure, preventing a nil
	// dereference on *mockUser in the Create call below.
	suite.Require().NoError(mockUserErr)

	actual, actualErr := suite.userWriteRepoAdapter.Create(ctx, *mockUser)
	suite.Require().NoError(actualErr)
	suite.Require().NotNil(actual)
}

package urlrepo

import (
	"context"
	"testing"
	"time"

	domainurl "github.com/sanctumlabs/curtz/app/internal/domain/url"
	urlmock "github.com/sanctumlabs/curtz/app/internal/domain/url/mocks"
	"github.com/sanctumlabs/curtz/app/pkg/infra/database"
	recoveryutils "github.com/sanctumlabs/curtz/app/pkg/utils/recover"
	"github.com/sanctumlabs/curtz/app/test"
	"github.com/stretchr/testify/suite"
)

type UrlWriteRepoAdapterIntegrationTestSuite struct {
	suite.Suite
	urlWriteRepoAdapter        domainurl.UrlWriteRepository
	config                     database.Config
	testPostgresDatabaseClient database.PostgresDatabaseClient
}

func (suite *UrlWriteRepoAdapterIntegrationTestSuite) SetupTest() {
	ctx := context.Background()
	testPostgresDatabaseClient := test.TestPostgresDatabaseClient(suite.T(), ctx)

	config := database.Config{
		OperationTimeout: 5 * time.Minute,
		RetryConfig:      recoveryutils.DefaultRetryConfig,
	}
	suite.testPostgresDatabaseClient = testPostgresDatabaseClient
	urlWriteRepoAdapter := NewUrlWriteRepoAdapter(
		testPostgresDatabaseClient, config,
	)
	suite.urlWriteRepoAdapter = urlWriteRepoAdapter
	suite.config = config
}

func TestUrlWriteRepoAdapterIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(UrlWriteRepoAdapterIntegrationTestSuite))
}

func (suite *UrlWriteRepoAdapterIntegrationTestSuite) AfterTest(_, _ string) {
	suite.testPostgresDatabaseClient.Close()
}

// TestCreate_CreatesNewUrlSuccessfully tests the Create method of the UrlWriteRepositoryAdapter
func (suite *UrlWriteRepoAdapterIntegrationTestSuite) TestCreate_CreatesNewUrlSuccessfully() {
	bcgCtx := context.Background()
	ctx, cancel := context.WithTimeout(bcgCtx, suite.config.OperationTimeout)
	defer cancel()

	mockUrl, mockUrlErr := urlmock.MockUrl(
		urlmock.WithExpiresOn(time.Now().Add(time.Hour*24)),
		urlmock.WithCustomAlias("custom"),
		urlmock.WithShortCode("shortcode"),
	)
	// Require stops the test immediately on failure, preventing a nil
	// dereference on *mockUrl in the Create call below.
	suite.Require().NoError(mockUrlErr)

	actual, actualErr := suite.urlWriteRepoAdapter.Create(ctx, *mockUrl)
	suite.Require().NoError(actualErr)

	suite.Equal(mockUrl.UserId(), actual.UserId())
	suite.Equal(mockUrl.ShortCode(), actual.ShortCode())
	suite.Equal(mockUrl.CustomAlias(), actual.CustomAlias())
	suite.Equal(mockUrl.OriginalURL(), actual.OriginalURL())
	suite.Equal(mockUrl.ExpiresOn(), actual.ExpiresOn())
	suite.Equal(mockUrl.Status(), actual.Status())
	suite.Equal(mockUrl.OgTitle(), actual.OgTitle())
	suite.Equal(mockUrl.OgDescription(), actual.OgDescription())
	suite.Equal(mockUrl.OgImageUrl(), actual.OgImageUrl())
}

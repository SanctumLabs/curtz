//go:build integration

package urldatastore

import (
	"context"
	"testing"
	"time"

	identitydatastore "github.com/sanctumlabs/curtz/app/internal/adapters/postgres/identity"
	"github.com/sanctumlabs/curtz/app/internal/core/entity"
	"github.com/sanctumlabs/curtz/app/internal/domain/identity"
	mockidentity "github.com/sanctumlabs/curtz/app/internal/domain/identity/mocks"
	domainurl "github.com/sanctumlabs/curtz/app/internal/domain/url"
	urlmock "github.com/sanctumlabs/curtz/app/internal/domain/url/mocks"
	"github.com/sanctumlabs/curtz/app/internal/pkg/common"
	"github.com/sanctumlabs/curtz/app/pkg/errdefs"
	"github.com/sanctumlabs/curtz/app/pkg/infra/database"
	recoveryutils "github.com/sanctumlabs/curtz/app/pkg/utils/recover"
	"github.com/sanctumlabs/curtz/app/test"
	"github.com/stretchr/testify/suite"
)

type UrlReadDatastoreAdapterIntegrationTestSuite struct {
	suite.Suite
	urlReadDatastoreAdapter    domainurl.UrlReadDatastore
	urlWriteDatastoreAdapter   domainurl.UrlWriteDatastore
	userWriteDatastore         identity.UserWriteDatastore
	config                     database.Config
	testPostgresDatabaseClient database.PostgresDatabaseClient
}

func (suite *UrlReadDatastoreAdapterIntegrationTestSuite) SetupTest() {
	ctx := context.Background()
	testPostgresDatabaseClient := test.TestPostgresDatabaseClientHelper(suite.T(), ctx)

	config := database.Config{
		OperationTimeout: 5 * time.Minute,
		RetryConfig:      recoveryutils.DefaultRetryConfig,
	}
	suite.testPostgresDatabaseClient = testPostgresDatabaseClient
	suite.userWriteDatastore = identitydatastore.NewUserWriteDatastoreAdapter(testPostgresDatabaseClient, config)
	suite.urlWriteDatastoreAdapter = NewUrlWriteDatastoreAdapter(testPostgresDatabaseClient, config)
	suite.urlReadDatastoreAdapter = NewUrlReadDatastoreAdapter(testPostgresDatabaseClient, config)
	suite.config = config
}

func TestUrlReadDatastoreAdapterIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(UrlReadDatastoreAdapterIntegrationTestSuite))
}

func (suite *UrlReadDatastoreAdapterIntegrationTestSuite) AfterTest(_, _ string) {
	suite.testPostgresDatabaseClient.Close()
}

// seedUrl persists a user and then a URL owned by that user, returning the stored URL.
// urls.user_id references users(id), so the owner has to exist before the URL is written.
func (suite *UrlReadDatastoreAdapterIntegrationTestSuite) seedUrl(ctx context.Context, options ...urlmock.MockUrlOption) domainurl.URL {
	mockUser, mockUserErr := mockidentity.MockUser()
	suite.Require().NoError(mockUserErr)

	_, saveUserErr := suite.userWriteDatastore.Save(ctx, *mockUser)
	suite.Require().NoError(saveUserErr)

	defaults := []urlmock.MockUrlOption{
		urlmock.WithUserId(mockUser.ID().String()),
		urlmock.WithExpiresOn(time.Now().Add(24 * time.Hour)),
	}

	mockUrl, mockUrlErr := urlmock.MockUrl(append(defaults, options...)...)
	suite.Require().NoError(mockUrlErr)

	savedUrl, savedUrlErr := suite.urlWriteDatastoreAdapter.Save(ctx, *mockUrl)
	suite.Require().NoError(savedUrlErr)

	return savedUrl
}

func (suite *UrlReadDatastoreAdapterIntegrationTestSuite) newContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), suite.config.OperationTimeout)
}

func (suite *UrlReadDatastoreAdapterIntegrationTestSuite) TestFetchById_ReturnsStoredUrl() {
	ctx, cancel := suite.newContext()
	defer cancel()

	saved := suite.seedUrl(ctx)

	actual, actualErr := suite.urlReadDatastoreAdapter.FetchById(ctx, saved.ID().String())
	suite.Require().NoError(actualErr)
	suite.Equal(saved.ID(), actual.ID())
	suite.Equal(saved.UserId(), actual.UserId())
	suite.Equal(saved.ShortCode(), actual.ShortCode())
	suite.Equal(saved.CustomAlias(), actual.CustomAlias())
	suite.Equal(saved.OriginalURL(), actual.OriginalURL())
	suite.Equal(domainurl.URLStatusActive, actual.Status())
	suite.WithinDuration(saved.ExpiresOn(), actual.ExpiresOn(), time.Microsecond)
}

func (suite *UrlReadDatastoreAdapterIntegrationTestSuite) TestFetchById_NotFoundForUnknownId() {
	ctx, cancel := suite.newContext()
	defer cancel()

	_, actualErr := suite.urlReadDatastoreAdapter.FetchById(ctx, entity.NewID().String())
	suite.Require().Error(actualErr)
	suite.True(errdefs.IsNotFound(actualErr), "expected a NotFound error, got %v", actualErr)
}

func (suite *UrlReadDatastoreAdapterIntegrationTestSuite) TestFetchByShortCode_ReturnsStoredUrl() {
	ctx, cancel := suite.newContext()
	defer cancel()

	saved := suite.seedUrl(ctx, urlmock.WithShortCode("abc1234"))

	actual, actualErr := suite.urlReadDatastoreAdapter.FetchByShortCode(ctx, "abc1234")
	suite.Require().NoError(actualErr)
	suite.Equal(saved.ID(), actual.ID())
	suite.Equal("abc1234", actual.ShortCode().Value())
}

func (suite *UrlReadDatastoreAdapterIntegrationTestSuite) TestFetchByShortCode_NotFound() {
	ctx, cancel := suite.newContext()
	defer cancel()

	_, actualErr := suite.urlReadDatastoreAdapter.FetchByShortCode(ctx, "nosuch")
	suite.Require().Error(actualErr)
	suite.True(errdefs.IsNotFound(actualErr), "expected a NotFound error, got %v", actualErr)
}

func (suite *UrlReadDatastoreAdapterIntegrationTestSuite) TestFetchByCustomAlias_ReturnsStoredUrl() {
	ctx, cancel := suite.newContext()
	defer cancel()

	saved := suite.seedUrl(ctx, urlmock.WithCustomAlias("my-alias"))

	actual, actualErr := suite.urlReadDatastoreAdapter.FetchByCustomAlias(ctx, "my-alias")
	suite.Require().NoError(actualErr)
	suite.Equal(saved.ID(), actual.ID())
	suite.Equal("my-alias", actual.CustomAlias().Value())
}

func (suite *UrlReadDatastoreAdapterIntegrationTestSuite) TestFetchByCustomAlias_NotFound() {
	ctx, cancel := suite.newContext()
	defer cancel()

	_, actualErr := suite.urlReadDatastoreAdapter.FetchByCustomAlias(ctx, "ghost")
	suite.Require().Error(actualErr)
	suite.True(errdefs.IsNotFound(actualErr), "expected a NotFound error, got %v", actualErr)
}

func (suite *UrlReadDatastoreAdapterIntegrationTestSuite) TestFetchByOriginalUrl_ReturnsStoredUrl() {
	ctx, cancel := suite.newContext()
	defer cancel()

	saved := suite.seedUrl(ctx, urlmock.WithOriginalUrl("https://example.com/original"))

	actual, actualErr := suite.urlReadDatastoreAdapter.FetchByOriginalUrl(ctx, "https://example.com/original")
	suite.Require().NoError(actualErr)
	suite.Equal(saved.ID(), actual.ID())
	suite.Equal("https://example.com/original", actual.OriginalURL().Value())
}

func (suite *UrlReadDatastoreAdapterIntegrationTestSuite) TestFetchByOriginalUrl_NotFound() {
	ctx, cancel := suite.newContext()
	defer cancel()

	_, actualErr := suite.urlReadDatastoreAdapter.FetchByOriginalUrl(ctx, "https://example.com/never-stored")
	suite.Require().Error(actualErr)
	suite.True(errdefs.IsNotFound(actualErr), "expected a NotFound error, got %v", actualErr)
}

func (suite *UrlReadDatastoreAdapterIntegrationTestSuite) TestFetchAll_ReturnsStoredUrlsWithTotal() {
	ctx, cancel := suite.newContext()
	defer cancel()

	first := suite.seedUrl(ctx)
	second := suite.seedUrl(ctx)

	actual, actualErr := suite.urlReadDatastoreAdapter.FetchAll(ctx, common.NewRequestParams())
	suite.Require().NoError(actualErr)
	suite.Equal(2, actual.Total)
	suite.Equal(2, actual.Size)
	suite.Equal(1, actual.Page)

	ids := make([]string, 0, len(actual.Records))
	for _, record := range actual.Records {
		ids = append(ids, record.ID().String())
	}
	suite.Contains(ids, first.ID().String())
	suite.Contains(ids, second.ID().String())
}

// TestFetchAll_PaginatesWithTotal checks that the total reflects every matching row while the page
// only carries the requested window.
func (suite *UrlReadDatastoreAdapterIntegrationTestSuite) TestFetchAll_PaginatesWithTotal() {
	ctx, cancel := suite.newContext()
	defer cancel()

	for range 3 {
		suite.seedUrl(ctx)
	}

	actual, actualErr := suite.urlReadDatastoreAdapter.FetchAll(ctx,
		common.NewRequestParams(common.WithRequestLimit(2), common.WithOffset(2)))
	suite.Require().NoError(actualErr)
	suite.Equal(3, actual.Total)
	suite.Equal(1, actual.Size)
	suite.Equal(2, actual.Page)
}

func (suite *UrlReadDatastoreAdapterIntegrationTestSuite) TestFetchByUserId_ReturnsOnlyThatUsersUrls() {
	ctx, cancel := suite.newContext()
	defer cancel()

	mine := suite.seedUrl(ctx)
	// A second seeded URL belongs to a different user, so it must not appear in the result.
	suite.seedUrl(ctx)

	actual, actualErr := suite.urlReadDatastoreAdapter.FetchByUserId(ctx, mine.UserId().String())
	suite.Require().NoError(actualErr)
	suite.Equal(1, actual.Total)
	suite.Require().Len(actual.Records, 1)
	suite.Equal(mine.ID(), actual.Records[0].ID())
	suite.Equal(mine.UserId(), actual.Records[0].UserId())
}

func (suite *UrlReadDatastoreAdapterIntegrationTestSuite) TestFetchByUserId_EmptyForUnknownUser() {
	ctx, cancel := suite.newContext()
	defer cancel()

	suite.seedUrl(ctx)

	actual, actualErr := suite.urlReadDatastoreAdapter.FetchByUserId(ctx, entity.NewID().String())
	suite.Require().NoError(actualErr)
	suite.Equal(0, actual.Total)
	suite.Empty(actual.Records)
}

func (suite *UrlReadDatastoreAdapterIntegrationTestSuite) TestFetchByStatus_ReturnsActiveUrls() {
	ctx, cancel := suite.newContext()
	defer cancel()

	saved := suite.seedUrl(ctx)

	actual, actualErr := suite.urlReadDatastoreAdapter.FetchByStatus(ctx, domainurl.URLStatusActive)
	suite.Require().NoError(actualErr)
	suite.Equal(1, actual.Total)
	suite.Require().Len(actual.Records, 1)
	suite.Equal(saved.ID(), actual.Records[0].ID())
	suite.Equal(domainurl.URLStatusActive, actual.Records[0].Status())
}

// TestFetchByStatus_EmptyWhenNoneMatch guards the status filter: an ACTIVE row must not leak into a
// query for another status.
func (suite *UrlReadDatastoreAdapterIntegrationTestSuite) TestFetchByStatus_EmptyWhenNoneMatch() {
	ctx, cancel := suite.newContext()
	defer cancel()

	suite.seedUrl(ctx)

	actual, actualErr := suite.urlReadDatastoreAdapter.FetchByStatus(ctx, domainurl.URLStatusSuspended)
	suite.Require().NoError(actualErr)
	suite.Equal(0, actual.Total)
	suite.Empty(actual.Records)
}

// TestFetchExpiredActive_ReturnsOnlyAlreadyExpiredUrls is the query the expiry job depends on: an
// ACTIVE URL past its expiry is returned, one still in the future is not.
func (suite *UrlReadDatastoreAdapterIntegrationTestSuite) TestFetchExpiredActive_ReturnsOnlyAlreadyExpiredUrls() {
	ctx, cancel := suite.newContext()
	defer cancel()

	expired := suite.seedUrl(ctx, urlmock.WithExpiresOn(time.Now().Add(-1*time.Hour)))
	suite.seedUrl(ctx, urlmock.WithExpiresOn(time.Now().Add(24*time.Hour)))

	actual, actualErr := suite.urlReadDatastoreAdapter.FetchExpiredActive(ctx, time.Now(), 10)
	suite.Require().NoError(actualErr)
	suite.Require().Len(actual, 1)
	suite.Equal(expired.ID(), actual[0].ID())
	suite.Equal(domainurl.URLStatusActive, actual[0].Status())
}

// TestFetchExpiredActive_RespectsLimit matters because the expiry job drains the backlog in batches.
func (suite *UrlReadDatastoreAdapterIntegrationTestSuite) TestFetchExpiredActive_RespectsLimit() {
	ctx, cancel := suite.newContext()
	defer cancel()

	for i := range 3 {
		suite.seedUrl(ctx, urlmock.WithExpiresOn(time.Now().Add(-time.Duration(i+1)*time.Hour)))
	}

	actual, actualErr := suite.urlReadDatastoreAdapter.FetchExpiredActive(ctx, time.Now(), 2)
	suite.Require().NoError(actualErr)
	suite.Len(actual, 2)
}

func (suite *UrlReadDatastoreAdapterIntegrationTestSuite) TestFetchExpiredActive_EmptyWhenNothingExpired() {
	ctx, cancel := suite.newContext()
	defer cancel()

	suite.seedUrl(ctx, urlmock.WithExpiresOn(time.Now().Add(24*time.Hour)))

	actual, actualErr := suite.urlReadDatastoreAdapter.FetchExpiredActive(ctx, time.Now(), 10)
	suite.Require().NoError(actualErr)
	suite.Empty(actual)
}

// TestFetchByOriginalUrl_FindsRespelledTarget closes the loop on normalization: the row is stored in
// canonical form, and a lookup spelled differently must still find it.
func (suite *UrlReadDatastoreAdapterIntegrationTestSuite) TestFetchByOriginalUrl_FindsRespelledTarget() {
	ctx, cancel := suite.newContext()
	defer cancel()

	saved := suite.seedUrl(ctx, urlmock.WithOriginalUrl("https://example.com/canonical"))

	for _, respelling := range []string{
		"https://example.com/canonical",
		"https://EXAMPLE.COM/canonical",
		"https://example.com:443/canonical",
		"https://example.com/canonical/",
		"https://example.com/canonical?utm_source=news",
		"HTTPS://Example.Com:443/canonical/?utm_medium=email",
	} {
		suite.Run(respelling, func() {
			actual, actualErr := suite.urlReadDatastoreAdapter.FetchByOriginalUrl(ctx, respelling)
			suite.Require().NoError(actualErr, "%q should have found the stored target", respelling)
			suite.Equal(saved.ID(), actual.ID())
		})
	}
}

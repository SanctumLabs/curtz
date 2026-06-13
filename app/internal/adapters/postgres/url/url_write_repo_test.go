package urlrepo

import (
	"context"
	"testing"
	"time"

	mockpostgresql "github.com/sanctumlabs/curtz/app/internal/adapters/postgres/sql/mocks"
	mockurlrepo "github.com/sanctumlabs/curtz/app/internal/adapters/postgres/url/mocks"
	"github.com/sanctumlabs/curtz/app/internal/domain/url"
	urlmock "github.com/sanctumlabs/curtz/app/internal/domain/url/mocks"
	"github.com/sanctumlabs/curtz/app/pkg/infra/database"
	mockdatabase "github.com/sanctumlabs/curtz/app/pkg/infra/database/mocks"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type UrlWriteRepoAdapterTestSuite struct {
	suite.Suite
	mockCtrl            *gomock.Controller
	mockDbClient        *mockdatabase.MockPostgresDatabaseClient
	mockUrlWriteQuerier *mockurlrepo.MockUrlWriteQuerier
	urlWriteRepoAdapter *urlWriteRepositoryAdapter
	config              database.Config
}

func (suite *UrlWriteRepoAdapterTestSuite) SetupTest() {
	config := database.Config{
		OperationTimeout: 15 * time.Second,
	}
	mockCtrl := gomock.NewController(suite.T())
	suite.mockCtrl = mockCtrl
	suite.mockDbClient = mockdatabase.NewMockPostgresDatabaseClient(mockCtrl)
	suite.mockUrlWriteQuerier = mockurlrepo.NewMockUrlWriteQuerier(mockCtrl)
	suite.urlWriteRepoAdapter = &urlWriteRepositoryAdapter{
		logPrefix: "UrlWriteRepoAdapter",
		dbClient:  suite.mockDbClient,
		config:    config,
	}
}

func TestUrlWriteRepoAdapterTestSuite(t *testing.T) {
	suite.Run(t, new(UrlWriteRepoAdapterTestSuite))
}

func (suite *UrlWriteRepoAdapterTestSuite) AfterTest(_, _ string) {
	suite.mockCtrl.Finish()
}

// TestCreate_CreatesNewUrlSuccessfully tests the Create method of the UrlWriteRepositoryAdapter
func (suite *UrlWriteRepoAdapterTestSuite) TestCreate_CreatesNewUrlSuccessfully() {
	bcgCtx := context.Background()
	ctx, _ := context.WithTimeout(bcgCtx, suite.config.OperationTimeout)
	mockUrl, mockUrlErr := urlmock.MockUrl(
		urlmock.WithExpiresOn(time.Now().Add(time.Hour*24)),
		urlmock.WithCustomAlias("custom"),
		urlmock.WithShortCode("shortcode"),
	)
	suite.NoError(mockUrlErr)

	mockCreatedUrl := mockpostgresql.MockUrl(
		mockpostgresql.WithUserId(mockUrl.UserId()),
		mockpostgresql.WithOriginalUrl(mockUrl.OriginalURL().Value()),
	)

	suite.mockUrlWriteQuerier.
		EXPECT().
		QueryUrlStatusByName(ctx, gomock.Any()).
		Return(mockpostgresql.MockUrlStatus(url.URLStatusActive), nil).
		Times(1)

	suite.mockUrlWriteQuerier.
		EXPECT().
		QueryCreateUrl(ctx, gomock.Any()).
		Return(mockCreatedUrl, nil).
		Times(1)

	suite.mockUrlWriteQuerier.EXPECT().QueryCreateKeyword(ctx, gomock.Any())

	injectMockUrlWriteTx(suite.urlWriteRepoAdapter, suite.mockUrlWriteQuerier)

	actual, actualErr := suite.urlWriteRepoAdapter.Create(ctx, *mockUrl)
	suite.Nil(actualErr)
	// suite.Equal(mockUrl.UserId(), actual.UserId())
	suite.Equal(mockUrl.OriginalURL(), actual.OriginalURL())
}

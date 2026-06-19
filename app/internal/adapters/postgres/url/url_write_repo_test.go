package urlrepo

import (
	"context"
	"testing"
	"time"

	mockurlrepo "github.com/sanctumlabs/curtz/app/internal/adapters/postgres/mocks"
	mockpostgresql "github.com/sanctumlabs/curtz/app/internal/adapters/postgres/sql/mocks"
	"github.com/sanctumlabs/curtz/app/internal/domain/url"
	urlmock "github.com/sanctumlabs/curtz/app/internal/domain/url/mocks"
	"github.com/sanctumlabs/curtz/app/pkg/infra/database"
	mockdatabase "github.com/sanctumlabs/curtz/app/pkg/infra/database/mocks"
	recoveryutils "github.com/sanctumlabs/curtz/app/pkg/utils/recover"
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
		OperationTimeout: 30 * time.Second,
		RetryConfig:      recoveryutils.DefaultRetryConfig,
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
	suite.config = config

	injectMockUrlWriteTx(suite.urlWriteRepoAdapter, suite.mockUrlWriteQuerier)
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
	ctx, cancel := context.WithTimeout(bcgCtx, suite.config.OperationTimeout)
	defer cancel()

	mockUrl, mockUrlErr := urlmock.MockUrl(
		urlmock.WithExpiresOn(time.Now().Add(time.Hour*24)),
		urlmock.WithCustomAlias("custom"),
		urlmock.WithShortCode("shortcode"),
	)
	suite.NoError(mockUrlErr)

	mockCreatedUrl := mockpostgresql.MockUrl(
		mockpostgresql.WithUrl(*mockUrl),
	)

	suite.mockUrlWriteQuerier.
		EXPECT().
		QueryUrlStatusByName(gomock.Any(), gomock.Any()).
		Return(mockpostgresql.MockUrlStatus(url.URLStatusActive), nil).
		Times(1)

	suite.mockUrlWriteQuerier.
		EXPECT().
		QueryCreateUrl(gomock.Any(), gomock.Any()).
		Return(mockCreatedUrl, nil).
		Times(1)

	suite.mockUrlWriteQuerier.
		EXPECT().
		QueryCreateKeyword(gomock.Any(), gomock.Any()).
		AnyTimes()

	actual, actualErr := suite.urlWriteRepoAdapter.Create(ctx, *mockUrl)
	suite.Nil(actualErr)
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

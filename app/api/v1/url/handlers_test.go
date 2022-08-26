package url

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/sanctumlabs/curtz/app/server/router"
	"github.com/sanctumlabs/curtz/app/test/mocks"
)

var baseURI = "/api/v1/curtz"

func createUrlRouter(t *testing.T) (*urlRouter, mocks.MockUrlService) {
	mockCtrl := gomock.NewController(t)
	mockUrlSvc := mocks.NewMockUrlService(mockCtrl)

	urlRouter := &urlRouter{
		svc:     mockUrlSvc,
		baseUri: baseURI,
		routes:  []router.Route{},
	}

	routes := []router.Route{
		router.NewPostRoute(fmt.Sprintf("%s/urls", baseURI), urlRouter.createShortUrl),
		router.NewGetRoute(fmt.Sprintf("%s/urls", baseURI), urlRouter.getAllUrls),
		router.NewGetRoute(fmt.Sprintf("%s/urls/:id", baseURI), urlRouter.getUrlById),
	}

	urlRouter.routes = append(urlRouter.routes, routes...)

	return urlRouter, *mockUrlSvc
}

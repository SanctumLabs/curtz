package auth

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/sanctumlabs/curtz/app/server/router"
	"github.com/sanctumlabs/curtz/app/test/mocks"
)

var baseUri = "/api/v1/curtz"

func setupMocks(t *testing.T) (*mocks.MockUserService, *mocks.MockAuthService) {
	mockCtrl := gomock.NewController(t)
	mockUserSvc := mocks.NewMockUserService(mockCtrl)
	mockAuthSvc := mocks.NewMockAuthService(mockCtrl)

	return mockUserSvc, mockAuthSvc
}

func setupTestAuthRouter(mockUserSvc *mocks.MockUserService, mockAuthSvc *mocks.MockAuthService) *authRouter {
	authRouter := &authRouter{
		baseUri: baseUri,
		svc:     mockUserSvc,
		authSvc: mockAuthSvc,
		routes:  []router.Route{},
	}

	routes := []router.Route{
		router.NewPostRoute(fmt.Sprintf("%s/auth/register", baseUri), authRouter.register),
		router.NewPostRoute(fmt.Sprintf("%s/auth/login", baseUri), authRouter.login),
	}

	authRouter.routes = append(authRouter.routes, routes...)
	return authRouter
}

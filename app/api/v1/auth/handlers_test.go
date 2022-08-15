package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	gin "github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/sanctumlabs/curtz/app/internal/core/entities"
	"github.com/sanctumlabs/curtz/app/server/router"
	"github.com/sanctumlabs/curtz/app/test/mocks"
	"github.com/sanctumlabs/curtz/app/test/utils"
	"github.com/stretchr/testify/assert"
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

func TestFailLoginOnNonExistentPayload(t *testing.T) {
	mockUserSvc, mockAuthSvc := setupMocks(t)
	authRouter := setupTestAuthRouter(mockUserSvc, mockAuthSvc)

	resp := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(resp)

	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("%s/auth/login", baseUri), nil)

	ctx.Request = req

	authRouter.login(ctx)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestFailLoginOnNonExistentUser(t *testing.T) {
	mockUserSvc, mockAuthSvc := setupMocks(t)
	authRouter := setupTestAuthRouter(mockUserSvc, mockAuthSvc)

	gin.SetMode(gin.TestMode)
	resp := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(resp)

	req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("%s/auth/login", baseUri), nil)
	ctx.Request = req

	email := "johndoe@example.com"
	password := "password"

	requestBody := loginRequestDto{
		Email:    email,
		Password: password,
	}

	mockUserSvc.
		EXPECT().
		GetUserByEmail(email).
		Return(entities.User{}, errors.New("user does not exist"))

	utils.MockRequestBody(ctx, requestBody)
	authRouter.login(ctx)

	expectedRespBody, err := json.Marshal(gin.H{
		"message": "Invalid Email or Password",
	})

	assert.NoError(t, err)

	assert.Equal(t, http.StatusUnauthorized, resp.Code)
	assert.Equal(t, expectedRespBody, resp.Body.Bytes())
}

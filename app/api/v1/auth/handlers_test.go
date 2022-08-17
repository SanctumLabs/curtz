package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	gin "github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/sanctumlabs/curtz/app/internal/core/entities"
	"github.com/sanctumlabs/curtz/app/server/router"
	"github.com/sanctumlabs/curtz/app/test/data"
	"github.com/sanctumlabs/curtz/app/test/mocks"
	"github.com/sanctumlabs/curtz/app/test/utils"
	"github.com/stretchr/testify/assert"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestAuthHandlers(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Auth Handler Suite")
}

var _ = Describe("Auth Handler", func() {
	var (
		mockCtrl    *gomock.Controller
		mockUserSvc *mocks.MockUserService
		mockAuthSvc *mocks.MockAuthService
	)
	authRouter := &authRouter{}

	BeforeEach(func() {
		gin.SetMode(gin.TestMode)

		mockCtrl = gomock.NewController(GinkgoT())
		mockUserSvc = mocks.NewMockUserService(mockCtrl)
		mockAuthSvc = mocks.NewMockAuthService(mockCtrl)
		authRouter.baseUri = baseUri
		authRouter.svc = mockUserSvc
		authRouter.authSvc = mockAuthSvc
		authRouter.routes = []router.Route{}

		routes := []router.Route{
			router.NewPostRoute(fmt.Sprintf("%s/auth/register", baseUri), authRouter.register),
			router.NewPostRoute(fmt.Sprintf("%s/auth/login", baseUri), authRouter.login),
		}

		authRouter.routes = append(authRouter.routes, routes...)
	})

	When("logging in", func() {
		httpRequest := httptest.NewRequest(http.MethodPost, fmt.Sprintf("%s/auth/login", baseUri), nil)

		Context("and payload is empty", func() {
			responseRecorder := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(responseRecorder)
			ctx.Request = httpRequest

			It("should return 400 response code", func() {
				authRouter.login(ctx)

				assert.Equal(GinkgoT(), http.StatusBadRequest, responseRecorder.Code)
			})
		})

		Context("and the user does not exist", func() {
			responseRecorder := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(responseRecorder)
			ctx.Request = httpRequest

			It("should return a 401 response code", func() {
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

				assert.NoError(GinkgoT(), err)
				assert.Equal(GinkgoT(), http.StatusUnauthorized, responseRecorder.Code)
				assert.Equal(GinkgoT(), expectedRespBody, responseRecorder.Body.Bytes())
			})
		})

		Context("and the user exists", func() {
			It("but the password supplied is invalid, should return a 401 response code", func() {
				responseRecorder := httptest.NewRecorder()
				ctx, _ := gin.CreateTestContext(responseRecorder)
				ctx.Request = httpRequest

				email := "johndoe@example.com"
				password := "wrong-password"
				correctPassword := "correct-password"

				requestBody := loginRequestDto{
					Email:    email,
					Password: password,
				}

				user, err := data.MockUser(email, correctPassword)
				assert.NoError(GinkgoT(), err)

				mockUserSvc.
					EXPECT().
					GetUserByEmail(email).
					Return(user, nil)

				utils.MockRequestBody(ctx, requestBody)

				authRouter.login(ctx)

				expectedRespBody, err := json.Marshal(gin.H{
					"message": "Invalid Email or Password",
				})

				assert.NoError(GinkgoT(), err)
				assert.Equal(GinkgoT(), http.StatusUnauthorized, responseRecorder.Code)
				assert.Equal(GinkgoT(), expectedRespBody, responseRecorder.Body.Bytes())
			})

			Context("and the password supplied is valid", func() {
				email := "johndoe@example.com"
				password := "password"
				accessToken := "header.payload.signature"
				refreshToken := "header.payload.signature"

				requestBody := loginRequestDto{
					Email:    email,
					Password: password,
				}

				user, _ := data.MockUser(email, password)

				When("there is a failure to generate", func() {

					It("an access token, should return a 400 response code", func() {
						responseRecorder := httptest.NewRecorder()
						ctx, _ := gin.CreateTestContext(responseRecorder)
						ctx.Request = httpRequest

						err := errors.New("failed to supply access token")

						mockUserSvc.
							EXPECT().
							GetUserByEmail(email).
							Return(user, nil)

						mockAuthSvc.
							EXPECT().
							GenerateToken(user.ID.String()).
							Return("", err)

						utils.MockRequestBody(ctx, requestBody)

						authRouter.login(ctx)

						expectedRespBody, err := json.Marshal(gin.H{
							"message": err.Error(),
						})

						assert.Equal(GinkgoT(), expectedRespBody, responseRecorder.Body.Bytes())
						assert.Equal(GinkgoT(), http.StatusBadRequest, responseRecorder.Code)
					})

					It("refresh token should return a 400 response code", func() {
						responseRecorder := httptest.NewRecorder()
						ctx, _ := gin.CreateTestContext(responseRecorder)
						ctx.Request = httpRequest

						err := errors.New("failed to supply refresh token")

						mockUserSvc.
							EXPECT().
							GetUserByEmail(email).
							Return(user, nil)

						mockAuthSvc.
							EXPECT().
							GenerateToken(user.ID.String()).
							Return(accessToken, nil)

						mockAuthSvc.
							EXPECT().
							GenerateRefreshToken(user.ID.String()).
							Return(refreshToken, err)

						utils.MockRequestBody(ctx, requestBody)

						authRouter.login(ctx)

						expectedRespBody, err := json.Marshal(gin.H{
							"message": err.Error(),
						})

						assert.Equal(GinkgoT(), expectedRespBody, responseRecorder.Body.Bytes())
						assert.Equal(GinkgoT(), http.StatusBadRequest, responseRecorder.Code)
					})
				})

				When("there is success generating both access and refresh tokens", func() {
					It("should return 200 response code with user information and tokens", func() {
						responseRecorder := httptest.NewRecorder()
						ctx, _ := gin.CreateTestContext(responseRecorder)
						ctx.Request = httpRequest
						mockUserSvc.
							EXPECT().
							GetUserByEmail(email).
							Return(user, nil)

						mockAuthSvc.
							EXPECT().
							GenerateToken(user.ID.String()).
							Return(accessToken, nil)

						mockAuthSvc.
							EXPECT().
							GenerateRefreshToken(user.ID.String()).
							Return(refreshToken, nil)

						utils.MockRequestBody(ctx, requestBody)

						authRouter.login(ctx)

						expectedRespBody := gin.H{
							"id":            user.ID.String(),
							"email":         user.Email.Value,
							"created_at":    user.CreatedAt.Format(time.RFC3339Nano),
							"updated_at":    user.UpdatedAt.Format(time.RFC3339Nano),
							"access_token":  accessToken,
							"refresh_token": refreshToken,
						}

						var actualResponse map[string]string
						err := json.Unmarshal([]byte(responseRecorder.Body.String()), &actualResponse)
						assert.NoError(GinkgoT(), err)

						assert.Equal(GinkgoT(), http.StatusOK, responseRecorder.Code)

						if id, ok := actualResponse["id"]; ok {
							assert.True(GinkgoT(), ok)
							assert.Equal(GinkgoT(), expectedRespBody["id"], id)
						}

						if id, ok := actualResponse["email"]; ok {
							assert.True(GinkgoT(), ok)
							assert.Equal(GinkgoT(), expectedRespBody["email"], id)
						}

						if id, ok := actualResponse["created_at"]; ok {
							assert.True(GinkgoT(), ok)
							assert.Equal(GinkgoT(), expectedRespBody["created_at"], id)
						}

						if id, ok := actualResponse["updated_at"]; ok {
							assert.True(GinkgoT(), ok)
							assert.Equal(GinkgoT(), expectedRespBody["updated_at"], id)
						}
						if id, ok := actualResponse["access_token"]; ok {
							assert.True(GinkgoT(), ok)
							assert.Equal(GinkgoT(), expectedRespBody["access_token"], id)
						}
						if id, ok := actualResponse["refresh_token"]; ok {
							assert.True(GinkgoT(), ok)
							assert.Equal(GinkgoT(), expectedRespBody["refresh_token"], id)
						}
					})
				})
			})
		})
	})
})

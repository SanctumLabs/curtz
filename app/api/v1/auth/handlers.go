package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (hdl *authRouter) register(ctx *gin.Context) {
	var request registerRequestDto
	err := ctx.BindJSON(&request)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}

	user, err := hdl.svc.CreateUser(request.Email, request.Password)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	response := userResponseDto{
		ID:        user.ID.String(),
		Email:     user.Email.Value,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	ctx.JSON(http.StatusOK, response)
}

func (hdl *authRouter) login(ctx *gin.Context) {
	var request loginRequestDto
	err := ctx.BindJSON(&request)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	user, err := hdl.svc.GetUserByEmail(request.Email)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Invalid Email or Password"})
		return
	}

	if ok, err := user.Compare(user.Password.Value, request.Password); err != nil && !ok {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Invalid Email or Password"})
		return
	}

	accessToken, err := hdl.authSvc.GenerateToken(user.ID.String())
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	refreshToken, err := hdl.authSvc.GenerateRefreshToken(user.ID.String())
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	response := loginResponseDto{
		userResponseDto: userResponseDto{
			ID:        user.ID.String(),
			Email:     user.Email.Value,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	ctx.JSON(http.StatusOK, response)
}

// oauthToken refreshes a token given a refresh token
func (hdl *authRouter) oauthToken(ctx *gin.Context) {
	grantType := ctx.Query("grant_type")
	refreshToken := ctx.Query("refresh_token")

	if grantType == "refresh_token" {
		if len(refreshToken) == 0 {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		uid, _, err := hdl.authSvc.Authenticate(refreshToken)

		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
			return
		}

		if _, err := hdl.svc.GetUserByID(uid); err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
			return
		}

		accessToken, err := hdl.authSvc.GenerateToken(uid)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		refreshToken, err := hdl.authSvc.GenerateRefreshToken(uid)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		response := oauthRefreshTokenResponseDto{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			TokenType:    "Bearer",
		}

		ctx.JSON(http.StatusOK, response)
	} else {
		ctx.AbortWithStatus(http.StatusUnauthorized)
	}
}

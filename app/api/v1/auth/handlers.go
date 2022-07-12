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

	token, err := hdl.authSvc.GenerateToken(user.ID.String())

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Invalid Email or Password"})
		return
	}

	response := userResponseDto{
		ID:          user.ID.String(),
		Email:       user.Email.Value,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		AccessToken: token,
	}

	ctx.JSON(http.StatusOK, response)
}

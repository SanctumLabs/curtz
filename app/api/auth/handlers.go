package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (hdl *authRouter) signUp(ctx *gin.Context) {
	var request signUpRequestDto
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

}

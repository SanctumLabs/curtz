package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sanctumlabs/curtz/app/internal/services/notifications"
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

	err = hdl.notificationSvc.SendNotification(user.Email.Value, "Welcome to Curtz", notifications.NotificationTypeEmail)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
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

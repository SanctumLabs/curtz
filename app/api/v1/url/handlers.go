package url

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/sanctumlabs/curtz/app/api"
	"github.com/sanctumlabs/curtz/app/internal/core/contracts"
	"github.com/sanctumlabs/curtz/app/pkg/validators"
)

// createShortUrl creates a new shortened url
func (hdl *urlRouter) createShortUrl(ctx *gin.Context) {
	payload := createShortUrlDto{}
	err := ctx.BindJSON(&payload)
	if err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			out := make([]api.ErrorDto, len(ve))
			for i, e := range ve {
				out[i] = api.ErrorDto{Field: e.Field(), Message: api.ParseDtoFieldError(e)}
			}
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": out})
			return
		}
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userId, ok := ctx.Get("userId")
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if err := validators.IsValidUrl(payload.OriginalUrl); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	uid := userId.(string)
	url, err := hdl.urlWriteSvc.CreateUrl(uid, payload.OriginalUrl, payload.CustomAlias, payload.ExpiresOn, payload.Keywords)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response := mapEntityToResponseDto(url)
	ctx.JSON(http.StatusCreated, response)
}

// getUrlById returns a url that is attached to a user
func (hdl *urlRouter) getUrlById(ctx *gin.Context) {
	urlId := ctx.Param("id")

	url, err := hdl.urlReadSvc.GetById(urlId)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	response := mapEntityToResponseDto(url)

	ctx.JSON(http.StatusOK, response)
}

func (hdl *urlRouter) getAllUrls(ctx *gin.Context) {
	userId, ok := ctx.Get("userId")
	if !ok {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	uid := userId.(string)
	urls, err := hdl.urlReadSvc.GetByUserId(uid)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response := []urlResponseDto{}
	for _, url := range urls {
		response = append(response, mapEntityToResponseDto(url))
	}

	ctx.JSON(http.StatusOK, response)
}

func (hdl *urlRouter) deleteUrl(ctx *gin.Context) {
	urlId := ctx.Param("id")

	err := hdl.urlWriteSvc.Remove(urlId)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Url with ID %s has been deleted", urlId)})
}

// updateUrl is a handler to update an existing short url
func (hdl *urlRouter) updateUrl(ctx *gin.Context) {
	userId, ok := ctx.Get("userId")
	if !ok {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	urlId := ctx.Param("id")

	if _, err := hdl.urlReadSvc.GetById(urlId); err != nil {
		ctx.AbortWithStatusJSON(http.StatusNotFound, map[string]any{"message": err.Error()})
		return
	}

	payload := updateShortUrlDto{}
	err := ctx.BindJSON(&payload)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	customAlias := payload.CustomAlias
	keywords := payload.Keywords
	expiresOn := payload.ExpiresOn

	updateCmd, err := contracts.NewUpdateUrlRequest(userId.(string), urlId, customAlias, keywords, expiresOn)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, map[string]any{"message": err.Error()})
		return
	}

	updatedUrl, err := hdl.urlWriteSvc.UpdateUrl(updateCmd)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}

	response := mapEntityToResponseDto(updatedUrl)

	ctx.JSON(http.StatusAccepted, response)
}

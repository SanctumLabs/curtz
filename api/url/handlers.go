package url

import (
	"github.com/gin-gonic/gin"
	"github.com/sanctumlabs/curtz/internal/core/contracts"
)

type UrlHandler struct {
	urlService contracts.UrlService
}

func NewUrlHandler(urlService contracts.UrlService) *UrlHandler {
	return &UrlHandler{urlService}
}

func (hdl *UrlHandler) CreateUrl(c *gin.Context) {
	request := CreateUrlDto{}
	c.BindJSON(&request)

	url, err := hdl.urlService.CreateUrl(request.owner, request.originalUrl, request.shortenedUrl)

	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, url)
}

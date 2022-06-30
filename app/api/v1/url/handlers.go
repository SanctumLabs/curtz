package url

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sanctumlabs/curtz/app/pkg/validators"
)

// createShortUrl creates a new shortened url
func (hdl *urlRouter) createShortUrl(c *gin.Context) {
	payload := createShortUrlDto{}
	err := c.BindJSON(&payload)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userId, ok := c.Get("userId")
	if !ok {
		c.Status(http.StatusUnauthorized)
		return
	}

	if err := validators.IsValidUrl(payload.OriginalUrl); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	uid := userId.(string)
	url, err := hdl.svc.CreateUrl(uid, payload.OriginalUrl, payload.CustomAlias, payload.ExpiresOn, payload.Keywords)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response := mapEntityToResponseDto(url)
	c.JSON(http.StatusCreated, response)
}

func (hdl *urlRouter) getUrlById(c *gin.Context) {
	urlId := c.Param("id")

	url, err := hdl.svc.GetById(urlId)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	response := mapEntityToResponseDto(url)

	c.JSON(http.StatusOK, response)
}

func (hdl *urlRouter) getUrlByShortCode(c *gin.Context) {
	urlId := c.Param("id")

	url, err := hdl.svc.GetById(urlId)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	response := mapEntityToResponseDto(url)

	c.JSON(http.StatusOK, response)
}

func (hdl *urlRouter) getAllUrls(c *gin.Context) {
	userId, ok := c.Get("userId")
	if !ok {
		c.Status(http.StatusUnauthorized)
		return
	}

	uid := userId.(string)
	urls, err := hdl.svc.GetByUserId(uid)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response := []urlResponseDto{}
	for _, url := range urls {
		response = append(response, mapEntityToResponseDto(url))
	}

	c.JSON(http.StatusOK, response)
}

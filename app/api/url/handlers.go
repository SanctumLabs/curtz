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

	if err := validators.IsValidUrl(payload.OriginalUrl); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err = hdl.svc.CreateUrlShortCode(payload.OriginalUrl)
}

func (hdl *urlRouter) getUrl(c *gin.Context) {
	request := createShortUrlDto{}
	err := c.BindJSON(&request)
	if err != nil {
		return
	}

	url, err := hdl.svc.GetByOriginalUrl(request.OriginalUrl)

	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, url)
}

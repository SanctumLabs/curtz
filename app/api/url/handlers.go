package url

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// createShortUrl creates a short url from a long url
func (hdl urlRouter) createShortUrl(c *gin.Context) {
	request := shortedUrlRequestDto{}
	err := c.BindJSON(&request)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err = hdl.svc.CreateUrlShortCode(request.url)
}

func (hdl *urlRouter) createUrl(c *gin.Context) {
	request := CreateUrlDto{}
	err := c.BindJSON(&request)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	//url, err := hdl.svc.CreateUrl(request.owner, request.originalUrl)

	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}
	//c.JSON(200, url)
}

func (hdl *urlRouter) getUrl(c *gin.Context) {
	request := CreateUrlDto{}
	err := c.BindJSON(&request)
	if err != nil {
		return
	}

	url, err := hdl.svc.GetByOriginalUrl(request.originalUrl)

	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, url)
}

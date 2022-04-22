package url

import (
	"github.com/gin-gonic/gin"
)

func (hdl *urlRouter) createUrl(c *gin.Context) {
	request := CreateUrlDto{}
	err := c.BindJSON(&request)
	if err != nil {
		return
	}

	url, err := hdl.svc.CreateUrl(request.owner, request.originalUrl, request.shortenedUrl)

	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, url)
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

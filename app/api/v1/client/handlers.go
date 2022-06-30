package client

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (hdl *clientRouter) handleRedirect(c *gin.Context) {
	shortCode := c.Param("shortCode")

	url, err := hdl.svc.GetByShortCode(shortCode)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.Redirect(http.StatusPermanentRedirect, url.OriginalUrl)
}

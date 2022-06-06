package url

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sanctumlabs/curtz/app/pkg"
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

	url, err := hdl.svc.CreateUrl("", payload.OriginalUrl, payload.CustomAlias, payload.ExpiresOn, payload.Keywords)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response := urlResponseDto{
		Id:          url.ID.String(),
		UserId:      url.UserId.String(),
		OriginalUrl: url.OriginalUrl,
		CustomAlias: url.CustomAlias,
		ShortCode:   url.ShortCode,
		Keywords:    payload.Keywords,
		ExpiresOn:   url.ExpiresOn.Format(pkg.DateLayout),
		DeletedAt:   url.DeletedAt.Format(pkg.DateLayout),
		CreatedAt:   url.CreatedAt.Format(pkg.DateLayout),
		UpdatedAt:   url.UpdatedAt.Format(pkg.DateLayout),
		Hits:        url.Hits,
	}

	c.JSON(http.StatusCreated, response)
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

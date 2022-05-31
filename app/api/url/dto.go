package url

// CreateUrlDto is the request body for creating a url
type createShortUrlDto struct {
	// OriginalUrl is the original url to shorten
	OriginalUrl string `json:"original_url" binding:"required"`

	// CustomAlias is the custom alias for the url provided for by the user
	CustomAlias string `json:"custom_alias" default:""`

	// ExpiresOn is the expiration date for the url
	ExpiresOn string `json:"expires_on" default:""`

	// Keywords is the list of keywords to be attached to the url
	Keywords []string `json:"keywords"`

	// Host is the host of the url
	Host string `json:"-"`
}

// createShortUrlResponseDto is the response body for creating a url
type createShortUrlResponseDto struct {
	// OriginalUrl is the original url to shorten
	OriginalUrl string `json:"original_url"`

	// CustomAlias is the custom alias for the url provided for by the user
	ShortenedUrl string `json:"shortend_url"`

	// ExpiresOn is the expiration date for the url
	ExpiresOn string `json:"expires_on"`

	// Keywords is the list of keywords to be attached to the url
	Keywords []string `json:"keywords"`

	// Host is the host of the url
	Host string `json:"-"`
}

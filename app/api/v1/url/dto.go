package url

import "time"

//urlResponseDto is the response dto for the url service
type urlResponseDto struct {
	// Id is the id of the created url to shorten
	Id string `json:"id"`

	// UserId is the id of the user whow created this url to shorten
	UserId string `json:"user_id"`

	// OriginalUrl is the original url to shorten
	OriginalUrl string `json:"original_url"`

	// CustomAlias is the custom alias for the url provided for by the user
	CustomAlias string `json:"custom_alias" default:""`

	// ShortCode is the short code generated for the url provided for by the user
	ShortCode string `json:"short_code"`

	// ExpiresOn is the expiration date for the url
	ExpiresOn time.Time `json:"expires_on" default:""`

	// Deletedat is the deleted at date for the url
	DeletedAt *time.Time `json:"deleted_at" default:""`

	// CreatedAt is the created date for the url
	CreatedAt time.Time `json:"created_at" default:""`

	// UpdatedAt is the updated date for the url
	UpdatedAt time.Time `json:"updated_at" default:""`

	// Keywords is the list of keywords to be attached to the url
	Keywords []string `json:"keywords"`

	// Hits is the number of hits of the url
	Hits uint `json:"hits"`
}

// CreateUrlDto is the request body for creating a url
type createShortUrlDto struct {
	// OriginalUrl is the original url to shorten
	OriginalUrl string `json:"original_url" binding:"required"`

	// CustomAlias is the custom alias for the url provided for by the user
	CustomAlias string `json:"custom_alias" default:""`

	// ExpiresOn is the expiration date for the url
	ExpiresOn string `json:"expires_on" binding:"required"`

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

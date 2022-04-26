package url

import "github.com/google/uuid"

type CreateUrlDto struct {
	owner       uuid.UUID
	originalUrl string
}

type shortedUrlRequestDto struct {
	url       string   `json:"url" binding:"required"`
	expiresOn string   `json:"expires_on"`
	keywords  []string `json:"keywords"`
	host      string   `json:"-"`
}

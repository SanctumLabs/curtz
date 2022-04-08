package url

import "github.com/google/uuid"

type CreateUrlDto struct {
	owner                     uuid.UUID
	originalUrl, shortenedUrl string
}

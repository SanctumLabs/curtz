package url

import "github.com/sanctumlabs/curtz/app/internal/core/entities"

func mapEntityToResponseDto(url entities.URL) urlResponseDto {
	keywords := []string{}

	for _, keyword := range url.Keywords {
		value := keyword.Value
		if value != "" {
			keywords = append(keywords, value)
		}
	}

	return urlResponseDto{
		Id:     url.ID.String(),
		UserId: url.UserId.String(),
		urlDto: urlDto{
			OriginalUrl: url.OriginalUrl,
			CustomAlias: url.CustomAlias,
			Keywords:    keywords,
			ExpiresOn:   url.ExpiresOn,
		},
		ShortCode: url.ShortCode,
		CreatedAt: url.CreatedAt,
		UpdatedAt: url.UpdatedAt,
		Hits:      url.Hits,
	}
}

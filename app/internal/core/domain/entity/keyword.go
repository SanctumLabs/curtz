package entity

import "github.com/sanctumlabs/curtz/app/pkg/identifier"

// Keyword is model for keywords attached to a url
type Keyword struct {
	identifier.ID
	Keyword string
}

// NewKeyword creates a new keyword
func NewKeyword(keyword string) Keyword {
	id := identifier.New()
	return Keyword{
		ID:      id,
		Keyword: keyword,
	}
}

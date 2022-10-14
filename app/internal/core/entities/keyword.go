package entities

import (
	"github.com/sanctumlabs/curtz/app/pkg"
	"github.com/sanctumlabs/curtz/app/pkg/errdefs"
	"github.com/sanctumlabs/curtz/app/pkg/identifier"
)

// Keyword is model for keywords attached to a url
type Keyword struct {
	identifier.ID
	value string
}

// NewKeyword creates a new keyword
func NewKeyword(keyword string) (*Keyword, error) {

	if l := len(keyword); l < 2 || l > 25 {
		return nil, errdefs.ErrKeywordLength
	}

	if !pkg.KeywordRegex.MatchString(keyword) {
		return nil, errdefs.ErrInvalidKeyword
	}

	id := identifier.New()
	return &Keyword{
		ID:    id,
		value: keyword,
	}, nil
}

// SetValue sets the value of a Keyword. Returns error if any
func (kw *Keyword) SetValue(value string) error {
	if l := len(value); l < 2 || l > 25 {
		return errdefs.ErrKeywordLength
	}

	if !pkg.KeywordRegex.MatchString(value) {
		return errdefs.ErrInvalidKeyword
	}

	kw.value = value
	return nil
}

// GetValue returns the value of a keyword
func (kw *Keyword) GetValue() string {
	return kw.value
}

func createKeywords(keywords []string) ([]Keyword, error) {
	kws := make([]Keyword, len(keywords))

	if len(keywords) > 10 {
		return kws, errdefs.ErrKeywordsCount
	}

	for _, kw := range keywords {
		if keyword, err := NewKeyword(kw); err == nil {
			kws = append(kws, *keyword)
		}
	}

	return kws, nil
}

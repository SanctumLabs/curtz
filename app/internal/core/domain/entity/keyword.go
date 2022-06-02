package entity

import (
	"regexp"

	"github.com/sanctumlabs/curtz/app/pkg/errdefs"
	"github.com/sanctumlabs/curtz/app/pkg/identifier"
)

var (
	KeywordRegex = `^[a-zA-Z0-9-_]+$`
	kwRe         = regexp.MustCompile(KeywordRegex)
)

// Keyword is model for keywords attached to a url
type Keyword struct {
	identifier.ID
	Keyword string
}

// NewKeyword creates a new keyword
func NewKeyword(keyword string) (Keyword, error) {

	if l := len(keyword); l < 2 || l > 25 {
		return Keyword{}, errdefs.ErrKeywordLength
	}

	if !kwRe.MatchString(keyword) {
		return Keyword{}, errdefs.ErrInvalidKeyword
	}

	id := identifier.New()
	return Keyword{
		ID:      id,
		Keyword: keyword,
	}, nil
}

func createKeywords(keywords []string) ([]Keyword, error) {
	kws := make([]Keyword, len(keywords))

	if len(keywords) > 10 {
		return kws, errdefs.ErrKeywordsCount
	}

	for idx, kw := range keywords {
		if kw, err := NewKeyword(kw); err != nil {
			kws[idx] = kw
		}
	}

	return kws, nil
}

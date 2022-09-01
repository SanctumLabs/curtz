package entities

import (
	"testing"
	"time"

	"github.com/sanctumlabs/curtz/app/pkg/errdefs"
	"github.com/sanctumlabs/curtz/app/pkg/identifier"
)

type urlTestCase struct {
	name          string
	userId        identifier.ID
	url           string
	alias         string
	expiresOn     time.Time
	keywords      []string
	expectedError error
}

var urlTestCases = []urlTestCase{
	{
		name:          "original url with invalid length should return error",
		userId:        identifier.New(),
		url:           "http://google",
		alias:         "",
		expiresOn:     time.Now().Add(time.Hour * 2),
		keywords:      []string{},
		expectedError: errdefs.ErrInvalidURLLen,
	},
	{
		name:          "original url with invalid localhost as host should return error",
		userId:        identifier.New(),
		url:           "http://localhost:8080",
		alias:         "",
		expiresOn:     time.Now().Add(time.Hour * 2),
		keywords:      []string{},
		expectedError: errdefs.ErrFilteredURL,
	},
	{
		name:          "original url with invalid 127.0.0.1 as host should return error",
		userId:        identifier.New(),
		url:           "http://127.0.0.1:8080",
		alias:         "",
		expiresOn:     time.Now().Add(time.Hour * 2),
		keywords:      []string{},
		expectedError: errdefs.ErrFilteredURL,
	},
	{
		name:          "original url with invalid scheme should return error",
		userId:        identifier.New(),
		url:           "xyz://xyx/xyz/z",
		alias:         "",
		expiresOn:     time.Now().Add(time.Hour * 2),
		keywords:      []string{},
		expectedError: errdefs.ErrInvalidURL,
	},
	{
		name:          "original url with filtered url should return error",
		userId:        identifier.New(),
		url:           "http://localhost/xxx",
		alias:         "",
		expiresOn:     time.Now().Add(time.Hour * 2),
		keywords:      []string{},
		expectedError: errdefs.ErrFilteredURL,
	},
	{
		name:          "should not return error with valid expires on date",
		userId:        identifier.New(),
		url:           "http://example.com",
		alias:         "",
		expiresOn:     time.Now().Add(time.Hour * 1),
		keywords:      []string{},
		expectedError: nil,
	},
}

func TestNewUrl(t *testing.T) {
	for _, tc := range urlTestCases {
		t.Run(tc.name, func(t *testing.T) {
			url, err := NewUrl(tc.userId, tc.url, tc.alias, tc.expiresOn, tc.keywords)
			if err != tc.expectedError {
				t.Errorf("NewUrl(%v, %s, %s, %s, %v) = (%v, %v) expected error: %v, got: %v", tc.userId, tc.url, tc.alias, tc.expiresOn, tc.keywords, url, err, tc.expectedError, err)
			}
		})
	}
}

func BenchmarkNewUrl(b *testing.B) {
	if testing.Short() {
		b.Skip("skipping benchmark")
	}

	for _, tc := range urlTestCases {
		b.Run(tc.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = NewUrl(tc.userId, tc.url, tc.alias, tc.expiresOn, tc.keywords)
			}
		})
	}
}

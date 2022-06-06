package entities

import (
	"testing"

	"github.com/sanctumlabs/curtz/app/pkg/errdefs"
	"github.com/sanctumlabs/curtz/app/pkg/identifier"
)

type urlTestCase struct {
	name          string
	userId        identifier.ID
	url           string
	alias         string
	expiresOn     string
	keywords      []string
	expectedError error
}

var urlTestCases = []urlTestCase{
	{
		name:          "original url with invalid length should return error",
		userId:        identifier.New(),
		url:           "http://google",
		alias:         "",
		expiresOn:     "",
		keywords:      []string{},
		expectedError: errdefs.ErrInvalidURLLen,
	},
	{
		name:          "original url with invalid localhost as host should return error",
		userId:        identifier.New(),
		url:           "http://localhost:8080",
		alias:         "",
		expiresOn:     "",
		keywords:      []string{},
		expectedError: errdefs.ErrFilteredURL,
	},
	{
		name:          "original url with invalid 127.0.0.1 as host should return error",
		userId:        identifier.New(),
		url:           "http://127.0.0.1:8080",
		alias:         "",
		expiresOn:     "",
		keywords:      []string{},
		expectedError: errdefs.ErrFilteredURL,
	},
	{
		name:          "original url with invalid scheme should return error",
		userId:        identifier.New(),
		url:           "xyz://xyx/xyz/z",
		alias:         "",
		expiresOn:     "",
		keywords:      []string{},
		expectedError: errdefs.ErrInvalidURL,
	},
	{
		name:          "original url with filtered url should return error",
		userId:        identifier.New(),
		url:           "http://localhost/xxx",
		alias:         "",
		expiresOn:     "",
		keywords:      []string{},
		expectedError: errdefs.ErrFilteredURL,
	},
	{
		name:          "should return error with invalid expires on date",
		userId:        identifier.New(),
		url:           "http://example.com",
		alias:         "",
		expiresOn:     "2030x01x01x00x00x00",
		keywords:      []string{},
		expectedError: errdefs.ErrInvalidDate,
	},
	{
		name:          "should return error with expires on date missing time",
		userId:        identifier.New(),
		url:           "http://example.com",
		alias:         "",
		expiresOn:     "2030-01-01",
		keywords:      []string{},
		expectedError: errdefs.ErrInvalidDate,
	},
	{
		name:          "should not return error with valid expires on date",
		userId:        identifier.New(),
		url:           "http://example.com",
		alias:         "",
		expiresOn:     "2030-01-01 00:00:00",
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
				NewUrl(tc.userId, tc.url, tc.alias, tc.expiresOn, tc.keywords)
			}
		})
	}
}

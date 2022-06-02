package validators

import (
	"testing"

	"github.com/sanctumlabs/curtz/app/pkg/errdefs"
)

type testCase struct {
	name        string
	url         string
	expectedErr error
}

var testCases = []testCase{
	{
		name:        "valid url with protocol",
		url:         "https://www.google.com",
		expectedErr: nil,
	},
	{
		name:        "valid url without protocol",
		url:         "www.google.com",
		expectedErr: nil,
	},
	{
		name:        "invalid url with invalid protocol",
		url:         "htt://www.google.com",
		expectedErr: errdefs.ErrInvalidURL,
	},
}

func TestIsValidUrl(t *testing.T) {
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := IsValidUrl(tc.url)
			if err != tc.expectedErr {
				t.Errorf("IsValidUrl(%s) = %v expected error %v, got %v", tc.url, err, tc.expectedErr, err)
			}
		})
	}
}

func BenchmarkIsValidUrl(b *testing.B) {
	if testing.Short() {
		b.Skip("skipping benchmark")
	}

	for i := 0; i < b.N; i++ {
		for _, tc := range testCases {
			IsValidUrl(tc.url)
		}
	}
}

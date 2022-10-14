package entities

import (
	"testing"

	"github.com/sanctumlabs/curtz/app/pkg/errdefs"
	"github.com/stretchr/testify/assert"
)

type keywordTestCase struct {
	name  string
	input string
	err   error
}

var keywordTestCases = []keywordTestCase{
	{
		name:  "empty keyword should return nil and error",
		input: "",
		err:   errdefs.ErrKeywordLength,
	},
	{
		name:  "invalid keyword should return nil and error",
		input: "`'/?.;,<]{}",
		err:   errdefs.ErrInvalidKeyword,
	},
	{
		name:  "valid keyword should return keyword and nil error",
		input: "Social",
		err:   nil,
	},
}

func TestNewKeyword(t *testing.T) {
	for _, tc := range keywordTestCases {
		t.Run(tc.name, func(t *testing.T) {
			keyword, err := NewKeyword(tc.input)
			if err != tc.err {
				t.Errorf("NewKeyword(%s) = (%v, %v), expected error %v, got %v", tc.input, keyword, err, tc.err, err)
			}
		})
	}
}

func BenchmarkNewKeyword(b *testing.B) {
	if testing.Short() {
		b.Skip("skipping benchmark in short mode")
	}

	for i := 0; i < b.N; i++ {
		for _, tc := range keywordTestCases {
			_, _ = NewKeyword(tc.input)
		}
	}
}

type setKeywordTestCase struct {
	keywordTestCase
	initialKw string
}

var setKeywordTestCases = []setKeywordTestCase{
	{
		keywordTestCase: keywordTestCase{
			name:  "empty keyword should return nil and error",
			input: "",
			err:   errdefs.ErrKeywordLength,
		},
		initialKw: "Brand",
	},
	{
		keywordTestCase: keywordTestCase{
			name:  "invalid keyword should return nil and error",
			input: "`'/?.;,<]{}",
			err:   errdefs.ErrInvalidKeyword,
		},
		initialKw: "google",
	},
	{
		keywordTestCase: keywordTestCase{
			name:  "valid keyword should return keyword and nil error",
			input: "Social",
			err:   nil,
		},
		initialKw: "NewNew",
	},
}

func TestSetKeyword(t *testing.T) {
	for _, tc := range setKeywordTestCases {
		t.Run(tc.name, func(t *testing.T) {
			keyword, err := NewKeyword(tc.initialKw)
			assert.NoError(t, err)

			actualErr := keyword.SetValue(tc.input)

			if tc.err != actualErr {
				t.Errorf("SetValue(%s) = %v, expected error %v, got %v", tc.input, actualErr, tc.err, actualErr)
			}
		})
	}
}

func BenchmarkSetKeyword(b *testing.B) {
	if testing.Short() {
		b.Skip("skipping benchmark in short mode")
	}

	for i := 0; i < b.N; i++ {
		for _, tc := range setKeywordTestCases {
			keyword, _ := NewKeyword(tc.initialKw)
			_ = keyword.SetValue(tc.input)
		}
	}
}

type createKeywordsTestCase struct {
	name  string
	input []string
	err   error
}

var createKeywordsTestCases = []createKeywordsTestCase{
	{
		name:  "empty keywords should return empty keywords and nil error",
		input: []string{},
		err:   nil,
	},
	{
		name:  "More than 10 keywords should return nil and error",
		input: []string{"Social", "Media", "Entertainment", "Sports", "Technology", "Politics", "Business", "Science", "Health", "Travel", "Food", "Fashion"},
		err:   errdefs.ErrKeywordsCount,
	},
	{
		name:  "valid keyword length return keywords and nil error",
		input: []string{"Social", "Media", "Entertainment", "Sports", "Technology", "Politics", "Business", "Science", "Health", "Travel"},
		err:   nil,
	},
}

func TestCreateKeywords(t *testing.T) {
	for _, tc := range createKeywordsTestCases {
		t.Run(tc.name, func(t *testing.T) {
			keywords, err := createKeywords(tc.input)
			if err != tc.err {
				t.Errorf("createKeywords(%s) = (%v, %v), expected error %v, got %v", tc.input, keywords, err, tc.err, err)
			}
		})
	}
}

func BenchmarkCreateKeywords(b *testing.B) {
	if testing.Short() {
		b.Skip("skipping benchmark in short mode")
	}

	for i := 0; i < b.N; i++ {
		for _, tc := range createKeywordsTestCases {
			_, _ = createKeywords(tc.input)
		}
	}
}

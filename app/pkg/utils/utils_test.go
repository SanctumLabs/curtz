package utils

import "testing"

type emailValidTest struct {
	name     string
	email    string
	expected bool
}

var emailValidityTests = []emailValidTest{
	{
		name:     "johndoe@example.com",
		email:    "johndoe@example.com",
		expected: true,
	},
	{
		name:     "invalid email johndoe@example",
		email:    "johndoe@example",
		expected: false,
	},
}

func TestIsEmailValid(t *testing.T) {
	for _, tc := range emailValidityTests {
		t.Run(tc.name, func(t *testing.T) {
			actual := IsEmailValid(tc.email)

			if actual != tc.expected {
				t.Errorf("IsEmailValid(%s) = %v, expected=%v", tc.email, actual, tc.expected)
			}
		})
	}
}

func BenchmarkIsEmailValid(b *testing.B) {
	if testing.Short() {
		b.Skip("skipping benchmark")
	}

	for i := 0; i < b.N; i++ {
		for _, tc := range emailValidityTests {
			IsEmailValid(tc.email)
		}
	}
}

package helpers

import (
	"testing"
)

func TestQuote(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Already double quoted",
			input:    `"foo bar"`,
			expected: `"foo bar"`,
		},
		{
			name:     "Already single quoted",
			input:    `'foo bar'`,
			expected: `'foo bar'`,
		},
		{
			name:     "Contains escaped space",
			input:    `foo\ bar`,
			expected: `foo\ bar`,
		},
		{
			name:     "No space, no quote needed",
			input:    "foobar",
			expected: "foobar",
		},
		{
			name:     "Tilde-prefixed path with spaces",
			input:    "~/foo bar/baz",
			expected: `~/foo\ bar/baz`,
		},
		{
			name:     "Normal string with spaces",
			input:    "foo bar",
			expected: `"foo bar"`,
		},
		{
			name:     "String with quotes",
			input:    `he said "hello"`,
			expected: `"he said \"hello\""`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Quote(tt.input)
			if result != tt.expected {
				t.Errorf("Quote(%q) = %q; expected %q", tt.input, result, tt.expected)
			}
		})
	}
}

package data

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEscapeString(t *testing.T) {
	t.Parallel()

	cases := []struct {
		original     string
		replacements []Replacement
		escaped      string
	}{
		{
			original:     "abc",
			replacements: []Replacement{},
			escaped:      "abc",
		},
		{
			original:     "abc",
			replacements: []Replacement{{"z", "b"}},
			escaped:      "abc",
		},
		{
			original:     "abc",
			replacements: []Replacement{{"b", "a"}},
			escaped:      "aac",
		},
		{
			original:     "abc",
			replacements: []Replacement{{"b", "a"}, {"b", "a"}},
			escaped:      "aac",
		},
		{
			original:     "abc",
			replacements: []Replacement{{"a", "b"}, {"b", "c"}},
			escaped:      "ccc",
		},
	}

	for _, testCase := range cases {
		assert.Equal(
			t,
			testCase.escaped,
			EscapeString(testCase.original, testCase.replacements),
		)
	}
}

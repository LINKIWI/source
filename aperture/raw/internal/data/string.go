package data

import (
	"strings"
)

// Replacement is a rule describing how a substring should be escaped.
type Replacement struct {
	From string
	To   string
}

// EscapeString escapes a set of special characters using an ordered list of substring replacements.
// Replacements are executed in the order in which they are provided.
func EscapeString(original string, replacements []Replacement) string {
	escaped := original

	for _, replacement := range replacements {
		escaped = strings.ReplaceAll(escaped, replacement.From, replacement.To)
	}

	return escaped
}

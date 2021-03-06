package cli

import (
	"fmt"
	"strconv"

	"github.com/buger/goterm"
)

// Bold returns asn ANSI-encoded sequence to boldface a particular string in a console.
func Bold(text string) string {
	return goterm.Bold(text)
}

// Highlight returns an ANSI-encoded sequence to highlight a particular string in a console.
func Highlight(text string) string {
	return goterm.Bold(goterm.Color(text, goterm.GREEN))
}

// Pad inserts whitespace after a string such that the resulting string is equal to an exact length.
func Pad(text string, length int) string {
	return fmt.Sprintf("%-"+strconv.Itoa(length)+"v", text)
}

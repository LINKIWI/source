package config

import (
	"fmt"
	"path/filepath"
	"regexp"
	"time"
)

// glob is a TOML unmarshaler implementation for a string with globbing convenience methods.
type glob struct {
	Pattern string
}

// Match attempts to return exactly one match from the glob pattern, and errors when not possible.
func (g *glob) Match() (string, error) {
	matches, err := g.Matches()
	if err != nil {
		return "", err
	} else if len(matches) == 0 {
		return "", fmt.Errorf("config: glob pattern matched no files")
	} else if len(matches) > 1 {
		return "", fmt.Errorf(
			"config: glob pattern matched more than one file: matches=%v",
			matches,
		)
	}

	return matches[0], nil
}

// Matches returns the names of all files that match the glob pattern.
func (g *glob) Matches() ([]string, error) {
	return filepath.Glob(g.Pattern)
}

// UnmarshalText stores the string value as a glob pattern.
func (g *glob) UnmarshalText(text []byte) error {
	g.Pattern = string(text)

	return nil
}

// duration is a TOML unmarshaler implementation for a time.Duration.
type duration struct {
	time.Duration
}

// UnmarshalText attempts to parse the input string into a time.Duration.
func (d *duration) UnmarshalText(text []byte) error {
	var err error

	d.Duration, err = time.ParseDuration(string(text))

	return err
}

// regex is a TOML unmarshaler implementation for a regexp.Regexp.
type regex struct {
	*regexp.Regexp
}

// UnmarshalText attempts to parse a non-empty input string into a regexp.Regexp.
// Empty input strings will leave the regexp.Regexp as nil.
func (r *regex) UnmarshalText(text []byte) error {
	var err error

	if string(text) == "" {
		return nil
	}

	r.Regexp, err = regexp.Compile(string(text))

	return err
}

package config

import (
	"fmt"
	"net"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// address is a TOML unmarshaler implementation for a net.Addr.
type address struct {
	network string
	address string
	net.Addr
}

// UnmarshalText attempts to parse the input string into the network and address components
// comprising a net.Addr. It supports address specifications for TCP and Unix socket transports.
func (a *address) UnmarshalText(text []byte) error {
	value := string(text)

	switch {
	case strings.HasPrefix(value, "unix:"):
		a.network = "unix"
		a.address = value[len("unix:"):]
	case strings.HasPrefix(value, "tcp:"):
		a.network = "tcp"
		a.address = value[len("tcp:"):]
	case strings.HasPrefix(value, "tcp4:"):
		a.network = "tcp4"
		a.address = value[len("tcp4:"):]
	case strings.HasPrefix(value, "tcp6:"):
		a.network = "tcp6"
		a.address = value[len("tcp6:"):]
	default:
		return fmt.Errorf("config: error parsing address format: address=%s", value)
	}

	switch a.network {
	case "tcp", "tcp4", "tcp6":
		if _, _, err := net.SplitHostPort(a.address); err != nil {
			return fmt.Errorf(
				"config: TCP address is malformed: address=%s err=%v",
				value,
				err,
			)
		}
	}

	return nil
}

// Network returns the address's network transport alias.
func (a address) Network() string {
	return a.network
}

// String returns the address itself, the format of which depends on the transport.
func (a address) String() string {
	return a.address
}

// Spec returns the address as a formatted specification including both the network and address.
func (a address) Spec() string {
	if a.Network() == "" && a.String() == "" {
		return "<nil>"
	}

	return fmt.Sprintf("%s:%s", a.Network(), a.String())
}

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

package config

import (
	"fmt"
	"net"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

// Address is a YAML unmarshaler implementation for a net.Addr.
type Address struct {
	network string
	address string
	net.Addr
}

// UnmarshalYAML attempts to parse the input string into the network and address components
// comprising a net.Addr. It supports address specifications for TCP and Unix socket transports.
func (a *Address) UnmarshalYAML(node *yaml.Node) error {
	switch {
	case strings.HasPrefix(node.Value, "unix:"):
		a.network = "unix"
		a.address = node.Value[len("unix:"):]
	case strings.HasPrefix(node.Value, "tcp:"):
		a.network = "tcp"
		a.address = node.Value[len("tcp:"):]
	case strings.HasPrefix(node.Value, "tcp4:"):
		a.network = "tcp4"
		a.address = node.Value[len("tcp4:"):]
	case strings.HasPrefix(node.Value, "tcp6:"):
		a.network = "tcp6"
		a.address = node.Value[len("tcp6:"):]
	default:
		return fmt.Errorf("config: error parsing address format: address=%s", node.Value)
	}

	switch a.network {
	case "tcp", "tcp4", "tcp6":
		if _, _, err := net.SplitHostPort(a.address); err != nil {
			return fmt.Errorf(
				"config: TCP address is malformed: address=%s err=%v",
				node.Value,
				err,
			)
		}
	}

	return nil
}

// Network returns the address's network transport alias.
func (a Address) Network() string {
	return a.network
}

// String returns the address itself, the format of which depends on the transport.
func (a Address) String() string {
	return a.address
}

// Spec returns the address as a formatted specification including both the network and address.
func (a Address) Spec() string {
	if a.Network() == "" && a.String() == "" {
		return "<nil>"
	}

	return fmt.Sprintf("%s:%s", a.Network(), a.String())
}

// Regex is YAML node with a regular expression.
type Regex struct {
	*regexp.Regexp
}

// UnmarshalYAML attempts to parse the node's string contents as either a regular expression or a
// string literal.
func (r *Regex) UnmarshalYAML(node *yaml.Node) error {
	var value string
	var err error

	switch {
	// String must be non-empty.
	case len(node.Value) == 0:
		return fmt.Errorf("config: regular expression is empty")
	// Strings wrapped with forward slashes are interpreted as standard regular expressions.
	case len(node.Value) > 2 && strings.HasPrefix(node.Value, "/") && strings.HasSuffix(node.Value, "/"):
		value = node.Value[1 : len(node.Value)-1]
	// All other strings are interpreted as a string literal compatible only with exact matches.
	default:
		value = fmt.Sprintf("^%s$", regexp.QuoteMeta(node.Value))
	}

	if r.Regexp, err = regexp.Compile(value); err != nil {
		return fmt.Errorf("config: failed to parse regular expression: value=%s", value)
	}

	return nil
}

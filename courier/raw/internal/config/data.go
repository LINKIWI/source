package config

import (
	"fmt"
	"net"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"

	"courier/internal/util"
)

// Regex is YAML node with a regular expression.
type Regex struct {
	negated bool
	*regexp.Regexp
}

// UnmarshalYAML attempts to parse the node's string contents as either a regular expression or a
// string literal.
func (r *Regex) UnmarshalYAML(node *yaml.Node) error {
	var value string
	var negated bool

	switch {
	// String must be non-empty.
	case len(node.Value) == 0:
		return &util.Error{
			Namespace: "config",
			Message:   "regular expression is empty",
		}
	// Strings wrapped with forward slashes are interpreted as standard regular expressions.
	case len(node.Value) > 2 && strings.HasPrefix(node.Value, "/") && strings.HasSuffix(node.Value, "/"):
		value = node.Value[1 : len(node.Value)-1]
	// Strings wrapped with backwards slashes are interpreted as negated regular expressions.
	case len(node.Value) > 2 && strings.HasPrefix(node.Value, "\\") && strings.HasSuffix(node.Value, "\\"):
		value = node.Value[1 : len(node.Value)-1]
		negated = true
	// All other strings are interpreted as a string literal compatible only with exact matches.
	default:
		value = fmt.Sprintf("^%s$", regexp.QuoteMeta(node.Value))
	}

	re, err := regexp.Compile(value)
	if err != nil {
		return &util.Error{
			Namespace:    "config",
			Message:      "failed to parse regular expression",
			Tags:         map[string]interface{}{"value": value},
			StackedError: err,
		}
	}

	r.negated = negated
	r.Regexp = re

	return nil
}

// Match matches the input byte array against the regular expression, respecting negation as needed.
func (r *Regex) Match(b []byte) bool {
	if r.negated {
		return !r.Regexp.Match(b)
	}

	return r.Regexp.Match(b)
}

// MatchString matches the input string against the regular expression, respecting negation as
// needed.
func (r *Regex) MatchString(s string) bool {
	if r.negated {
		return !r.Regexp.MatchString(s)
	}

	return r.Regexp.MatchString(s)
}

// URL is a YAML node with a URL.
type URL struct {
	*url.URL
}

// UnmarshalYAML attempts to parse the node's string contents as a URL.
func (u *URL) UnmarshalYAML(node *yaml.Node) error {
	parsed, err := url.Parse(node.Value)
	if err != nil {
		return &util.Error{
			Namespace:    "config",
			Message:      "failed to parse upstream URL",
			Tags:         map[string]interface{}{"upstream": node.Value},
			StackedError: err,
		}
	}

	u.URL = parsed

	return nil
}

// Address is a YAML node representing a network address.
type Address struct {
	Net    string
	Host   string
	Port   int
	Socket string
}

// UnmarshalYAML attempts to parse the node's string contents as a host ands port combination.
func (a *Address) UnmarshalYAML(node *yaml.Node) error {
	var address string

	switch {
	case strings.HasPrefix(node.Value, "unix:"):
		a.Net = "unix"
		address = node.Value[5:]
	case strings.HasPrefix(node.Value, "tcp4:"):
		a.Net = "tcp4"
		address = node.Value[5:]
	case strings.HasPrefix(node.Value, "tcp6:"):
		a.Net = "tcp6"
		address = node.Value[5:]
	case strings.HasPrefix(node.Value, "tcp:"):
		a.Net = "tcp"
		address = node.Value[4:]
	default:
		a.Net = "tcp"
		address = node.Value
	}

	if a.Net == "unix" {
		a.Socket = address
		return nil
	}

	host, port, err := net.SplitHostPort(address)
	if err != nil {
		return &util.Error{
			Namespace:    "config",
			Message:      "failed to parse address",
			Tags:         map[string]interface{}{"address": node.Value},
			StackedError: err,
		}
	}

	a.Host = host
	a.Port, _ = strconv.Atoi(port)

	return nil
}

// Address returns a dial target as a (network, address) pair.
func (a *Address) Address() (string, string) {
	return a.Net, a.String()
}

// String returns the address in standard form.
func (a *Address) String() string {
	switch a.Net {
	case "tcp", "tcp4", "tcp6":
		return net.JoinHostPort(a.Host, strconv.Itoa(a.Port))
	case "unix":
		return a.Socket
	default:
		return "<nil>"
	}
}

// CIDR is a YAML node representing an IP and subnet mask in CIDR notation.
type CIDR struct {
	negated bool
	*net.IPNet
}

// UnmarshalYAML attempts to parse the node's string contents as an IP address and subnet mask.
func (c *CIDR) UnmarshalYAML(node *yaml.Node) error {
	cidr := node.Value

	if len(cidr) > 0 && cidr[0] == '!' {
		c.negated = true
		cidr = cidr[1:]
	}

	_, ipnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return &util.Error{
			Namespace:    "config",
			Message:      "failed to parse CIDR block",
			Tags:         map[string]interface{}{"cidr": node.Value},
			StackedError: err,
		}
	}

	c.IPNet = ipnet

	return nil
}

// Contains returns whether the target IP address is contained within this CIDR block, respecting
// negotiation as needed.
func (c *CIDR) Contains(ip net.IP) bool {
	if c.negated {
		return !c.IPNet.Contains(ip)
	}

	return c.IPNet.Contains(ip)
}

// String returns the CIDR in standard form.
func (c CIDR) String() string {
	if c.negated {
		return fmt.Sprintf("!%s", c.IPNet.String())
	}

	return c.IPNet.String()
}

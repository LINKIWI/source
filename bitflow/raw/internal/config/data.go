package config

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

// Address is a YAML node representing a network address.
type Address struct {
	Net    string
	Host   string
	Port   int
	Socket string
	FD     int
}

// UnmarshalYAML attempts to parse the node's string contents as a host and port combination.
func (a *Address) UnmarshalYAML(node *yaml.Node) error {
	var address string

	switch {
	case strings.HasPrefix(node.Value, "fd:"):
		a.Net = "fd"
		address = node.Value[3:]
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

	if a.Net == "fd" {
		fd, err := strconv.Atoi(address)
		if err != nil {
			return fmt.Errorf("config: failed to parse file descriptor: err=%v", err)
		}

		a.FD = fd
		return nil
	}

	if a.Net == "unix" {
		a.Socket = address
		return nil
	}

	host, port, err := net.SplitHostPort(address)
	if err != nil {
		return fmt.Errorf("config: failed to parse address: err=%v", err)
	}

	a.Host = host
	a.Port, _ = strconv.Atoi(port)

	return nil
}

// Listen creates a net.Listener for the address on the requested network.
func (a *Address) Listen(tlsCfg *tls.Config) (net.Listener, error) {
	switch a.Net {
	case "fd":
		ln, err := net.FileListener(os.NewFile(uintptr(a.FD), a.String()))
		if err != nil {
			return nil, fmt.Errorf(
				"config: failed to create file listener from FD: fd=%v err=%v",
				a.FD,
				err,
			)
		}

		if tlsCfg != nil {
			return tls.NewListener(ln, tlsCfg), nil
		}

		return ln, nil
	default:
		network, addr := a.Address()

		if tlsCfg != nil {
			return tls.Listen(network, addr, tlsCfg)
		}

		return net.Listen(network, addr)
	}
}

// Address returns a dial target as a (network, address) pair.
func (a *Address) Address() (string, string) {
	switch a.Net {
	case "tcp", "tcp4", "tcp6":
		return a.Net, net.JoinHostPort(a.Host, strconv.Itoa(a.Port))
	case "unix":
		return a.Net, a.Socket
	case "fd":
		return a.Net, strconv.Itoa(a.FD)
	default:
		return "", ""
	}
}

// Resolve resolves the address as an explicit IPv4 or IPv6 address, as specified by the IP network.
func (a *Address) Resolve(ctx context.Context, ipNet string) (*Address, error) {
	var resolvedNet string

	switch a.Net {
	case "tcp":
	case "tcp4", "tcp6":
		return nil, fmt.Errorf("config: address is already resolved as explicit IPv4/IPv6")
	default:
		return nil, fmt.Errorf("config: network is not supported for IPv4/IPv6 resolution")
	}

	ips, err := net.DefaultResolver.LookupIP(ctx, ipNet, a.Host)
	if err != nil {
		return nil, fmt.Errorf(
			"config: error in IP resolution: network=%s host=%s err=%v",
			ipNet,
			a.Host,
			err,
		)
	}

	switch ipNet {
	case "ip4":
		resolvedNet = "tcp4"
	case "ip6":
		resolvedNet = "tcp6"
	default:
		resolvedNet = "tcp"
	}

	return &Address{
		Net:  resolvedNet,
		Host: ips[0].String(),
		Port: a.Port,
	}, nil
}

// String returns the address in standard form.
func (a *Address) String() string {
	network, addr := a.Address()
	return fmt.Sprintf("%s:%s", network, addr)
}

// CIDR is a YAML node representing an IP and subnet mask in CIDR notation.
type CIDR struct {
	*net.IPNet
}

// UnmarshalYAML attempts to parse the node's string contents as an IP address and subnet mask.
func (c *CIDR) UnmarshalYAML(node *yaml.Node) error {
	_, ipnet, err := net.ParseCIDR(node.Value)
	if err != nil {
		return fmt.Errorf("config: failed to parse CIDR: err=%v", err)
	}

	c.IPNet = ipnet

	return nil
}

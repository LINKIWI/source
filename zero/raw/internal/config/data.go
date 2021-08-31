package config

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"syscall"

	"gopkg.in/yaml.v3"
)

// Address is a YAML node representing a network address.
type Address struct {
	Net    string
	Host   string
	Port   int
	Socket string
}

// UnmarshalYAML attempts to parse the node's string contents as a host and port combination.
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
		return fmt.Errorf("config: failed to parse address: err=%v", err)
	}

	a.Host = host
	a.Port, _ = strconv.Atoi(port)

	return nil
}

// Listen creates a net.Listener for the address on the requested network.
func (a *Address) Listen() (net.Listener, error) {
	network, addr := a.Address()
	return net.Listen(network, addr)
}

// Address returns a dial target as a (network, address) pair.
func (a *Address) Address() (string, string) {
	switch a.Net {
	case "tcp", "tcp4", "tcp6":
		return a.Net, net.JoinHostPort(a.Host, strconv.Itoa(a.Port))
	case "unix":
		return a.Net, a.Socket
	default:
		return "", ""
	}
}

// String returns the address in standard form.
func (a *Address) String() string {
	network, addr := a.Address()
	return fmt.Sprintf("%s:%s", network, addr)
}

// Signal is a YAML node representing a process signal.
type Signal struct {
	key string
	os.Signal
}

// UnmarshalYAML attempts to parse the node's string contents as an os.Signal in a case-insensitive
// manner. All Unix process signals in the syscall package are supported.
func (s *Signal) UnmarshalYAML(node *yaml.Node) error {
	signals := map[string]os.Signal{
		"SIGABRT":   syscall.SIGABRT,
		"SIGALRM":   syscall.SIGALRM,
		"SIGBUS":    syscall.SIGBUS,
		"SIGCHLD":   syscall.SIGCHLD,
		"SIGCLD":    syscall.SIGCLD,
		"SIGCONT":   syscall.SIGCONT,
		"SIGFPE":    syscall.SIGFPE,
		"SIGHUP":    syscall.SIGHUP,
		"SIGILL":    syscall.SIGILL,
		"SIGINT":    syscall.SIGINT,
		"SIGIO":     syscall.SIGIO,
		"SIGIOT":    syscall.SIGIOT,
		"SIGKILL":   syscall.SIGKILL,
		"SIGPIPE":   syscall.SIGPIPE,
		"SIGPOLL":   syscall.SIGPOLL,
		"SIGPROF":   syscall.SIGPROF,
		"SIGPWR":    syscall.SIGPWR,
		"SIGQUIT":   syscall.SIGQUIT,
		"SIGSEGV":   syscall.SIGSEGV,
		"SIGSTKFLT": syscall.SIGSTKFLT,
		"SIGSTOP":   syscall.SIGSTOP,
		"SIGSYS":    syscall.SIGSYS,
		"SIGTERM":   syscall.SIGTERM,
		"SIGTRAP":   syscall.SIGTRAP,
		"SIGTSTP":   syscall.SIGTSTP,
		"SIGTTIN":   syscall.SIGTTIN,
		"SIGTTOU":   syscall.SIGTTOU,
		"SIGUNUSED": syscall.SIGUNUSED,
		"SIGURG":    syscall.SIGURG,
		"SIGUSR1":   syscall.SIGUSR1,
		"SIGUSR2":   syscall.SIGUSR2,
		"SIGVTALRM": syscall.SIGVTALRM,
		"SIGWINCH":  syscall.SIGWINCH,
		"SIGXCPU":   syscall.SIGXCPU,
		"SIGXFSZ":   syscall.SIGXFSZ,
	}

	signal, ok := signals[strings.ToUpper(node.Value)]
	if !ok {
		return fmt.Errorf("config: unknown or unsupported signal: signal=%s", node.Value)
	}

	s.key = strings.ToUpper(node.Value)
	s.Signal = signal

	return nil
}

// String returns the signal Unix name.
func (s *Signal) String() string {
	return s.key
}

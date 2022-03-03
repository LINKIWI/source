package server

import (
	"crypto/tls"
	"fmt"
	"net"
	"sync/atomic"

	"bitflow/internal/config"
	"bitflow/pkg/proxy"
)

// targetConnFactory is an alias for a factory that dials a target, given a net.Conn representing
// the source connection. It returns a net.Conn for the destination as well as an I/O pipe
// implementation of proxy.Pipe to use for the connection proxy.
type targetConnFactory func(net.Conn) (net.Conn, proxy.Pipe, error)

// targetCandidate is a pair of a target's configuration and its associated factory capable of
// supplying a connection for that target.
type targetCandidate struct {
	cfg     *config.Target
	factory targetConnFactory
}

// newFailoverLoadBalancer creates a load-balanced connection factory using a per-target failover
// policy. Each target is dialed in order until one successfully supplies a connection. An error is
// returned if all targets fail.
func newFailoverLoadBalancer(candidates []targetCandidate) targetConnFactory {
	return func(src net.Conn) (conn net.Conn, pipe proxy.Pipe, err error) {
		for _, candidate := range candidates {
			conn, pipe, err = candidate.factory(src)
			if err == nil {
				return conn, pipe, nil
			}
		}

		return nil, nil, fmt.Errorf(
			"server: failover load balancer exhausted all available candidate "+
				"targets without a successful connection: targets=%d err=%v",
			len(candidates),
			err,
		)
	}
}

// newRoundRobinBalancer creates a load-balanced connection factory with a round-robin rotation
// policy. The factory rotates through all targets fairly on each successive connection attempt.
func newRoundRobinBalancer(candidates []targetCandidate) targetConnFactory {
	id := new(uint64)

	return func(src net.Conn) (net.Conn, proxy.Pipe, error) {
		defer atomic.AddUint64(id, 1)

		return candidates[int(atomic.LoadUint64(id))%len(candidates)].factory(src)
	}
}

// newSNILoadBalancer creates a load-balanced connection factory for TLS connections that routes
// requests to targets based on the SNI (Server Name Indication) header in the client hello
// information revealed during the TLS handshake. The first target whose name is an exact match
// with the SNI value is selected to provide a connection. An error is returned if no targets match
// the client-indicated SNI.
func newSNILoadBalancer(candidates []targetCandidate) targetConnFactory {
	return func(src net.Conn) (net.Conn, proxy.Pipe, error) {
		srcTLSConn, ok := src.(*tls.Conn)
		if !ok {
			return nil, nil, fmt.Errorf(
				"server: cannot extract TLS metadata from source connection; "+
					"client transport is most likely not TLS: src=%v",
				src.RemoteAddr(),
			)
		}

		// At this point in the connection, a handshake most likely has not yet been
		// executed as there has not yet been any user I/O on the source connection. Force
		// completion of the TLS handshake at this stage in order to expose the SNI
		// requested by the client in its hello payload for purposes of target routing.
		if err := srcTLSConn.Handshake(); err != nil {
			return nil, nil, fmt.Errorf(
				"server: error in TLS handshake: src=%v err=%v",
				src.RemoteAddr(),
				err,
			)
		}

		state := srcTLSConn.ConnectionState()

		if state.ServerName == "" {
			return nil, nil, fmt.Errorf(
				"server: no SNI available from client TLS handshake: src=%v",
				src.RemoteAddr(),
			)
		}

		// Select the first matching target whose name is equal to the client-requested SNI.
		for _, candidate := range candidates {
			if candidate.cfg.Name == state.ServerName {
				return candidate.factory(src)
			}
		}

		return nil, nil, fmt.Errorf(
			"server: no targets match requested server name: src=%v sni=%s",
			src.RemoteAddr(),
			state.ServerName,
		)
	}
}

// newStaticLoadBalancer creates a connection factory that always uses the first available target
// to provide a connection.
func newStaticLoadBalancer(candidates []targetCandidate) targetConnFactory {
	return candidates[0].factory
}

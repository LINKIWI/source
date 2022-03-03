package server

import (
	"math"
	"net"
	"time"

	"golang.org/x/net/proxy"
)

// retryDialer is a proxy.Dialer that automatically retries dials with an exponential backoff retry
// policy, using a fixed number of total attempts and a seed backoff delay.
type retryDialer struct {
	attempts int
	backoff  time.Duration
	proxy.Dialer
}

// Dial uses the underlying proxy.Dialer to dial the specified network and address with built-in
// optional retry.
func (d *retryDialer) Dial(network string, addr string) (conn net.Conn, err error) {
	if d.attempts <= 1 {
		return d.Dialer.Dial(network, addr)
	}

	for i := 0; i < d.attempts; i++ {
		if i > 0 && d.backoff > 0 {
			time.Sleep(d.backoff * time.Duration(math.Pow(2, float64(i-1))))
		}

		conn, err = d.Dialer.Dial(network, addr)
		if err != nil {
			continue
		}

		return
	}

	return
}

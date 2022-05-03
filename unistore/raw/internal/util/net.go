package util

import (
	"context"
	"net"
	"time"

	"go.uber.org/zap"
)

const (
	defaultAttempts = 1
	defaultDelay    = 5 * time.Second
)

// RetryDialer is a net.Dialer abstraction that wraps multiple connection attempts with exponential
// backoff between retries.
type RetryDialer struct {
	attempts int
	delay    time.Duration
	net.Dialer
}

// NewRetryDialer creates a new retryDialer on top of an existing net.Dialer, with the requested
// number of total connection attempts and initial backoff time delay.
func NewRetryDialer(dialer net.Dialer, attempts int, delay time.Duration) *RetryDialer {
	if attempts <= 0 {
		attempts = defaultAttempts
	}
	if delay == 0 {
		delay = defaultDelay
	}

	return &RetryDialer{
		attempts: attempts,
		delay:    delay,
		Dialer:   dialer,
	}
}

// DialContext dials the specified address target. with transparent connection retries.
func (r *RetryDialer) DialContext(ctx context.Context, network string, addr string) (net.Conn, error) {
	var conn net.Conn
	var err error

	for i := 0; i < r.attempts; i++ {
		if i > 0 {
			select {
			case <-time.After(r.delay):
				r.delay *= 2
			case <-ctx.Done():
				return nil, ctx.Err()
			}

		}

		conn, err = r.Dialer.DialContext(ctx, network, addr)
		if err == nil {
			return conn, nil
		}

		zap.L().Warn(
			"dial encountered error; will retry",
			zap.String("network", network),
			zap.String("address", addr),
			zap.Int("attempt", i),
			zap.Error(err),
		)
	}

	return nil, err
}

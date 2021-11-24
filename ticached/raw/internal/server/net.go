package server

import (
	"crypto/tls"
	"crypto/x509"
	"net"
	"os"
	"sync"
	"time"
)

const (
	// certificateReloadInterval is the interval at which certificates are reloaded from disk.
	certificateReloadInterval = 30 * time.Minute
)

// Listener is an abstraction over net.Listen and tls.Listen that handles plaintext/TLS and
// certificate management.
type Listener struct {
	Network       string
	Address       string
	TLSKey        string
	TLSCert       string
	TLSCACert     string
	TLSClientAuth tls.ClientAuthType
}

// Listen creates a net.Listener with the current configuration.
func (l *Listener) Listen() (net.Listener, error) {
	if l.TLSKey == "" && l.TLSCert == "" {
		return net.Listen(l.Network, l.Address)
	}

	cs := &certificateStore{
		tlsKey:  l.TLSKey,
		tlsCert: l.TLSCert,
		tlsCA:   l.TLSCACert,
	}

	cfg := &tls.Config{
		GetCertificate: cs.getCertificate,
		ClientCAs:      cs.caPool,
		MinVersion:     tls.VersionTLS12,
		ClientAuth:     l.TLSClientAuth,
	}

	return tls.Listen(l.Network, l.Address, cfg)
}

// certificateStore abstracts the loading of TLS certificates with automatic reloading.
type certificateStore struct {
	tlsKey  string
	tlsCert string
	tlsCA   string

	keypair    tls.Certificate
	caPool     *x509.CertPool
	lastReload time.Time
	mu         sync.Mutex
}

// getCertificate implements the tls.Config callback for retrieving a server certificate to use for
// new connections.
func (c *certificateStore) getCertificate(hello *tls.ClientHelloInfo) (*tls.Certificate, error) {
	if err := c.reload(); err != nil {
		return nil, err
	}

	return &c.keypair, nil
}

// reload checks if the currently loaded certificates are stale, and if needed, reloads the key,
// certificate, and optionally CA certificate from disk while updating internal state.
func (c *certificateStore) reload() error {
	var err error

	c.mu.Lock()
	defer c.mu.Unlock()

	if time.Since(c.lastReload) < certificateReloadInterval {
		return nil
	}

	c.keypair, err = tls.LoadX509KeyPair(c.tlsCert, c.tlsKey)
	if err != nil {
		return err
	}

	c.caPool, err = x509.SystemCertPool()
	if err != nil {
		return err
	}

	if c.tlsCA != "" {
		ca, err := os.ReadFile(c.tlsCA)
		if err != nil {
			return err
		}

		c.caPool.AppendCertsFromPEM(ca)
	}

	return nil
}

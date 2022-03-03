package config

import (
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"os"
	"sync"
	"time"

	"go.uber.org/zap"

	"courier/internal/util"
)

// certificateStore is an abstraction over a set of TLS certificates that automatically manages
// periodic reloads.
type certificateStore struct {
	ctx        *TLSContext
	seedCAPool func() (*x509.CertPool, error)

	certs      []*certificate
	cas        *x509.CertPool
	lastReload time.Time
	mutex      sync.Mutex
}

// newCertificateStore creates a new certificate store.
func newCertificateStore(ctx *TLSContext, seedCAPool func() (*x509.CertPool, error)) (*certificateStore, error) {
	if seedCAPool == nil {
		seedCAPool = func() (*x509.CertPool, error) {
			return x509.NewCertPool(), nil
		}
	}

	cs := &certificateStore{
		ctx:        ctx,
		seedCAPool: seedCAPool,
	}

	if err := cs.reload(); err != nil {
		return nil, err
	}

	return cs, nil
}

// reload loads the certificates from disk and replaces the current store with the new certificates,
// in one atomic, concurrent-safe operation. If the reload fails, the old store's contents are
// retained without modification.
func (cs *certificateStore) reload() error {
	cs.mutex.Lock()
	defer cs.mutex.Unlock()

	var certs []*certificate

	cas, err := cs.seedCAPool()
	if err != nil {
		return &util.Error{
			Namespace:    "config",
			Message:      "failed to create seed CA pool",
			StackedError: err,
		}
	}

	for _, cert := range cs.ctx.Certificates {
		zap.L().Debug(
			"loading TLS certificate",
			zap.String("certificate", cert.Certificate),
			zap.String("key", cert.Key),
			zap.String("authority", cert.Authority),
			zap.Stringer("peer", cert.Peer),
		)

		keypair, err := newCertificate(cert)
		if err != nil {
			return &util.Error{
				Namespace: "config",
				Message:   "error reading X509 keypair from disk",
				Tags: map[string]interface{}{
					"cert": cert.Certificate,
					"key":  cert.Key,
				},
				StackedError: err,
			}
		}

		if cert.Authority != "" {
			ca, err := os.ReadFile(cert.Authority)
			if err != nil {
				return &util.Error{
					Namespace: "config",
					Message:   "error reading CA certificate from disk",
					Tags: map[string]interface{}{
						"authority": cert.Authority,
					},
					StackedError: err,
				}
			}

			cas.AppendCertsFromPEM(ca)
		}

		certs = append(certs, keypair)
	}

	cs.certs = certs
	cs.cas = cas
	cs.lastReload = time.Now()

	zap.L().Debug(
		"successfully completed TLS certificate reload routine",
		zap.Duration("interval", cs.ctx.ReloadInterval),
	)

	return nil
}

// checkStaleness checks if the certificate store is currently stale and initiates a reload if
// necessary. It returns an error indicating whether the reload was successful.
func (cs *certificateStore) checkStaleness() error {
	if cs.ctx.ReloadInterval <= 0 {
		return nil
	}

	if time.Since(cs.lastReload) > cs.ctx.ReloadInterval {
		zap.L().Debug(
			"reloading TLS certificates due to stale store",
			zap.Duration("age", time.Since(cs.lastReload)),
		)

		if err := cs.reload(); err != nil {
			zap.L().Error("failed to reload TLS certificates store", zap.Error(err))
			return err
		}
	}

	return nil
}

// getClientCertificate attempts to retrieve a client certificate from the store that is eligible
// for TLS origination based on the inbound tls.CertificateRequestInfo from the server. This routine
// will also statefully reload the certificate store if it is deemed stale by the configured reload
// interval.
func (cs *certificateStore) getClientCertificate(request *tls.CertificateRequestInfo) (*tls.Certificate, error) {
	cs.checkStaleness() // Ignore errors

	for _, cert := range cs.certs {
		if cert.supportsServer(request) {
			return &cert.keypair, nil
		}
	}

	return nil, &util.Error{
		Namespace: "config",
		Message:   "no client certificate available that satisfies handshake requirements",
	}
}

// getServerCertificate attempts to retrieve a certificate from the store that is eligible for TLS
// termination based on the inbound tls.ClientHelloInfo from the client. This routine will also
// statefully reload the certificate store if it is deemed stale by the configured reload interval.
func (cs *certificateStore) getServerCertificate(hello *tls.ClientHelloInfo) (*tls.Certificate, error) {
	cs.checkStaleness() // Ignore errors

	for _, cert := range cs.certs {
		if cert.supportsClient(hello) {
			return &cert.keypair, nil
		}
	}

	return nil, &util.Error{
		Namespace: "config",
		Message:   "no server certificate available that satisfies handshake requirements",
		Tags:      map[string]interface{}{"server_name": hello.ServerName},
	}
}

// certificate is an abstraction over a tls.Certificate that supports a peer regular expression
// pattern to determine eligibility for the certificate to be offered in a TLS handshake.
type certificate struct {
	keypair tls.Certificate
	pattern *Regex
}

// newCertificate loads the X509 keypair in the configuration specification from disk.
func newCertificate(cert *TLSCertificate) (*certificate, error) {
	keypair, err := tls.LoadX509KeyPair(cert.Certificate, cert.Key)
	if err != nil {
		return nil, &util.Error{
			Namespace:    "config",
			Message:      "failed to load certificate keypair",
			StackedError: err,
		}
	}

	keypair.Leaf, err = x509.ParseCertificate(keypair.Certificate[0])
	if err != nil {
		return nil, &util.Error{
			Namespace:    "config",
			Message:      "failed to parse certificate leaf",
			StackedError: err,
		}
	}

	return &certificate{
		keypair: keypair,
		pattern: cert.Peer,
	}, nil
}

// supportsClient determines whether this certificate is eligible for use in a TLS handshake against
// a client-supplied *tls.ClientHelloInfo.
// When no peer pattern is explicitly specified, this defers to Go's native certificate selection
// logic, which examines whether the client SNI (server name) matches a SAN in the certificate.
// When a peer pattern is explicitly specified, the certificate is deemed eligible as long as the
// pattern matches the client SNI.
func (c *certificate) supportsClient(hello *tls.ClientHelloInfo) bool {
	if c.pattern == nil {
		return hello.SupportsCertificate(&c.keypair) == nil
	}

	return c.pattern.MatchString(hello.ServerName)
}

// supportsServer determines whether this certificate is eligible for use in a TLS handshake against
// a server-supplied *tls.CertificateRequestInfo.
// When no peer pattern is explicitly specified, this defers to Go's native certificate selection
// logic, which offers the certificate if it signed by a CA the server advertises as acceptable.
// When a peer pattern is explicitly specified, the certificate is deemed eligible if the pattern
// matches the common name of any of the server's advertised acceptable CAs.
func (c *certificate) supportsServer(request *tls.CertificateRequestInfo) bool {
	if c.pattern == nil {
		return request.SupportsCertificate(&c.keypair) == nil
	}

	for _, ca := range request.AcceptableCAs {
		var name pkix.Name
		var seq pkix.RDNSequence

		if _, err := asn1.Unmarshal(ca, &seq); err != nil {
			return false
		}

		name.FillFromRDNSequence(&seq)

		if c.pattern.MatchString(name.CommonName) {
			return true
		}
	}

	return false
}

// peerVerifier provides functions for verifying both client and server connections against an
// authorized name pattern. Its implementation first defers Go's native x509 verification routine
// before checking peer certificate SANs against the specified regular expression.
// Reference: https://github.com/golang/go/blob/70deaa33ebd91944484526ab368fa19c499ff29f/src/crypto/tls/example_test.go#L186
type peerVerifier struct {
	pattern *Regex
	roots   *x509.CertPool
}

// verifyClient verifies a client connection.
func (pv *peerVerifier) verifyClient(state tls.ConnectionState) error {
	if len(state.PeerCertificates) < 2 {
		return &util.Error{
			Namespace: "config",
			Message:   "client did not supply a certificate",
		}
	}

	opts := x509.VerifyOptions{
		Intermediates: x509.NewCertPool(),
		Roots:         pv.roots,
		KeyUsages:     []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
	}

	for _, cert := range state.PeerCertificates[1:] {
		opts.Intermediates.AddCert(cert)
	}

	if _, err := state.PeerCertificates[0].Verify(opts); err != nil {
		return &util.Error{
			Namespace:    "config",
			Message:      "native client certificate verification failed",
			StackedError: err,
		}
	}

	for _, name := range state.PeerCertificates[0].DNSNames {
		if pv.pattern.MatchString(name) {
			return nil
		}
	}

	return &util.Error{
		Namespace: "config",
		Message:   "client certificate name verification error",
		Tags: map[string]interface{}{
			"presented": state.PeerCertificates[0].DNSNames,
			"pattern":   pv.pattern,
		},
	}
}

// verifyServer verifies a server connection.
func (pv *peerVerifier) verifyServer(state tls.ConnectionState) error {
	opts := x509.VerifyOptions{
		Intermediates: x509.NewCertPool(),
		Roots:         pv.roots,
	}

	for _, cert := range state.PeerCertificates[1:] {
		opts.Intermediates.AddCert(cert)
	}

	if _, err := state.PeerCertificates[0].Verify(opts); err != nil {
		return &util.Error{
			Namespace:    "config",
			Message:      "native server certificate verification failed",
			StackedError: err,
		}
	}

	for _, name := range state.PeerCertificates[0].DNSNames {
		if pv.pattern.MatchString(name) {
			return nil
		}
	}

	return &util.Error{
		Namespace: "config",
		Message:   "server certificate name verification error",
		Tags: map[string]interface{}{
			"presented": state.PeerCertificates[0].DNSNames,
			"pattern":   pv.pattern,
		},
	}
}

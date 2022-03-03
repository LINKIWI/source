package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"

	"unistore/internal/config"
	"unistore/internal/meta"
	"unistore/pkg/client/unistore"
)

// newClientConfig initializes the application client configuration.
func newClientConfig(path string, store string) (*config.Store, *config.RPC, error) {
	cfg, err := config.New(path)
	if err != nil {
		return nil, nil, err
	}

	for _, s := range cfg.Client.Stores {
		if s.Name == store {
			return &s, &cfg.Client.RPC, nil
		}
	}

	return nil, nil, fmt.Errorf(
		"no such store defined in configuration: path=%s store=%s",
		path,
		store,
	)
}

// newUnistoreClient initializes a Unistore gRPC client based on the client configuration.
func newUnistoreClient(store *config.Store, rpc *config.RPC) (*unistore.Client, error) {
	var cancel context.CancelFunc
	var opts []grpc.DialOption

	instance := "cli"
	if store.Connection.Identity != "" {
		instance = store.Connection.Identity
	}
	clientID := fmt.Sprintf(
		"unistore-cli/%s (instance:%s; authority:%s)",
		meta.Version,
		instance,
		store.Authority,
	)

	opts = append(opts, grpc.WithAuthority(store.Authority))
	opts = append(opts, grpc.WithUserAgent(clientID))
	opts = append(opts, grpc.WithUnaryInterceptor(newRequestTimeoutUnaryInterceptor(store.Connection.RequestTimeout)))

	if rpc.MaxRecvMessageSize > 0 {
		opts = append(opts, grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(rpc.MaxRecvMessageSize)))
	}
	if rpc.MaxSendMessageSize > 0 {
		opts = append(opts, grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(rpc.MaxSendMessageSize)))
	}

	if tlsCtx := store.Connection.TLSContext; tlsCtx != nil {
		tlsCfg, err := newTLSConfig(tlsCtx)
		if err != nil {
			return nil, err
		}

		tlsCfg.ServerName = store.Authority
		opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(tlsCfg)))
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	ctx := context.Background()
	if store.Connection.ConnectTimeout > 0 {
		ctx, cancel = context.WithTimeout(ctx, store.Connection.ConnectTimeout)
		defer cancel()
	}

	return unistore.New(ctx, store.Address, opts...)
}

// newTLSConfig creates a new *tls.Config from the configuration-declared TLS context.
func newTLSConfig(cfg *config.TLSContext) (*tls.Config, error) {
	var certificates []tls.Certificate

	caPool, err := x509.SystemCertPool()
	if err != nil {
		return nil, fmt.Errorf(
			"error retrieving system CA certificate pool: err=%v",
			err,
		)
	}

	for _, cert := range cfg.Certificates {
		keypair, err := tls.LoadX509KeyPair(cert.Certificate, cert.Key)
		if err != nil {
			return nil, fmt.Errorf(
				"error loading client keypair: key=%s cert=%s err=%v",
				cert.Key,
				cert.Certificate,
				err,
			)
		}

		if cert.Authority != "" {
			ca, err := os.ReadFile(cert.Authority)
			if err != nil {
				return nil, fmt.Errorf(
					"error reading CA certificate: cert=%s err=%v",
					cert.Authority,
					err,
				)
			}

			caPool.AppendCertsFromPEM(ca)
		}

		certificates = append(certificates, keypair)
	}

	tlsCfg := &tls.Config{
		Certificates: certificates,
		RootCAs:      caPool,
	}

	if cfg.VerifyPeer != "" {
		tlsCfg.InsecureSkipVerify = true
		tlsCfg.VerifyConnection = func(state tls.ConnectionState) error {
			verifyOpts := x509.VerifyOptions{
				Intermediates: x509.NewCertPool(),
				Roots:         caPool,
			}

			for _, cert := range state.PeerCertificates[1:] {
				verifyOpts.Intermediates.AddCert(cert)
			}

			if _, err := state.PeerCertificates[0].Verify(verifyOpts); err != nil {
				return fmt.Errorf(
					"native server certificate verification failed: err=%v",
					err,
				)
			}

			return state.PeerCertificates[0].VerifyHostname(cfg.VerifyPeer)
		}
	}

	return tlsCfg, nil
}

// newRequestTimeoutUnaryInterceptor creates a new grpc.UnaryClientInterceptor that automatically
// applies a request timeout on all unary requests by augmenting the parent context passed to the
// invoker.
func newRequestTimeoutUnaryInterceptor(timeout time.Duration) grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req interface{},
		reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		var cancel context.CancelFunc

		if timeout > 0 {
			ctx, cancel = context.WithTimeout(ctx, timeout)
			defer cancel()
		}

		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

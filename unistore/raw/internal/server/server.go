package server

import (
	"context"
	"fmt"
	"net"

	"go.uber.org/zap"
	"golang.org/x/net/netutil"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"unistore/internal/config"
	"unistore/internal/server/middleware"
	"unistore/schemas/service"
)

// Unistore wraps a grpc.Server to provide the core Unistore gRPC services.
type Unistore struct {
	cfg *config.Server
	*grpc.Server
}

// New creates a new Unistore gRPC server.
func New(cfg *config.Server) (*Unistore, error) {
	var opts []grpc.ServerOption

	if cfg.RPC.MaxRecvMessageSize > 0 {
		opts = append(opts, grpc.MaxRecvMsgSize(cfg.RPC.MaxRecvMessageSize))
	}
	if cfg.RPC.MaxSendMessageSize > 0 {
		opts = append(opts, grpc.MaxSendMsgSize(cfg.RPC.MaxSendMessageSize))
	}
	if cfg.RPC.MaxConcurrentStreams > 0 {
		opts = append(opts, grpc.MaxConcurrentStreams(uint32(cfg.RPC.MaxConcurrentStreams)))
	}

	opts = append(
		opts,
		grpc.ChainUnaryInterceptor(
			middleware.NewErrorCaptureUnaryInterceptor(),
			middleware.NewLogUnaryInterceptor(),
		),
		grpc.ChainStreamInterceptor(
			middleware.NewErrorCaptureStreamInterceptor(),
			middleware.NewLogStreamInterceptor(),
		),
	)

	srv := grpc.NewServer(opts...)

	unistoreSrv, err := newUnistoreService(cfg)
	if err != nil {
		return nil, fmt.Errorf("server: error initializing Unistore service: err=%v", err)
	}

	metaSrv, err := newMetaService(cfg, unistoreSrv.(*unistoreService).backend)
	if err != nil {
		return nil, fmt.Errorf("server: error initializing meta service: err=%v", err)
	}

	service.RegisterUnistoreServer(srv, unistoreSrv)
	service.RegisterMetaServer(srv, metaSrv)
	reflection.Register(srv)

	for name, info := range srv.GetServiceInfo() {
		zap.L().Debug(
			"initialized gRPC service",
			zap.String("service", name),
			zap.Any("methods", info.Methods),
		)
	}

	return &Unistore{
		cfg:    cfg,
		Server: srv,
	}, nil
}

// Serve serves the gRPC server over the specified listener. The context is used for graceful
// shutdown of the server when canceled.
func (u *Unistore) Serve(ctx context.Context) error {
	ln, err := net.Listen(u.cfg.Listener.Address.Network(), u.cfg.Listener.Address.String())
	if err != nil {
		return fmt.Errorf("server: error creating listener: err=%v", err)
	}

	if u.cfg.Listener.MaxActiveConnections > 0 {
		ln = netutil.LimitListener(ln, u.cfg.Listener.MaxActiveConnections)
	}

	go func() {
		<-ctx.Done()

		if err := u.Close(); err != nil {
			zap.L().Error("error during server close", zap.Error(err))
		}
	}()

	zap.L().Info(
		"serving gRPC indefinitely",
		zap.String("network", ln.Addr().Network()),
		zap.Stringer("address", ln.Addr()),
	)

	return u.Server.Serve(ln)
}

// Close closes the server by stopping the backing gRPC server gracefully.
func (u *Unistore) Close() error {
	u.Server.GracefulStop()

	return nil
}

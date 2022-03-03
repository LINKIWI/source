package server

import (
	"context"

	"unistore/internal/backend"
	"unistore/internal/config"
	"unistore/internal/meta"
	"unistore/schemas/service"
)

// metaService is an implementation of the service.MetaServer gRPC service.
type metaService struct {
	cfg     *config.Server
	backend backend.Backend
	service.UnimplementedMetaServer
}

// newMetaService creates a new meta gRPC service.
func newMetaService(cfg *config.Server, backend backend.Backend) (service.MetaServer, error) {
	return &metaService{
		cfg:     cfg,
		backend: backend,
	}, nil
}

// HealthCheck is a noop and always returns a successful, empty response.
func (m *metaService) HealthCheck(ctx context.Context, request *service.HealthCheckRequest) (*service.HealthCheckResponse, error) {
	return &service.HealthCheckResponse{}, nil
}

// Info returns static information about the server.
func (m *metaService) Info(ctx context.Context, request *service.InfoRequest) (*service.InfoResponse, error) {
	return &service.InfoResponse{
		Version: meta.Version,
		Address: m.cfg.Listener.Address.Spec(),
		Backend: m.backend.Descriptor(),
	}, nil
}

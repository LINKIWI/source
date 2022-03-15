package server

import (
	"context"
	"strings"

	"go.uber.org/zap"

	"unistore/internal/backend"
	"unistore/internal/config"
	"unistore/schemas/common"
	"unistore/schemas/service"
)

// unistoreService is an implementation of the service.UnistoreServer gRPC service.
type unistoreService struct {
	backend backend.Backend
	service.UnimplementedUnistoreServer
}

// newUnistoreService creates a new Unistore gRPC service from storage configuration.
func newUnistoreService(cfg *config.Server) (*unistoreService, error) {
	var err error
	var tree backend.Backend

	mux := make(map[common.Backend]backend.Backend)

	if cfg.Storage.Disk != nil {
		disk, err := backend.NewDisk(cfg.Storage.Disk)
		if err != nil {
			return nil, err
		}

		if cfg.Storage.Disk.Buffer != nil {
			disk = backend.NewBuffer(
				cfg.Storage.Disk.Buffer.GetChunkSize,
				cfg.Storage.Disk.Buffer.PutChunkSize,
				disk,
			)
		}

		if cfg.Storage.Disk.Encryption != nil {
			disk = backend.NewEncryption(
				cfg.Storage.Disk.Encryption.Mechanism,
				cfg.Storage.Disk.Encryption.PrivateKey,
				cfg.Storage.Disk.Encryption.PublicKey,
				disk,
			)
			if err != nil {
				return nil, err
			}
		}

		if cfg.Storage.Disk.Compression != nil {
			disk = backend.NewCompression(
				cfg.Storage.Disk.Compression.Algorithm,
				cfg.Storage.Disk.Compression.Level,
				disk,
			)
		}

		if cfg.Storage.Disk.Checksum != nil {
			disk = backend.NewChecksum(cfg.Storage.Disk.Checksum.Algorithm, disk)
		}

		if len(cfg.Storage.Disk.BucketPermissions) > 0 {
			disk = backend.NewPermission(cfg.Storage.Disk.BucketPermissions, disk)
		}

		if len(cfg.Storage.Disk.BucketAliases) > 0 {
			disk = backend.NewAlias(cfg.Storage.Disk.BucketAliases, disk)
		}

		disk, err = backend.NewMetadata(
			cfg.Storage.Disk.Compression,
			cfg.Storage.Disk.Encryption,
			disk,
		)
		if err != nil {
			return nil, err
		}

		if cfg.Storage.Disk.Log != nil {
			tags := make(map[string]string)
			for _, tag := range cfg.Storage.Disk.Log.Tags {
				tags[tag.Key] = tag.Value
			}

			disk, err = backend.NewLog(
				cfg.Storage.Disk.Log.Level,
				cfg.Storage.Disk.Log.Outputs,
				tags,
				disk,
			)
			if err != nil {
				return nil, err
			}
		}

		zap.L().Info(
			"configured storage backend",
			zap.Stringer("type", common.Backend_DISK),
			zap.Stringer("backend", disk),
		)

		mux[common.Backend_DISK] = disk
	}

	if cfg.Storage.Unistore != nil {
		unistore, err := backend.NewUnistore(cfg.Storage.Unistore, cfg.RPC)
		if err != nil {
			return nil, err
		}

		if cfg.Storage.Unistore.Buffer != nil {
			unistore = backend.NewBuffer(
				cfg.Storage.Unistore.Buffer.GetChunkSize,
				cfg.Storage.Unistore.Buffer.PutChunkSize,
				unistore,
			)
		}

		if cfg.Storage.Unistore.Threshold != nil {
			unistore = backend.NewThreshold(
				cfg.Storage.Unistore.Threshold.MinStreamPutSize,
				cfg.Storage.Unistore.Threshold.MaxUnaryPutSize,
				unistore,
			)
		}

		if cfg.Storage.Unistore.Encryption != nil {
			unistore = backend.NewEncryption(
				cfg.Storage.Unistore.Encryption.Mechanism,
				cfg.Storage.Unistore.Encryption.PrivateKey,
				cfg.Storage.Unistore.Encryption.PublicKey,
				unistore,
			)
		}

		if cfg.Storage.Unistore.Compression != nil {
			unistore = backend.NewCompression(
				cfg.Storage.Unistore.Compression.Algorithm,
				cfg.Storage.Unistore.Compression.Level,
				unistore,
			)
		}

		if cfg.Storage.Unistore.Checksum != nil {
			unistore = backend.NewChecksum(
				cfg.Storage.Unistore.Checksum.Algorithm,
				unistore,
			)
		}

		if len(cfg.Storage.Unistore.BucketPermissions) > 0 {
			unistore = backend.NewPermission(
				cfg.Storage.Unistore.BucketPermissions,
				unistore,
			)
		}

		if len(cfg.Storage.Unistore.BucketAliases) > 0 {
			unistore = backend.NewAlias(cfg.Storage.Unistore.BucketAliases, unistore)
		}

		unistore, err = backend.NewMetadata(
			cfg.Storage.Unistore.Compression,
			cfg.Storage.Unistore.Encryption,
			unistore,
		)
		if err != nil {
			return nil, err
		}

		if cfg.Storage.Unistore.Log != nil {
			tags := make(map[string]string)
			for _, tag := range cfg.Storage.Unistore.Log.Tags {
				tags[tag.Key] = tag.Value
			}

			unistore, err = backend.NewLog(
				cfg.Storage.Unistore.Log.Level,
				cfg.Storage.Unistore.Log.Outputs,
				tags,
				unistore,
			)
			if err != nil {
				return nil, err
			}
		}

		zap.L().Info(
			"configured storage backend",
			zap.Stringer("type", common.Backend_UNISTORE),
			zap.Stringer("backend", unistore),
		)

		mux[common.Backend_UNISTORE] = unistore
	}

	if cfg.Storage.B2 != nil {
		b2, err := backend.NewB2(cfg.Storage.B2)
		if err != nil {
			return nil, err
		}

		if cfg.Storage.B2.Buffer != nil {
			b2 = backend.NewBuffer(
				cfg.Storage.B2.Buffer.GetChunkSize,
				cfg.Storage.B2.Buffer.PutChunkSize,
				b2,
			)
		}

		if cfg.Storage.B2.Threshold != nil {
			b2 = backend.NewThreshold(
				cfg.Storage.B2.Threshold.MinStreamPutSize,
				cfg.Storage.B2.Threshold.MaxUnaryPutSize,
				b2,
			)
		}

		if cfg.Storage.B2.Encryption != nil {
			b2 = backend.NewEncryption(
				cfg.Storage.B2.Encryption.Mechanism,
				cfg.Storage.B2.Encryption.PrivateKey,
				cfg.Storage.B2.Encryption.PublicKey,
				b2,
			)
		}

		if cfg.Storage.B2.Compression != nil {
			b2 = backend.NewCompression(
				cfg.Storage.B2.Compression.Algorithm,
				cfg.Storage.B2.Compression.Level,
				b2,
			)
		}

		if cfg.Storage.B2.Checksum != nil {
			b2 = backend.NewChecksum(cfg.Storage.B2.Checksum.Algorithm, b2)
		}

		if len(cfg.Storage.B2.BucketPermissions) > 0 {
			b2 = backend.NewPermission(cfg.Storage.B2.BucketPermissions, b2)
		}

		if len(cfg.Storage.B2.BucketAliases) > 0 {
			b2 = backend.NewAlias(cfg.Storage.B2.BucketAliases, b2)
		}

		b2, err = backend.NewMetadata(
			cfg.Storage.B2.Compression,
			cfg.Storage.B2.Encryption,
			b2,
		)
		if err != nil {
			return nil, err
		}

		if cfg.Storage.B2.Log != nil {
			tags := make(map[string]string)
			for _, tag := range cfg.Storage.B2.Log.Tags {
				tags[tag.Key] = tag.Value
			}

			b2, err = backend.NewLog(
				cfg.Storage.B2.Log.Level,
				cfg.Storage.B2.Log.Outputs,
				tags,
				b2,
			)
			if err != nil {
				return nil, err
			}
		}

		zap.L().Info(
			"configured storage backend",
			zap.Stringer("type", common.Backend_B2),
			zap.Stringer("backend", b2),
		)

		mux[common.Backend_B2] = b2
	}

	if cfg.Storage.Composite != nil {
		var children []backend.Backend

		for _, name := range cfg.Storage.Composite.Backends {
			bid := common.Backend(common.Backend_value[strings.ToUpper(name)])
			children = append(children, mux[bid])
		}

		composite := backend.NewComposite(
			cfg.Storage.Composite.ReadDispatch,
			cfg.Storage.Composite.WriteDispatch,
			children,
		)

		if cfg.Storage.Composite.Log != nil {
			tags := make(map[string]string)
			for _, tag := range cfg.Storage.Composite.Log.Tags {
				tags[tag.Key] = tag.Value
			}

			composite, err = backend.NewLog(
				cfg.Storage.Composite.Log.Level,
				cfg.Storage.Composite.Log.Outputs,
				tags,
				composite,
			)
			if err != nil {
				return nil, err
			}
		}

		zap.L().Info(
			"configured storage backend",
			zap.Stringer("type", common.Backend_COMPOSITE),
			zap.Stringer("backend", composite),
		)

		mux[common.Backend_COMPOSITE] = composite
	}

	tree = backend.NewMux(mux)
	tree = backend.NewValidation(tree)

	return &unistoreService{backend: tree}, nil
}

// HeadBucket calls into the backend's HeadBucket.
func (u *unistoreService) HeadBucket(ctx context.Context, request *service.HeadBucketRequest) (*service.HeadBucketResponse, error) {
	return u.backend.HeadBucket(ctx, request)
}

// HeadObject calls into the backend's HeadObject.
func (u *unistoreService) HeadObject(ctx context.Context, request *service.HeadObjectRequest) (*service.HeadObjectResponse, error) {
	return u.backend.HeadObject(ctx, request)
}

// GetObject calls into the backend's GetObject.
func (u *unistoreService) GetObject(ctx context.Context, request *service.GetObjectRequest) (*service.GetObjectResponse, error) {
	return u.backend.GetObject(ctx, request)
}

// StreamGetObject calls into the backend's StreamGetObject.
func (u *unistoreService) StreamGetObject(request *service.GetObjectStreamRequest, stream service.Unistore_StreamGetObjectServer) error {
	responses, errs := u.backend.StreamGetObject(stream.Context(), request)
	if responses == nil {
		return <-errs
	}

	for {
		select {
		case err := <-errs:
			if err != nil {
				return err
			}
		case response, ok := <-responses:
			if !ok {
				// Channel closed; there are no remaining responses to read
				return nil
			}

			if err := stream.Send(response); err != nil {
				return err
			}
		}
	}
}

// PutObject calls into the backend's PutObject.
func (u *unistoreService) PutObject(ctx context.Context, request *service.PutObjectRequest) (*service.PutObjectResponse, error) {
	return u.backend.PutObject(ctx, request)
}

// StreamPutObject calls into the backend's StreamPutObject.
func (u *unistoreService) StreamPutObject(stream service.Unistore_StreamPutObjectServer) error {
	puts := make(chan *service.PutObjectStreamRequest)

	go func() {
		for {
			request, err := stream.Recv()
			if err != nil {
				break
			}

			puts <- request
		}

		close(puts)
	}()

	response, err := u.backend.StreamPutObject(stream.Context(), puts)
	if err != nil {
		return err
	}

	return stream.SendAndClose(response)
}

// DeleteObject calls into the backend's DeleteObject.
func (u *unistoreService) DeleteObject(ctx context.Context, request *service.DeleteObjectRequest) (*service.DeleteObjectResponse, error) {
	return u.backend.DeleteObject(ctx, request)
}

// ListBuckets calls into the backend's ListBuckets.
func (u *unistoreService) ListBuckets(ctx context.Context, request *service.ListBucketsRequest) (*service.ListBucketsResponse, error) {
	return u.backend.ListBuckets(ctx, request)
}

// ListObjects calls into the backend's ListObjects.
func (u *unistoreService) ListObjects(ctx context.Context, request *service.ListObjectsRequest) (*service.ListObjectsResponse, error) {
	return u.backend.ListObjects(ctx, request)
}

// Close closes the entire backend tree.
func (u *unistoreService) Close() error {
	return u.backend.Close()
}

// String provides a human-readable representation of the service by borrowing the string
// representation of the associated backend tree.
func (u *unistoreService) String() string {
	return u.backend.String()
}

package backend

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"unistore/internal/config"
	"unistore/internal/schemas"
	"unistore/pkg/client/disk"
	"unistore/schemas/common"
	"unistore/schemas/service"
)

// Disk is a Backend that implements an object storage server on top of a local filesystem.
type Disk struct {
	cfg    *config.Disk
	client *disk.Client
}

// NewDisk creates a new Disk backend.
func NewDisk(cfg *config.Disk) (Backend, error) {
	client, err := disk.New(cfg.Root)
	if err != nil {
		return nil, err
	}

	d := &Disk{
		cfg:    cfg,
		client: client,
	}

	return newInstrumentation("disk", d), nil
}

// HeadBucket stats the requested bucket directory.
func (d *Disk) HeadBucket(ctx context.Context, request *service.HeadBucketRequest) (*service.HeadBucketResponse, error) {
	stat, err := d.client.StatFile(request.Resource.Bucket)
	if err != nil {
		return nil, d.createError(err)
	}

	return &service.HeadBucketResponse{
		Bucket: &common.Bucket{
			Id:        stat.BaseName,
			Name:      stat.BaseName,
			Timestamp: timestamppb.New(stat.ModifiedTime),
			Tags:      stat.Attributes,
		},
	}, nil
}

// HeadObject stats the requested file.
func (d *Disk) HeadObject(ctx context.Context, request *service.HeadObjectRequest) (*service.HeadObjectResponse, error) {
	var objectType common.Object

	target, err := newFileTarget(d.client, request.Resource.Bucket, request.Key)
	if err != nil {
		return nil, status.Error(codes.PermissionDenied, err.Error())
	}

	stat, err := d.client.StatFile(target.resolve())
	if err != nil {
		return nil, d.createError(err)
	}

	attributes, err := schemas.UnmarshalAttributes(
		stat.Attributes,
		&schemas.AttributesEncodingOptions{Prefix: "user"},
	)
	if err != nil {
		return nil, d.createError(err)
	}

	if stat.Mode.IsRegular() {
		objectType = common.Object_REGULAR
	} else if stat.Mode.IsDir() {
		objectType = common.Object_DIRECTORY
	}

	return &service.HeadObjectResponse{
		Metadata: &common.Metadata{
			Id:           stat.BaseName,
			Key:          target.key,
			Size:         stat.Size,
			ObjectType:   objectType,
			ModifiedTime: timestamppb.New(stat.ModifiedTime),
			Attributes:   attributes,
		},
	}, nil
}

// GetObject opens a file handle for the requested file followed by reading its contents in full.
func (d *Disk) GetObject(ctx context.Context, request *service.GetObjectRequest) (*service.GetObjectResponse, error) {
	target, err := newFileTarget(d.client, request.Resource.Bucket, request.Key)
	if err != nil {
		return nil, status.Error(codes.PermissionDenied, err.Error())
	}

	file, stat, err := d.client.ReadFile(target.resolve())
	if err != nil {
		return nil, d.createError(err)
	}

	defer file.Close()

	if !stat.Mode.IsRegular() {
		return nil, status.Errorf(
			codes.InvalidArgument,
			"backend: cannot serve non-regular file: file=%s mode=%v",
			stat.Path,
			stat.Mode,
		)
	}

	n := stat.Size

	if request.Range != nil {
		if request.Range.Unit != "bytes" {
			return nil, status.Errorf(
				codes.InvalidArgument,
				"backend: unsupported range unit: unit=%s",
				request.Range.Unit,
			)
		} else if request.Range.End > 0 && request.Range.End <= request.Range.Start {
			return nil, status.Errorf(
				codes.OutOfRange,
				"backend: invalid range specification: start=%d end=%d",
				request.Range.Start,
				request.Range.End,
			)
		} else if request.Range.Start >= stat.Size {
			return nil, status.Errorf(
				codes.OutOfRange,
				"backend: range start index is beyond file boundary: start=%d size=%d",
				request.Range.Start,
				stat.Size,
			)
		}

		if _, err := file.Seek(int64(request.Range.Start), io.SeekStart); err != nil {
			return nil, d.createError(err)
		}

		// Read only the number of bytes requested by the range specification, up to the
		// end of the file. A zero end index implicitly reads to the end of the file.
		end := request.Range.End
		if end == 0 || end > stat.Size {
			end = stat.Size
		}
		n = end - request.Range.Start
	}

	data := make([]byte, n)
	if _, err := io.ReadFull(file, data); err != nil {
		return nil, d.createError(err)
	}

	attributes, err := schemas.UnmarshalAttributes(
		stat.Attributes,
		&schemas.AttributesEncodingOptions{Prefix: "user"},
	)
	if err != nil {
		return nil, d.createError(err)
	}

	return &service.GetObjectResponse{
		Data: data,
		Metadata: &common.Metadata{
			Id:           stat.BaseName,
			Key:          target.key,
			Size:         stat.Size,
			ObjectType:   common.Object_REGULAR,
			ModifiedTime: timestamppb.New(stat.ModifiedTime),
			Attributes:   attributes,
		},
	}, nil
}

// StreamGetObject reads the requested file in chunks, each of which is packaged as a single
// response instance over the return channel.
func (d *Disk) StreamGetObject(ctx context.Context, request *service.GetObjectStreamRequest) (<-chan *service.GetObjectStreamResponse, chan error) {
	stream := make(chan *service.GetObjectStreamResponse)
	errs := make(chan error, 1)

	if request.ChunkSize == 0 {
		errs <- status.Error(
			codes.InvalidArgument,
			"backend: disk streaming get chunk size must be greater than zero",
		)
		return nil, errs
	}

	target, err := newFileTarget(d.client, request.Request.Resource.Bucket, request.Request.Key)
	if err != nil {
		errs <- status.Error(codes.PermissionDenied, err.Error())
		return nil, errs
	}

	file, stat, err := d.client.ReadFile(target.resolve())
	if err != nil {
		errs <- d.createError(err)
		return nil, errs
	}

	if !stat.Mode.IsRegular() {
		file.Close()
		errs <- status.Errorf(
			codes.InvalidArgument,
			"backend: cannot serve non-regular file: file=%s mode=%v",
			stat.Path,
			stat.Mode,
		)
		return nil, errs
	}

	attributes, err := schemas.UnmarshalAttributes(
		stat.Attributes,
		&schemas.AttributesEncodingOptions{Prefix: "user"},
	)
	if err != nil {
		file.Close()
		errs <- d.createError(err)
		return nil, errs
	}

	if request.Request.Range != nil {
		if request.Request.Range.Unit != "bytes" {
			file.Close()
			errs <- status.Errorf(
				codes.InvalidArgument,
				"backend: unsupported range unit: unit=%s",
				request.Request.Range.Unit,
			)
			return nil, errs
		} else if request.Request.Range.End > 0 && request.Request.Range.End <= request.Request.Range.Start {
			file.Close()
			errs <- status.Errorf(
				codes.OutOfRange,
				"backend: invalid range specification: start=%d end=%d",
				request.Request.Range.Start,
				request.Request.Range.End,
			)
			return nil, errs
		} else if request.Request.Range.Start >= stat.Size {
			file.Close()
			errs <- status.Errorf(
				codes.OutOfRange,
				"backend: range start index is beyond file boundary: start=%d size=%d",
				request.Request.Range.Start,
				stat.Size,
			)
			return nil, errs
		}
	}

	go func() {
		defer file.Close()

		offset := uint64(0)
		end := stat.Size
		total := stat.Size

		if request.Request.Range != nil {
			offset = request.Request.Range.Start
			end = request.Request.Range.End
			if end == 0 || end > stat.Size {
				end = stat.Size
			}
			total = end - request.Request.Range.Start

			if _, err := file.Seek(int64(offset), io.SeekStart); err != nil {
				errs <- d.createError(err)
				return
			}
		}

		// Wrap the file with a buffered reader in order to normalize the number of actual
		// underlying read(2) syscalls on the file descriptor, regardless of the chunk size
		// requested by the client.
		reader := bufio.NewReader(file)
		metadata := &common.Metadata{
			Id:           stat.BaseName,
			Key:          target.key,
			Size:         stat.Size,
			ObjectType:   common.Object_REGULAR,
			ModifiedTime: timestamppb.New(stat.ModifiedTime),
			Attributes:   attributes,
		}

		for {
			n := request.ChunkSize
			if n > end-offset {
				n = end - offset
			}

			chunk := make([]byte, n)

			if _, err := io.ReadFull(reader, chunk); err != nil {
				errs <- d.createError(err)
				break
			}

			response := &service.GetObjectStreamResponse{
				Response: &service.GetObjectResponse{
					Metadata: metadata,
					Data:     chunk,
				},
				Range: &common.Range{
					Unit:  "bytes",
					Start: offset,
					End:   offset + n,
					Total: total,
				},
			}

			stream <- response
			offset += n

			if offset >= end {
				break
			}
		}

		close(stream)
	}()

	return stream, errs
}

// PutObject opens a file handle at the requested path followed by writing the request data payload
// to the file in full.
func (d *Disk) PutObject(ctx context.Context, request *service.PutObjectRequest) (*service.PutObjectResponse, error) {
	target, err := newFileTarget(d.client, request.Resource.Bucket, request.Key)
	if err != nil {
		return nil, status.Error(codes.PermissionDenied, err.Error())
	}

	if err := d.client.CreateDirectory(target.parent()); err != nil {
		return nil, d.createError(err)
	}

	w, err := d.client.WriteFile(target.resolve())
	if err != nil {
		return nil, d.createError(err)
	}

	if _, err := io.Copy(w, bytes.NewReader(request.Data)); err != nil {
		w.Close()
		return nil, d.createError(err)
	}

	w.Close()

	attributes, err := schemas.MarshalAttributes(
		request.Attributes,
		&schemas.AttributesEncodingOptions{Prefix: "user"},
	)
	if err != nil {
		return nil, d.createError(err)
	}

	if err := d.client.WriteAttributes(target.resolve(), attributes); err != nil {
		return nil, d.createError(err)
	}

	stat, err := d.client.StatFile(target.resolve())
	if err != nil {
		return nil, d.createError(err)
	}

	return &service.PutObjectResponse{
		Metadata: &common.Metadata{
			Id:           stat.BaseName,
			Key:          target.key,
			Size:         stat.Size,
			ObjectType:   common.Object_REGULAR,
			ModifiedTime: timestamppb.New(stat.ModifiedTime),
			Attributes:   request.Attributes,
		},
	}, nil
}

// StreamPutObject writes the request file in chunks, where each request over the request channel
// specifies a single chunk of data.
func (d *Disk) StreamPutObject(ctx context.Context, stream chan *service.PutObjectStreamRequest) (*service.PutObjectStreamResponse, error) {
	var err error
	var target *fileTarget
	var file disk.WriteSeekCloser
	var request *service.PutObjectStreamRequest

	for request = range stream {
		// Initialize the file descriptor for only the first request in the channel
		if file == nil {
			target, err = newFileTarget(d.client, request.Request.Resource.Bucket, request.Request.Key)
			if err != nil {
				return nil, status.Error(codes.PermissionDenied, err.Error())
			}

			if err := d.client.CreateDirectory(target.parent()); err != nil {
				return nil, d.createError(err)
			}

			file, err = d.client.WriteFile(target.resolve())
			if err != nil {
				return nil, d.createError(err)
			}

			defer file.Close()
		}

		if request.Range != nil && request.Range.Start > 0 && request.Range.End > request.Range.Start {
			if uint64(len(request.Request.Data)) != request.Range.End-request.Range.Start {
				return nil, status.Error(
					codes.InvalidArgument,
					"backend: write payload size is inconsistent with range specification",
				)
			}

			if _, err := file.Seek(int64(request.Range.Start), io.SeekStart); err != nil {
				return nil, status.Error(codes.OutOfRange, err.Error())
			}
		}

		if _, err := file.Write(request.Request.Data); err != nil {
			return nil, d.createError(err)
		}
	}

	// Use attributes from the latest (most recent) request from the stream.
	attributes, err := schemas.MarshalAttributes(
		request.Request.Attributes,
		&schemas.AttributesEncodingOptions{Prefix: "user"},
	)
	if err != nil {
		return nil, d.createError(err)
	}

	if err := d.client.WriteAttributes(target.resolve(), attributes); err != nil {
		return nil, d.createError(err)
	}

	stat, err := d.client.StatFile(target.resolve())
	if err != nil {
		return nil, d.createError(err)
	}

	return &service.PutObjectStreamResponse{
		Response: &service.PutObjectResponse{
			Metadata: &common.Metadata{
				Id:           stat.BaseName,
				Key:          target.key,
				Size:         stat.Size,
				ObjectType:   common.Object_REGULAR,
				ModifiedTime: timestamppb.New(stat.ModifiedTime),
				Attributes:   request.Request.Attributes,
			},
		},
	}, nil
}

// DeleteObject deletes the requested file.
func (d *Disk) DeleteObject(ctx context.Context, request *service.DeleteObjectRequest) (*service.DeleteObjectResponse, error) {
	target, err := newFileTarget(d.client, request.Resource.Bucket, request.Key)
	if err != nil {
		return nil, status.Error(codes.PermissionDenied, err.Error())
	}

	stat, err := d.client.StatFile(target.resolve())
	if err != nil {
		return nil, d.createError(err)
	}

	if !request.Recursive || stat.Mode.IsRegular() {
		if err := d.client.DeleteFile(target.resolve()); err != nil {
			return nil, d.createError(err)
		}

		return &service.DeleteObjectResponse{Deleted: 1}, nil
	}

	// Recursive semantics apply only to directory object types.

	files, err := d.client.ListFiles(target.resolve(), true)
	if err != nil {
		return nil, d.createError(err)
	}

	for _, file := range files {
		if err := d.client.DeleteFile(file.Path); err != nil {
			return nil, d.createError(err)
		}
	}

	return &service.DeleteObjectResponse{Deleted: int32(len(files))}, nil
}

// ListBuckets lists all bucket directories.
func (d *Disk) ListBuckets(ctx context.Context, request *service.ListBucketsRequest) (*service.ListBucketsResponse, error) {
	var buckets []*service.HeadBucketResponse

	files, err := d.client.ListFiles("", false)
	if err != nil {
		return nil, d.createError(err)
	}

	for _, file := range files {
		if file.Mode.IsDir() {
			head := &service.HeadBucketResponse{
				Bucket: &common.Bucket{
					Id:        file.BaseName,
					Name:      file.BaseName,
					Timestamp: timestamppb.New(file.ModifiedTime),
					Tags:      file.Attributes,
				},
			}

			buckets = append(buckets, head)
		}
	}

	return &service.ListBucketsResponse{Buckets: buckets}, nil
}

// ListObjects lists files matching the requested prefix.
func (d *Disk) ListObjects(ctx context.Context, request *service.ListObjectsRequest) (*service.ListObjectsResponse, error) {
	var objects []*service.HeadObjectResponse

	target, err := newFileTarget(d.client, request.Resource.Bucket, request.Prefix)
	if err != nil {
		return nil, status.Error(codes.PermissionDenied, err.Error())
	}

	files, err := d.client.ListFiles(target.parent(), request.Recursive)
	if err != nil {
		return nil, d.createError(err)
	}

	for _, file := range files {
		var objectType common.Object

		// Exclude the bucket name in the file path key for purposes of prefix evaluation
		key := filepath.Join(strings.Split(file.Path, string(os.PathSeparator))[1:]...)

		if target.key != "." && !strings.HasPrefix(key, target.key) {
			continue
		}

		attributes, err := schemas.UnmarshalAttributes(
			file.Attributes,
			&schemas.AttributesEncodingOptions{Prefix: "user"},
		)
		if err != nil {
			return nil, d.createError(err)
		}

		if file.Mode.IsRegular() {
			objectType = common.Object_REGULAR
		} else if file.Mode.IsDir() {
			objectType = common.Object_DIRECTORY
		}

		head := &service.HeadObjectResponse{
			Metadata: &common.Metadata{
				Id:           file.BaseName,
				Key:          key,
				Size:         file.Size,
				ObjectType:   objectType,
				ModifiedTime: timestamppb.New(file.ModifiedTime),
				Attributes:   attributes,
			},
		}

		objects = append(objects, head)
	}

	return &service.ListObjectsResponse{Objects: objects}, nil
}

// Descriptor returns a structured Protobuf-defined descriptor of this backend.
func (d *Disk) Descriptor() *common.Node {
	return &common.Node{
		Name: "disk",
		Children: map[string]*common.Node_Value{
			"root": {Child: &common.Node_Value_Value{Value: d.cfg.Root}},
		},
	}
}

// Close is a noop.
func (d *Disk) Close() error {
	return nil
}

// String returns a human-consumable representation of this backend.
func (d *Disk) String() string {
	return schemas.MarshalDescriptor(d.Descriptor())
}

// createError creates a semantically accurate gRPC status for errored filesystem operations.
func (d *Disk) createError(err error) error {
	switch {
	case os.IsNotExist(err):
		return status.Error(codes.NotFound, err.Error())
	case os.IsPermission(err):
		return status.Error(codes.PermissionDenied, err.Error())
	case os.IsTimeout(err):
		return status.Error(codes.DeadlineExceeded, err.Error())
	case errors.Is(err, os.ErrInvalid):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, syscall.ENOTEMPTY):
		return status.Error(codes.FailedPrecondition, err.Error())
	default:
		return status.Error(codes.Unknown, err.Error())
	}
}

// fileTarget is an abstraction over a file target on disk, as constructed from a bucket name and
// file key.
type fileTarget struct {
	client *disk.Client
	bucket string
	key    string
}

// newFileTarget creates a new file target while validating the desired path.
func newFileTarget(client *disk.Client, bucket string, key string) (*fileTarget, error) {
	zap.L().Debug(
		"creating disk file target",
		zap.String("bucket", bucket),
		zap.String("key", key),
	)

	if strings.Contains(bucket, string(os.PathSeparator)) {
		return nil, fmt.Errorf("bucket name contains illegal characters: bucket=%s", bucket)
	}

	canonicalKey := filepath.Clean(key)
	resolved := filepath.Join(bucket, canonicalKey)

	// Require the resolved path to start with the bucket after a pass through filepath.Clean.
	// Since at this point it is known that the bucket does not contain nested directories, this
	// thwarts potential directory traversal attacks by enforcing that the requested key does
	// not resolve to be outside the allowed file tree.
	if canonicalKey != "." && !strings.HasPrefix(resolved, fmt.Sprintf("%s%c", bucket, os.PathSeparator)) {
		return nil, fmt.Errorf("resolved file path is invalid: key=%s", key)
	}

	return &fileTarget{
		client: client,
		bucket: bucket,
		key:    canonicalKey,
	}, nil
}

// resolve resolves the file path targeted by the supplied bucket and key.
func (f *fileTarget) resolve() string {
	return filepath.Join(f.bucket, f.key)
}

// parent returns the parent directory of the target. If the target is itself a directory, the
// return value is the same as that of resolve. If the target is a file or file prefix, the nearest
// path ancestor is returned.
func (f *fileTarget) parent() string {
	resolved := f.resolve()

	// The requested path may itself already be a directory, in which case the parsed dirname
	// should be discarded.
	p := filepath.Dir(resolved)
	stat, err := f.client.StatFile(resolved)
	if p == "." || (err == nil && stat.Mode.IsDir()) {
		p = resolved
	}

	return p
}

package backend

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"lib.kevinlin.info/aperture/lib"

	"unistore/internal/metrics"
	"unistore/internal/schemas"
	"unistore/schemas/common"
	"unistore/schemas/service"
)

// Metric names emitted by the instrumentation middleware.
const (
	metricInvoke             = "backend.invoke"
	metricDuration           = "backend.duration"
	metricRequestSize        = "backend.request.size"
	metricResponseSize       = "backend.response.size"
	metricStreamStart        = "backend.stream.start"
	metricStreamEnd          = "backend.stream.end"
	metricStreamMessageSend  = "backend.stream.message.send"
	metricStreamMessageRecv  = "backend.stream.message.recv"
	metricStreamMessageSize  = "backend.stream.message.size"
	metricStreamMessageCount = "backend.stream.message.count"
)

// instrumentation is a Backend that provides transparent instrumentation over an existing Backend.
type instrumentation struct {
	name string
	Backend
}

// newInstrumentation creates a new instrumentation backend with the specified name and wrapped
// Backend. The name is typically a label associated with the underlying backend.
func newInstrumentation(name string, backend Backend) Backend {
	return &instrumentation{
		name:    name,
		Backend: backend,
	}
}

// HeadBucket defers to the underlying backend and emits relevant metrics.
func (i *instrumentation) HeadBucket(ctx context.Context, request *service.HeadBucketRequest) (*service.HeadBucketResponse, error) {
	stopwatch := lib.NewStopwatch()
	base := lib.MapTag(map[string]interface{}{
		"method":  "HeadBucket",
		"name":    i.name,
		"backend": request.Resource.Backend,
		"bucket":  request.Resource.Bucket,
	})

	metrics.Client.Size(metricRequestSize, int64(proto.Size(request)), lib.CombineTags(base))

	response, err := i.Backend.HeadBucket(ctx, request)
	if err != nil {
		st := status.Convert(err)
		tags := lib.CombineTags(base, lib.Tag("code", st.Code()))

		metrics.Client.Incr(metricInvoke, tags)
		metrics.Client.Timing(metricDuration, stopwatch.Elapsed(), tags)

		return nil, err
	}

	tags := lib.CombineTags(base, lib.Tag("code", codes.OK))
	metrics.Client.Incr(metricInvoke, tags)
	metrics.Client.Timing(metricDuration, stopwatch.Elapsed(), tags)
	metrics.Client.Size(metricResponseSize, int64(proto.Size(response)), tags)

	return response, err
}

// HeadObject defers to the underlying backend and emits relevant metrics.
func (i *instrumentation) HeadObject(ctx context.Context, request *service.HeadObjectRequest) (*service.HeadObjectResponse, error) {
	stopwatch := lib.NewStopwatch()
	base := lib.MapTag(map[string]interface{}{
		"method":  "HeadObject",
		"name":    i.name,
		"backend": request.Resource.Backend,
		"bucket":  request.Resource.Bucket,
	})

	metrics.Client.Size(metricRequestSize, int64(proto.Size(request)), lib.CombineTags(base))

	response, err := i.Backend.HeadObject(ctx, request)
	if err != nil {
		st := status.Convert(err)
		tags := lib.CombineTags(base, lib.Tag("code", st.Code()))

		metrics.Client.Incr(metricInvoke, tags)
		metrics.Client.Timing(metricDuration, stopwatch.Elapsed(), tags)

		return nil, err
	}

	tags := lib.CombineTags(base, lib.Tag("code", codes.OK))
	metrics.Client.Incr(metricInvoke, tags)
	metrics.Client.Timing(metricDuration, stopwatch.Elapsed(), tags)
	metrics.Client.Size(metricResponseSize, int64(proto.Size(response)), tags)

	return response, err
}

// GetObject defers to the underlying backend and emits relevant metrics.
func (i *instrumentation) GetObject(ctx context.Context, request *service.GetObjectRequest) (*service.GetObjectResponse, error) {
	stopwatch := lib.NewStopwatch()
	base := lib.MapTag(map[string]interface{}{
		"method":  "GetObject",
		"name":    i.name,
		"backend": request.Resource.Backend,
		"bucket":  request.Resource.Bucket,
	})

	metrics.Client.Size(metricRequestSize, int64(proto.Size(request)), lib.CombineTags(base))

	response, err := i.Backend.GetObject(ctx, request)
	if err != nil {
		st := status.Convert(err)
		tags := lib.CombineTags(base, lib.Tag("code", st.Code()))

		metrics.Client.Incr(metricInvoke, tags)
		metrics.Client.Timing(metricDuration, stopwatch.Elapsed(), tags)

		return nil, err
	}

	tags := lib.CombineTags(base, lib.Tag("code", codes.OK))
	metrics.Client.Incr(metricInvoke, tags)
	metrics.Client.Timing(metricDuration, stopwatch.Elapsed(), tags)
	metrics.Client.Size(metricResponseSize, int64(proto.Size(response)), tags)

	return response, err
}

// StreamGetObject defers to the underlying backend and emits relevant metrics.
func (i *instrumentation) StreamGetObject(ctx context.Context, request *service.GetObjectStreamRequest) (<-chan *service.GetObjectStreamResponse, chan error) {
	responseProxy := make(chan *service.GetObjectStreamResponse)
	errsProxy := make(chan error, 1)

	messages := 0
	stopwatch := lib.NewStopwatch()
	base := map[string]interface{}{
		"method":  "StreamGetObject",
		"name":    i.name,
		"backend": request.Request.Resource.Backend,
		"bucket":  request.Request.Resource.Bucket,
	}

	metrics.Client.Incr(metricStreamStart, base)
	metrics.Client.Size(metricRequestSize, int64(proto.Size(request)), base)

	stream, errs := i.Backend.StreamGetObject(ctx, request)
	if stream == nil {
		err := <-errs
		st := status.Convert(err)
		tags := lib.CombineTags(lib.MapTag(base), lib.Tag("code", st.Code()))

		metrics.Client.Incr(metricInvoke, tags)
		metrics.Client.Timing(metricDuration, stopwatch.Elapsed(), tags)

		errsProxy <- err
		return nil, errsProxy
	}

	go func() {
		for {
			select {
			case err := <-errs:
				st := status.Convert(err)
				tags := lib.CombineTags(lib.MapTag(base), lib.Tag("code", st.Code()))

				metrics.Client.Incr(metricInvoke, tags)
				metrics.Client.Incr(metricStreamEnd, tags)
				metrics.Client.Timing(metricDuration, stopwatch.Elapsed(), tags)
				metrics.Client.Size(metricStreamMessageCount, int64(messages), tags)

				errsProxy <- err
				return
			case response, ok := <-stream:
				if !ok {
					tags := lib.CombineTags(lib.MapTag(base), lib.Tag("code", codes.OK))
					metrics.Client.Incr(metricInvoke, tags)
					metrics.Client.Incr(metricStreamEnd, tags)
					metrics.Client.Timing(metricDuration, stopwatch.Elapsed(), tags)
					metrics.Client.Size(metricStreamMessageCount, int64(messages), tags)

					close(responseProxy)
					return
				}

				metrics.Client.Incr(metricStreamMessageRecv, base)
				metrics.Client.Size(
					metricStreamMessageSize,
					int64(proto.Size(response)),
					base,
				)

				messages++
				responseProxy <- response
			}
		}
	}()

	return responseProxy, errsProxy
}

// PutObject defers to the underlying backend and emits relevant metrics.
func (i *instrumentation) PutObject(ctx context.Context, request *service.PutObjectRequest) (*service.PutObjectResponse, error) {
	stopwatch := lib.NewStopwatch()
	base := lib.MapTag(map[string]interface{}{
		"method":  "PutObject",
		"name":    i.name,
		"backend": request.Resource.Backend,
		"bucket":  request.Resource.Bucket,
	})

	metrics.Client.Size(metricRequestSize, int64(proto.Size(request)), lib.CombineTags(base))

	response, err := i.Backend.PutObject(ctx, request)
	if err != nil {
		st := status.Convert(err)
		tags := lib.CombineTags(base, lib.Tag("code", st.Code()))

		metrics.Client.Incr(metricInvoke, tags)
		metrics.Client.Timing(metricDuration, stopwatch.Elapsed(), tags)

		return nil, err
	}

	tags := lib.CombineTags(base, lib.Tag("code", codes.OK))
	metrics.Client.Incr(metricInvoke, tags)
	metrics.Client.Timing(metricDuration, stopwatch.Elapsed(), tags)
	metrics.Client.Size(metricResponseSize, int64(proto.Size(response)), tags)

	return response, err
}

// StreamPutObject defers to the underlying backend and emits relevant metrics.
func (i *instrumentation) StreamPutObject(ctx context.Context, stream chan *service.PutObjectStreamRequest) (*service.PutObjectStreamResponse, error) {
	var request *service.PutObjectStreamRequest

	proxy := make(chan *service.PutObjectStreamRequest)

	messages := 0
	stopwatch := lib.NewStopwatch()
	base := lib.MapTag(map[string]interface{}{
		"method": "StreamPutObject",
		"name":   i.name,
	})

	metrics.Client.Incr(metricStreamStart, lib.CombineTags(base))

	go func() {
		for request = range stream {
			tags := lib.CombineTags(
				base,
				lib.Tag("backend", request.Request.Resource.Backend),
				lib.Tag("bucket", request.Request.Resource.Bucket),
			)

			metrics.Client.Incr(metricStreamMessageSend, tags)
			metrics.Client.Size(metricStreamMessageSize, int64(proto.Size(request)), tags)

			messages++
			proxy <- request
		}

		close(proxy)
	}()

	response, err := i.Backend.StreamPutObject(ctx, proxy)
	if err != nil {
		st := status.Convert(err)

		backend := common.Backend_EMPTY
		bucket := "unknown"
		if request != nil {
			backend = request.Request.Resource.Backend
			bucket = request.Request.Resource.Bucket
		}

		tags := lib.CombineTags(
			base,
			lib.Tag("backend", backend),
			lib.Tag("bucket", bucket),
			lib.Tag("code", st.Code()),
		)

		metrics.Client.Incr(metricInvoke, tags)
		metrics.Client.Incr(metricStreamEnd, tags)
		metrics.Client.Timing(metricDuration, stopwatch.Elapsed(), tags)
		metrics.Client.Size(metricResponseSize, int64(proto.Size(response)), tags)
		metrics.Client.Size(metricStreamMessageCount, int64(messages), tags)

		return nil, err
	}

	tags := lib.CombineTags(
		base,
		lib.Tag("backend", request.Request.Resource.Backend),
		lib.Tag("bucket", request.Request.Resource.Bucket),
		lib.Tag("code", codes.OK),
	)
	metrics.Client.Incr(metricInvoke, tags)
	metrics.Client.Incr(metricStreamEnd, tags)
	metrics.Client.Timing(metricDuration, stopwatch.Elapsed(), tags)
	metrics.Client.Size(metricResponseSize, int64(proto.Size(response)), tags)
	metrics.Client.Size(metricStreamMessageCount, int64(messages), tags)

	return response, nil
}

// DeleteObject defers to the underlying backend and emits relevant metrics.
func (i *instrumentation) DeleteObject(ctx context.Context, request *service.DeleteObjectRequest) (*service.DeleteObjectResponse, error) {
	stopwatch := lib.NewStopwatch()
	base := lib.MapTag(map[string]interface{}{
		"method":  "DeleteObject",
		"name":    i.name,
		"backend": request.Resource.Backend,
		"bucket":  request.Resource.Bucket,
	})

	metrics.Client.Size(metricRequestSize, int64(proto.Size(request)), lib.CombineTags(base))

	response, err := i.Backend.DeleteObject(ctx, request)
	if err != nil {
		st := status.Convert(err)
		tags := lib.CombineTags(base, lib.Tag("code", st.Code()))

		metrics.Client.Incr(metricInvoke, tags)
		metrics.Client.Timing(metricDuration, stopwatch.Elapsed(), tags)

		return nil, err
	}

	tags := lib.CombineTags(base, lib.Tag("code", codes.OK))
	metrics.Client.Incr(metricInvoke, tags)
	metrics.Client.Timing(metricDuration, stopwatch.Elapsed(), tags)
	metrics.Client.Size(metricResponseSize, int64(proto.Size(response)), tags)

	return response, err
}

// ListBuckets defers to the underlying backend and emits relevant metrics.
func (i *instrumentation) ListBuckets(ctx context.Context, request *service.ListBucketsRequest) (*service.ListBucketsResponse, error) {
	stopwatch := lib.NewStopwatch()
	base := lib.MapTag(map[string]interface{}{
		"method":  "ListBuckets",
		"name":    i.name,
		"backend": request.Resource.Backend,
	})

	metrics.Client.Size(metricRequestSize, int64(proto.Size(request)), lib.CombineTags(base))

	response, err := i.Backend.ListBuckets(ctx, request)
	if err != nil {
		st := status.Convert(err)
		tags := lib.CombineTags(base, lib.Tag("code", st.Code()))

		metrics.Client.Incr(metricInvoke, tags)
		metrics.Client.Timing(metricDuration, stopwatch.Elapsed(), tags)

		return nil, err
	}

	tags := lib.CombineTags(base, lib.Tag("code", codes.OK))
	metrics.Client.Incr(metricInvoke, tags)
	metrics.Client.Timing(metricDuration, stopwatch.Elapsed(), tags)
	metrics.Client.Size(metricResponseSize, int64(proto.Size(response)), tags)

	return response, err
}

// ListObjects defers to the underlying backend and emits relevant metrics.
func (i *instrumentation) ListObjects(ctx context.Context, request *service.ListObjectsRequest) (*service.ListObjectsResponse, error) {
	stopwatch := lib.NewStopwatch()
	base := lib.MapTag(map[string]interface{}{
		"method":  "ListObjects",
		"name":    i.name,
		"backend": request.Resource.Backend,
		"bucket":  request.Resource.Bucket,
	})

	metrics.Client.Size(metricRequestSize, int64(proto.Size(request)), lib.CombineTags(base))

	response, err := i.Backend.ListObjects(ctx, request)
	if err != nil {
		st := status.Convert(err)
		tags := lib.CombineTags(base, lib.Tag("code", st.Code()))

		metrics.Client.Incr(metricInvoke, tags)
		metrics.Client.Timing(metricDuration, stopwatch.Elapsed(), tags)

		return nil, err
	}

	tags := lib.CombineTags(base, lib.Tag("code", codes.OK))
	metrics.Client.Incr(metricInvoke, tags)
	metrics.Client.Timing(metricDuration, stopwatch.Elapsed(), tags)
	metrics.Client.Size(metricResponseSize, int64(proto.Size(response)), tags)

	return response, err
}

// Descriptor returns a structured Protobuf-defined descriptor of this backend.
func (i *instrumentation) Descriptor() *common.Node {
	return &common.Node{
		Name: "instrumentation",
		Children: map[string]*common.Node_Value{
			"name":    {Child: &common.Node_Value_Value{Value: i.name}},
			"backend": {Child: &common.Node_Value_Node{Node: i.Backend.Descriptor()}},
		},
	}
}

// String returns a human-consumable representation of this backend.
func (i *instrumentation) String() string {
	return schemas.MarshalDescriptor(i.Descriptor())
}

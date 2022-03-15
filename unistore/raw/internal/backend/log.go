package backend

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"syscall"
	"time"

	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
	"lib.kevinlin.info/aperture/lib"

	"unistore/internal/schemas"
	"unistore/internal/util"
	"unistore/schemas/common"
	"unistore/schemas/service"
)

// Log is a Backend that logs all requests with request and response parameters.
type Log struct {
	level  zapcore.Level
	files  []string
	output io.WriteCloser
	tags   map[string]string
	mutex  sync.Mutex
	closed bool
	Backend
}

// NewLog creates a new Log with the specified log level, outputs, and extra static tags.
func NewLog(level string, outputs []string, tags map[string]string, backend Backend) (Backend, error) {
	wc := make([]io.WriteCloser, len(outputs))

	// Use the convenience methods from Zap for log eligibility based on level.
	zapLevel, err := zapcore.ParseLevel(level)
	if err != nil {
		return nil, err
	}

	for i, output := range outputs {
		switch output {
		case "stdout":
			fd, err := syscall.Dup(int(os.Stdout.Fd()))
			if err != nil {
				return nil, err
			}

			wc[i] = os.NewFile(uintptr(fd), "stdout")
		case "stderr":
			fd, err := syscall.Dup(int(os.Stderr.Fd()))
			if err != nil {
				return nil, err
			}

			wc[i] = os.NewFile(uintptr(fd), "stderr")
		default:
			wc[i], err = os.OpenFile(
				output,
				syscall.O_WRONLY|syscall.O_APPEND|syscall.O_CREAT,
				0644,
			)
			if err != nil {
				return nil, err
			}
		}
	}

	l := &Log{
		level:   zapLevel,
		files:   outputs,
		output:  util.NewMultiWriteCloser(wc...),
		tags:    tags,
		Backend: backend,
	}

	return newInstrumentation("log", l), nil
}

// HeadBucket invokes the underlying backend and logs the request on completion.
func (l *Log) HeadBucket(ctx context.Context, request *service.HeadBucketRequest) (*service.HeadBucketResponse, error) {
	stopwatch := lib.NewStopwatch()
	response, err := l.Backend.HeadBucket(ctx, request)

	e := entry{
		Code:     status.Convert(err).Code().String(),
		Error:    status.Convert(err).Message(),
		Duration: float64(stopwatch.Elapsed().Microseconds()) / 1000.0,
		Method:   "HeadBucket",
		Backend:  request.Resource.Backend.String(),
		Bucket:   request.Resource.Bucket,
	}

	go l.log(ctx, e)

	return response, err
}

// HeadObject invokes the underlying backend and logs the request on completion.
func (l *Log) HeadObject(ctx context.Context, request *service.HeadObjectRequest) (*service.HeadObjectResponse, error) {
	stopwatch := lib.NewStopwatch()
	response, err := l.Backend.HeadObject(ctx, request)

	e := entry{
		Code:     status.Convert(err).Code().String(),
		Error:    status.Convert(err).Message(),
		Duration: float64(stopwatch.Elapsed().Microseconds()) / 1000.0,
		Method:   "HeadObject",
		Backend:  request.Resource.Backend.String(),
		Bucket:   request.Resource.Bucket,
		Key:      request.Key,
	}
	if response != nil && response.Metadata != nil {
		e.Size = response.Metadata.Size
	}

	go l.log(ctx, e)

	return response, err
}

// GetObject invokes the underlying backend and logs the request on completion.
func (l *Log) GetObject(ctx context.Context, request *service.GetObjectRequest) (*service.GetObjectResponse, error) {
	stopwatch := lib.NewStopwatch()
	response, err := l.Backend.GetObject(ctx, request)

	e := entry{
		Code:     status.Convert(err).Code().String(),
		Error:    status.Convert(err).Message(),
		Duration: float64(stopwatch.Elapsed().Microseconds()) / 1000.0,
		Method:   "GetObject",
		Backend:  request.Resource.Backend.String(),
		Bucket:   request.Resource.Bucket,
		Key:      request.Key,
	}
	if response != nil && response.Metadata != nil {
		e.Size = response.Metadata.Size
	}

	go l.log(ctx, e)

	return response, err
}

// StreamGetObject invokes the underlying backend and logs the request on completion.
func (l *Log) StreamGetObject(ctx context.Context, request *service.GetObjectStreamRequest) (<-chan *service.GetObjectStreamResponse, chan error) {
	responseProxy := make(chan *service.GetObjectStreamResponse)
	errsProxy := make(chan error, 1)
	stopwatch := lib.NewStopwatch()

	e := entry{
		Method:  "StreamGetObject",
		Backend: request.Request.Resource.Backend.String(),
		Bucket:  request.Request.Resource.Bucket,
		Key:     request.Request.Key,
	}

	stream, errs := l.Backend.StreamGetObject(ctx, request)
	if stream == nil {
		err := <-errs
		st := status.Convert(err)

		e.Code = st.Code().String()
		e.Error = st.Message()
		e.Duration = float64(stopwatch.Elapsed().Microseconds()) / 1000.0
		go l.log(ctx, e)

		errsProxy <- err
		return nil, errsProxy
	}

	go func() {
		var size uint64

		for {
			select {
			case err := <-errs:
				st := status.Convert(err)

				e.Code = st.Code().String()
				e.Error = st.Message()
				e.Duration = float64(stopwatch.Elapsed().Microseconds()) / 1000.0
				e.Size = size
				go l.log(ctx, e)

				errsProxy <- err
				return
			case response, ok := <-stream:
				if !ok {
					e.Code = codes.OK.String()
					e.Duration = float64(stopwatch.Elapsed().Microseconds()) / 1000.0
					e.Size = size
					go l.log(ctx, e)

					close(responseProxy)
					return
				}

				size = response.Response.Metadata.Size

				responseProxy <- response
			}
		}
	}()

	return responseProxy, errsProxy
}

// PutObject invokes the underlying backend and logs the request on completion.
func (l *Log) PutObject(ctx context.Context, request *service.PutObjectRequest) (*service.PutObjectResponse, error) {
	stopwatch := lib.NewStopwatch()
	response, err := l.Backend.PutObject(ctx, request)

	e := entry{
		Code:     status.Convert(err).Code().String(),
		Error:    status.Convert(err).Message(),
		Duration: float64(stopwatch.Elapsed().Microseconds()) / 1000.0,
		Method:   "PutObject",
		Backend:  request.Resource.Backend.String(),
		Bucket:   request.Resource.Bucket,
		Key:      request.Key,
		Size:     uint64(len(request.Data)),
	}

	go l.log(ctx, e)

	return response, err
}

// StreamPutObject invokes the underlying backend and logs the request on completion.
func (l *Log) StreamPutObject(ctx context.Context, stream chan *service.PutObjectStreamRequest) (*service.PutObjectStreamResponse, error) {
	var request *service.PutObjectStreamRequest
	var size uint64

	proxy := make(chan *service.PutObjectStreamRequest)
	stopwatch := lib.NewStopwatch()

	e := entry{
		Method:  "StreamPutObject",
		Backend: common.Backend_EMPTY.String(),
		Bucket:  "unknown",
	}

	go func() {
		for request = range stream {
			size += uint64(len(request.Request.Data))

			proxy <- request
		}

		close(proxy)
	}()

	response, err := l.Backend.StreamPutObject(ctx, proxy)
	if err != nil {
		st := status.Convert(err)

		e.Code = st.Code().String()
		e.Error = st.Message()
		e.Duration = float64(stopwatch.Elapsed().Microseconds()) / 1000.0
		e.Size = size
		if request != nil {
			e.Backend = request.Request.Resource.Backend.String()
			e.Bucket = request.Request.Resource.Bucket
			e.Key = request.Request.Key
		}
		go l.log(ctx, e)

		return nil, err
	}

	e.Code = codes.OK.String()
	e.Duration = float64(stopwatch.Elapsed().Microseconds()) / 1000.0
	e.Backend = request.Request.Resource.Backend.String()
	e.Bucket = request.Request.Resource.Bucket
	e.Key = request.Request.Key
	e.Size = size
	go l.log(ctx, e)

	return response, err
}

// DeleteObject invokes the underlying backend and logs the request on completion.
func (l *Log) DeleteObject(ctx context.Context, request *service.DeleteObjectRequest) (*service.DeleteObjectResponse, error) {
	stopwatch := lib.NewStopwatch()
	response, err := l.Backend.DeleteObject(ctx, request)

	e := entry{
		Code:     status.Convert(err).Code().String(),
		Error:    status.Convert(err).Message(),
		Duration: float64(stopwatch.Elapsed().Microseconds()) / 1000.0,
		Method:   "DeleteObject",
		Backend:  request.Resource.Backend.String(),
		Bucket:   request.Resource.Bucket,
		Key:      request.Key,
	}

	go l.log(ctx, e)

	return response, err
}

// ListBuckets invokes the underlying backend and logs the request on completion.
func (l *Log) ListBuckets(ctx context.Context, request *service.ListBucketsRequest) (*service.ListBucketsResponse, error) {
	stopwatch := lib.NewStopwatch()
	response, err := l.Backend.ListBuckets(ctx, request)

	e := entry{
		Code:     status.Convert(err).Code().String(),
		Error:    status.Convert(err).Message(),
		Duration: float64(stopwatch.Elapsed().Microseconds()) / 1000.0,
		Method:   "ListBuckets",
		Backend:  request.Resource.Backend.String(),
		Bucket:   request.Resource.Bucket,
	}
	if response != nil {
		e.Size = uint64(len(response.Buckets))
	}

	go l.log(ctx, e)

	return response, err
}

// ListObjects invokes the underlying backend and logs the request on completion.
func (l *Log) ListObjects(ctx context.Context, request *service.ListObjectsRequest) (*service.ListObjectsResponse, error) {
	stopwatch := lib.NewStopwatch()
	response, err := l.Backend.ListObjects(ctx, request)

	e := entry{
		Code:     status.Convert(err).Code().String(),
		Error:    status.Convert(err).Message(),
		Duration: float64(stopwatch.Elapsed().Microseconds()) / 1000.0,
		Method:   "ListObjects",
		Backend:  request.Resource.Backend.String(),
		Bucket:   request.Resource.Bucket,
		Key:      request.Prefix,
	}
	if response != nil {
		e.Size = uint64(len(response.Objects))
	}

	go l.log(ctx, e)

	return response, err
}

// Descriptor returns a structured Protobuf-defined descriptor of this backend.
func (l *Log) Descriptor() *common.Node {
	return &common.Node{
		Name: "log",
		Children: map[string]*common.Node_Value{
			"level": {
				Child: &common.Node_Value_Value{
					Value: l.level.String(),
				},
			},
			"outputs": {
				Child: &common.Node_Value_Value{
					Value: strings.Join(l.files, ", "),
				},
			},
			"backend": {
				Child: &common.Node_Value_Node{
					Node: l.Backend.Descriptor(),
				},
			},
		},
	}
}

// Close idempotently closes the io.WriteClosers used for logging outputs.
func (l *Log) Close() error {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	if l.closed {
		return nil
	}

	l.closed = true
	return l.output.Close()
}

// String returns a human-consumable representation of this backend.
func (l *Log) String() string {
	return schemas.MarshalDescriptor(l.Descriptor())
}

// log populates derived fields in the log entry, serializes the entry, and writes it to the
// the declared output targets when eligible per the configured log level.
func (l *Log) log(ctx context.Context, e entry) {
	if e.Error != "" && !l.level.Enabled(zapcore.ErrorLevel) {
		return
	} else if e.Error == "" && !l.level.Enabled(zapcore.InfoLevel) {
		return
	}

	// Populate some fields automatically

	e.Timestamp = time.Now()
	e.Tags = l.tags
	e.ClientID = "unknown"
	if p, ok := peer.FromContext(ctx); ok {
		e.Peer = fmt.Sprintf("%s:%s", p.Addr.Network(), p.Addr.String())
	}
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if ua := md.Get("User-Agent"); len(ua) > 0 {
			e.ClientID = strings.Join(ua, ", ")
		}
	}

	serialized, err := e.MarshalBinary()
	if err != nil {
		return
	}

	l.output.Write(append(serialized, []byte("\n")...))
}

// entry represents a single structured log entry for a request.
type entry struct {
	Timestamp time.Time         `json:"timestamp"`
	Code      string            `json:"code"`
	Error     string            `json:"error"`
	Duration  float64           `json:"duration"`
	Method    string            `json:"method"`
	Backend   string            `json:"backend"`
	Bucket    string            `json:"bucket"`
	Key       string            `json:"key"`
	Size      uint64            `json:"size"`
	Peer      string            `json:"peer"`
	ClientID  string            `json:"client_id"`
	Tags      map[string]string `json:"tags"`
}

// MarshalBinary implements the encoding.BinaryMarshaler interface. It serializes the log entry as
// a JSON string using json.Marshal.
func (e *entry) MarshalBinary() ([]byte, error) {
	return json.Marshal(e)
}

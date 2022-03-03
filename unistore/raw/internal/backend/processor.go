package backend

import (
	"bufio"
	"bytes"
	"context"
	"io"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	"unistore/internal/util"
	"unistore/schemas/common"
	"unistore/schemas/service"
)

var (
	// defaultStreamingWriteChunkSize is the default chunk size used for streaming puts to the
	// backend. An explicit chunk size is required since the volume of put chunks to the
	// backend may be mismatched relative to the input if there is a non-unity size
	// amplification factor.
	defaultStreamingWriteChunkSize = uint64(5 * 1024 * 1024) // 5 MB
)

// readerFactory is a type alias for a function that creates an io.ReadCloser from an io.Reader.
type readerFactory func(io.Reader) (io.ReadCloser, error)

// writerFactory is a type alias for a function that creates an io.WriteCloser from an io.Writer.
type writerFactory func(io.Writer) (io.WriteCloser, error)

// ioProcessor is a wrapper over a Backend that performs in-flight mutation of data streams using
// the provided io.Reader and io.Writer factories. In particular, reads (gets and streaming gets)
// proxy their data through the provided io.Reader factory; writes (puts and streaming puts) proxy
// their data through the provided io.Writer factory.
//
// The io.Reader factory accepts an io.Reader for original input, and provides an io.Reader from
// which the transformed data stream can be read.
// The io.Writer factory accepts an io.Writer into which the transformed output is written, and
// provides an io.Writer into which the original data should be written.
//
// The streaming read and write chunk sizes are used to specify a desired chunk size for reads and
// writes from and to the proxied io.Reader and io.Writer, respectively. When 0, the streaming read
// routine retains the original chunk size specified in the request, and the streaming write routine
// uses a default chunk size.
type ioProcessor struct {
	reader                  func(io.Reader) (io.ReadCloser, error)
	writer                  func(io.Writer) (io.WriteCloser, error)
	streamingReadChunkSize  uint64
	streamingWriteChunkSize uint64
	Backend
}

// newIOProcessor creates a new ioProcessor.
func newIOProcessor(
	reader readerFactory,
	writer writerFactory,
	streamingReadChunkSize uint64,
	streamingWriteChunkSize uint64,
	backend Backend,
) Backend {
	if reader == nil {
		reader = func(raw io.Reader) (io.ReadCloser, error) {
			return io.NopCloser(raw), nil
		}
	}

	if writer == nil {
		writer = func(raw io.Writer) (io.WriteCloser, error) {
			return util.NopWriteCloser(raw), nil
		}
	}

	if streamingWriteChunkSize == 0 {
		streamingWriteChunkSize = defaultStreamingWriteChunkSize
	}

	return &ioProcessor{
		reader:                  reader,
		writer:                  writer,
		streamingReadChunkSize:  streamingReadChunkSize,
		streamingWriteChunkSize: streamingWriteChunkSize,
		Backend:                 backend,
	}
}

// GetObject defers to the underlying backend followed by performing transparent transformation of
// the full payload.
func (p *ioProcessor) GetObject(ctx context.Context, request *service.GetObjectRequest) (*service.GetObjectResponse, error) {
	resp, err := p.Backend.GetObject(ctx, request)
	if err != nil {
		return nil, err
	}

	proxyReader, err := p.reader(bytes.NewReader(resp.Data))
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	defer proxyReader.Close()

	resp.Data, err = io.ReadAll(proxyReader)
	if err != nil {
		return nil, status.Error(codes.Unknown, err.Error())
	}

	return resp, nil
}

// StreamGetObject is a streaming implementation of GetObject that transforms data chunk-by-chunk.
func (p *ioProcessor) StreamGetObject(ctx context.Context, request *service.GetObjectStreamRequest) (<-chan *service.GetObjectStreamResponse, chan error) {
	var metadata *common.Metadata

	proxyRequest := request
	if p.streamingReadChunkSize > 0 {
		// Clone the inbound request in order to overwrite the chunk size to the desired
		// streaming read chunk size. The request-level chunk reshaping applies only to the
		// chunk reads done internally by the backend; the original chunk size requested by
		// the client is still respected for outbound stream responses.
		proxyRequest = proto.Clone(request).(*service.GetObjectStreamRequest)
		proxyRequest.ChunkSize = p.streamingReadChunkSize
	}

	stream, errs := p.Backend.StreamGetObject(ctx, proxyRequest)
	if stream == nil {
		return nil, errs
	}

	responseProxy := make(chan *service.GetObjectStreamResponse)
	errsProxy := make(chan error, 1)
	r, w := io.Pipe()

	if request.ChunkSize == 0 {
		errsProxy <- status.Error(
			codes.InvalidArgument,
			"backend: I/O processor streaming get chunk size must be greater than zero",
		)
		return nil, errsProxy
	}

	go func() {
		// The stream has concluded gracefully; close the writer to signal to the client
		// pipe that control should be transferred to the reader.
		defer w.Close()

		for {
			select {
			case err := <-errs:
				errsProxy <- err
				return
			case response, ok := <-stream:
				if !ok {
					return
				}

				// Capture the first available metadata payload to mirror in the
				// underlying response.
				if metadata == nil {
					metadata = response.Response.Metadata
				}

				// Simply relay the response data to the pipe for asynchronous
				// consumption by the reader routine.
				if _, err := w.Write(response.Response.Data); err != nil {
					errsProxy <- err
					return
				}
			}
		}
	}()

	go func() {
		var offset int

		proxyReader, err := p.reader(r)
		if err != nil {
			errsProxy <- err
			close(responseProxy)
			return
		}

		defer proxyReader.Close()

		for {
			chunk := make([]byte, request.ChunkSize)
			n, err := io.ReadFull(proxyReader, chunk)

			if n > 0 {
				responseProxy <- &service.GetObjectStreamResponse{
					Response: &service.GetObjectResponse{
						Metadata: metadata,
						Data:     chunk[:n],
					},
					Range: &common.Range{
						Unit:  "bytes",
						Start: uint64(offset),
						End:   uint64(offset + n),
						// Due to the streaming nature of the operation, the
						// total size is not known in advance.
					},
				}

				offset += n
			}

			if err != nil {
				if err != io.EOF && err != io.ErrUnexpectedEOF {
					errsProxy <- err
				}
				break
			}
		}

		close(responseProxy)
	}()

	return responseProxy, errsProxy
}

// PutObject transparently transforms the full data payload followed by replacing the data in the
// outbound request to the underlying backend.
func (p *ioProcessor) PutObject(ctx context.Context, request *service.PutObjectRequest) (*service.PutObjectResponse, error) {
	errs := make(chan error, 2)
	r, w := io.Pipe()

	go func() {
		defer w.Close()

		proxyWriter, err := p.writer(w)
		if err != nil {
			errs <- err
			return
		}

		defer proxyWriter.Close()

		if _, err = proxyWriter.Write(request.Data); err != nil {
			errs <- err
			return
		}

		errs <- nil
	}()

	go func() {
		transformed, err := io.ReadAll(r)
		if err != nil {
			errs <- err
			return
		}

		request.Data = transformed
		errs <- nil
	}()

	for i := 0; i < 2; i++ {
		if err := <-errs; err != nil {
			return nil, status.Error(codes.Unknown, err.Error())
		}
	}

	return p.Backend.PutObject(ctx, request)
}

// StreamPutObject is a streaming implementation of PutObject that transforms data chunk-by-chunk.
func (p *ioProcessor) StreamPutObject(ctx context.Context, stream chan *service.PutObjectStreamRequest) (*service.PutObjectStreamResponse, error) {
	var request *service.PutObjectStreamRequest

	proxyRequest := make(chan *service.PutObjectStreamRequest)
	responses := make(chan *service.PutObjectStreamResponse)
	errs := make(chan error)
	peek := <-stream

	if peek == nil {
		return nil, status.Error(
			codes.InvalidArgument,
			"backend: error reading first chunk from stream",
		)
	}

	if peek.Range != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			"backend: range specifications are not supported for streaming puts: range=%v",
			peek.Range,
		)
	}

	r, w := io.Pipe()

	// Wrap the pipe writer in a buffered writer in order to reduce the number of individual
	// reads (and thus backend put requests) that need to be invoked on the opposite end of the
	// pipe. Note that this requires explicitly flushing the buffer after no additional writes
	// are expected.
	wb := bufio.NewWriter(w)

	go func() {
		for {
			chunk := make([]byte, p.streamingWriteChunkSize)
			n, err := io.ReadFull(r, chunk)

			if n > 0 {
				// Pass along a duplicated instance of the request from the original
				// stream, but with the data overridden to the output from the I/O
				// processor. Note that there is almost necessarily a mismatch
				// between the number of requests in the original stream and the
				// number of outbound requests after applying mid-stream I/O; it is
				// impossible to achieve a 1:1 mapping of metadata between the
				// original stream and the transformed stream. Instead, this routine
				// makes a best-effort attempt to fetch the latest processed inbound
				// request as the base onto which the transformed data chunk is
				// applied.
				current := request
				if current == nil {
					current = peek
				}

				message := proto.Clone(current.Request).(*service.PutObjectRequest)
				message.Data = chunk[:n]
				proxyRequest <- &service.PutObjectStreamRequest{Request: message}
			}

			if err != nil {
				if err != io.EOF && err != io.ErrUnexpectedEOF {
					errs <- err
				}
				break
			}
		}

		close(proxyRequest)
	}()

	go func() {
		defer func() {
			// At this point, the streaming put has completed and there are no remaining
			// bytes to write to the output writer. Flush any remaining bytes from the
			// buffered writer and signal to the pipe reader via io.EOF that no
			// additional writes are to be expected.
			wb.Flush()
			w.Close()
		}()

		proxyWriter, err := p.writer(wb)
		if err != nil {
			errs <- err
			return
		}

		defer proxyWriter.Close()

		if _, err := proxyWriter.Write(peek.Request.Data); err != nil {
			errs <- err
			return
		}

		for request = range stream {
			if _, err := proxyWriter.Write(request.Request.Data); err != nil {
				errs <- err
				return
			}
		}
	}()

	go func() {
		response, err := p.Backend.StreamPutObject(ctx, proxyRequest)
		if err != nil {
			errs <- err
		} else {
			responses <- response
		}
	}()

	select {
	case err := <-errs:
		return nil, err
	case response := <-responses:
		return response, nil
	}
}

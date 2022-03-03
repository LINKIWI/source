package middleware

import (
	"context"
	"fmt"
	"io"
	"net"
	"strings"
	"sync/atomic"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"lib.kevinlin.info/aperture/lib"
)

// NewLogUnaryInterceptor creates a new grpc.UnaryServerInterceptor for logging metadata about unary
// requests.
func NewLogUnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		request interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		clientAddr := net.Addr(&net.TCPAddr{})
		clientID := "unknown"

		if p, ok := peer.FromContext(ctx); ok {
			clientAddr = p.Addr
		}

		if md, ok := metadata.FromIncomingContext(ctx); ok {
			if ua := md.Get("User-Agent"); len(ua) > 0 {
				clientID = strings.Join(ua, ", ")
			}
		}

		stopwatch := lib.NewStopwatch()
		response, err := handler(ctx, request)

		if err != nil {
			st := status.Convert(err)
			zap.L().Warn(
				"error serving gRPC unary request",
				zap.String("peer", fmt.Sprintf("%s:%s", clientAddr.Network(), clientAddr.String())),
				zap.String("client_id", clientID),
				zap.String("method", info.FullMethod),
				zap.Duration("duration", stopwatch.Elapsed()),
				zap.String("request", fmt.Sprintf("%T", request)),
				zap.Int("request_size", proto.Size(request.(proto.Message))),
				zap.Stringer("code", st.Code()),
				zap.String("error", st.Message()),
			)
		} else {
			zap.L().Debug(
				"served gRPC unary request",
				zap.String("peer", fmt.Sprintf("%s:%s", clientAddr.Network(), clientAddr.String())),
				zap.String("client_id", clientID),
				zap.String("method", info.FullMethod),
				zap.Duration("duration", stopwatch.Elapsed()),
				zap.String("request", fmt.Sprintf("%T", request)),
				zap.Int("request_size", proto.Size(request.(proto.Message))),
				zap.String("response", fmt.Sprintf("%T", response)),
				zap.Int("response_size", proto.Size(response.(proto.Message))),
			)
		}

		return response, err
	}
}

// NewLogStreamInterceptor creates a new grpc.StreamServerInterceptor for logging metadata about
// stream requests, including individual messages exchanged in the stream.
func NewLogStreamInterceptor() grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		clientAddr := net.Addr(&net.TCPAddr{})
		clientID := "unknown"

		if p, ok := peer.FromContext(ss.Context()); ok {
			clientAddr = p.Addr
		}

		if md, ok := metadata.FromIncomingContext(ss.Context()); ok {
			if ua := md.Get("User-Agent"); len(ua) > 0 {
				clientID = strings.Join(ua, ", ")
			}
		}

		stopwatch := lib.NewStopwatch()
		err := handler(srv, newLogServerStream(info, ss))

		if err != nil {
			st := status.Convert(err)
			zap.L().Warn(
				"error serving gRPC stream request",
				zap.String("peer", fmt.Sprintf("%s:%s", clientAddr.Network(), clientAddr.String())),
				zap.String("client_id", clientID),
				zap.String("method", info.FullMethod),
				zap.Duration("duration", stopwatch.Elapsed()),
				zap.Bool("client_stream", info.IsClientStream),
				zap.Bool("server_stream", info.IsServerStream),
				zap.Stringer("code", st.Code()),
				zap.String("error", st.Message()),
			)
		} else {
			zap.L().Debug(
				"served gRPC stream request",
				zap.String("peer", fmt.Sprintf("%s:%s", clientAddr.Network(), clientAddr.String())),
				zap.String("client_id", clientID),
				zap.String("method", info.FullMethod),
				zap.Duration("duration", stopwatch.Elapsed()),
				zap.Bool("client_stream", info.IsClientStream),
				zap.Bool("server_stream", info.IsServerStream),
			)
		}

		return err
	}
}

// logServerStream is a grpc.ServerStream wrapper that logs metadata about messages exchanged
// through the stream.
type logServerStream struct {
	info         *grpc.StreamServerInfo
	stopwatch    *lib.Stopwatch
	recvMessages *uint64
	sentMessages *uint64
	grpc.ServerStream
}

// newLogServerStream creates a new logServerStream while initializing internal state.
func newLogServerStream(info *grpc.StreamServerInfo, stream grpc.ServerStream) grpc.ServerStream {
	return &logServerStream{
		info:         info,
		stopwatch:    lib.NewStopwatch(),
		recvMessages: new(uint64),
		sentMessages: new(uint64),
		ServerStream: stream,
	}
}

// SendMsg wraps the underlying stream's SendMsg while providing logging.
func (l *logServerStream) SendMsg(m interface{}) error {
	err := l.ServerStream.SendMsg(m)

	if err != nil {
		st := status.Convert(err)
		zap.L().Warn(
			"error sending message over gRPC stream",
			zap.String("method", l.info.FullMethod),
			zap.Uint64("sequence", atomic.LoadUint64(l.sentMessages)),
			zap.Duration("elapsed", l.stopwatch.Elapsed()),
			zap.String("message", fmt.Sprintf("%T", m)),
			zap.Int("message_size", proto.Size(m.(proto.Message))),
			zap.Stringer("code", st.Code()),
			zap.String("error", st.Message()),
		)
	} else {
		atomic.AddUint64(l.sentMessages, 1)
		zap.L().Debug(
			"sent message over gRPC stream",
			zap.String("method", l.info.FullMethod),
			zap.Uint64("sequence", atomic.LoadUint64(l.sentMessages)),
			zap.Duration("elapsed", l.stopwatch.Elapsed()),
			zap.String("message", fmt.Sprintf("%T", m)),
			zap.Int("message_size", proto.Size(m.(proto.Message))),
		)
	}

	return err
}

// RecvMsg wraps the underlying stream's RecvMsg while providing logging.
func (l *logServerStream) RecvMsg(m interface{}) error {
	err := l.ServerStream.RecvMsg(m)

	if err != nil && err != io.EOF {
		st := status.Convert(err)
		zap.L().Warn(
			"error receiving message over gRPC stream",
			zap.String("method", l.info.FullMethod),
			zap.Uint64("sequence", atomic.LoadUint64(l.recvMessages)),
			zap.Duration("elapsed", l.stopwatch.Elapsed()),
			zap.String("message", fmt.Sprintf("%T", m)),
			zap.Int("message_size", proto.Size(m.(proto.Message))),
			zap.Stringer("code", st.Code()),
			zap.String("error", st.Message()),
		)
	} else {
		atomic.AddUint64(l.recvMessages, 1)
		zap.L().Debug(
			"received message over gRPC stream",
			zap.String("method", l.info.FullMethod),
			zap.Uint64("sequence", atomic.LoadUint64(l.recvMessages)),
			zap.Duration("elapsed", l.stopwatch.Elapsed()),
			zap.String("message", fmt.Sprintf("%T", m)),
			zap.Int("message_size", proto.Size(m.(proto.Message))),
		)
	}

	return err
}

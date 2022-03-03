package middleware

import (
	"context"
	"fmt"
	"strings"

	"github.com/getsentry/sentry-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

// NewErrorCaptureUnaryInterceptor creates a grpc.UnaryServerInterceptor that transparently reports
// gRPC service errors to Sentry.
func NewErrorCaptureUnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		request interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		response, err := handler(ctx, request)
		if err != nil {
			st := status.Convert(err)
			hub := sentry.CurrentHub().Clone()
			hub.ConfigureScope(func(scope *sentry.Scope) {
				scope.SetTag("code", st.Code().String())

				service, method, ok := parseServiceMethod(info.FullMethod)
				if ok {
					scope.SetTag("service", service)
					scope.SetTag("method", method)
				}

				p, ok := peer.FromContext(ctx)
				if ok {
					scope.SetExtra(
						"peer",
						fmt.Sprintf("%s:%s", p.Addr.Network(), p.Addr),
					)
				}
			})
			hub.CaptureException(err)
		}

		return response, err
	}
}

// NewErrorCaptureStreamInterceptor creates a grpc.StreamServerInterceptor that transparently
// reports gRPC service errors to Sentry when the stream concludes with an error.
func NewErrorCaptureStreamInterceptor() grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		ss grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		err := handler(srv, ss)
		if err != nil {
			st := status.Convert(err)
			hub := sentry.CurrentHub().Clone()
			hub.ConfigureScope(func(scope *sentry.Scope) {
				scope.SetTag("code", st.Code().String())

				service, method, ok := parseServiceMethod(info.FullMethod)
				if ok {
					scope.SetTag("service", service)
					scope.SetTag("method", method)
				}

				p, ok := peer.FromContext(ss.Context())
				if ok {
					scope.SetExtra(
						"peer",
						fmt.Sprintf("%s:%s", p.Addr.Network(), p.Addr),
					)
				}
			})
			hub.CaptureException(err)
		}

		return err
	}
}

func parseServiceMethod(fullMethod string) (string, string, bool) {
	components := strings.Split(fullMethod, "/")
	if len(components) != 3 {
		return "", "", false
	}

	return components[1], components[2], true
}

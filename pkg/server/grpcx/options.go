// Package grpcx provides ...
package grpcx

import (
	"context"

	"github.com/deepzz0/go-van/pkg/server"

	"google.golang.org/grpc"
)

type grpcOptsKey struct{}

// WithDialOpt grpc client option
func WithDialOpt(opts ...grpc.DialOption) server.DialOption {
	return func(opts *server.DialOptions) {
		if opts.Context == nil {
			opts.Context = context.Background()
		}
		opts.Context = context.WithValue(opts.Context, grpcOptsKey{}, opts)
	}
}

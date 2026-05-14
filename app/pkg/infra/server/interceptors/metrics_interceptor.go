package interceptors

import (
	"carduka/bidsvc/pkg/infra/monitoring/metrics"
	prometheusmetrics "carduka/bidsvc/pkg/infra/monitoring/metrics/prometheus"
	envutils "carduka/bidsvc/pkg/utils/env"
	"context"

	"google.golang.org/grpc"
)

// grpcServerMetricsInterceptor adds metrics to gRPC requests
func GrpcServerMetricsInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		// track metrics per handler if metrics are enabled, this has been defaulted to false by default
		if envutils.EnvBoolOr(metrics.EnvMetricsEnabled, false) {
			methodName := extractMethodName(info.FullMethod)
			defer prometheusmetrics.TrackMetrics(methodName, info.FullMethod)()
		}

		return handler(ctx, req)
	}
}

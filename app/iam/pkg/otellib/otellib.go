package otellib

import (
	"context"
	"fmt"
	"github.com/Nanhtu187/Online-Judge/app/iam/config"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

func newJaegerExporter(conf config.JaegerConfig) sdktrace.SpanExporter {
	if conf.Host == "" {
		conf.Host = "localhost"
	}
	if conf.Port == 0 {
		conf.Port = 6831
	}

	fmt.Printf("JAEGER URL: %s:%d\n", conf.Host, conf.Port)
	exporter, err := jaeger.New(
		jaeger.WithAgentEndpoint(
			jaeger.WithAgentHost(conf.Host),
			jaeger.WithAgentPort(fmt.Sprint(conf.Port)),
		),
	)
	if err != nil {
		panic(err)
	}
	return exporter
}

func newResource(name string, env string) *resource.Resource {
	r := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(name),
		semconv.ServiceVersionKey.String("v0.1.0"),
		attribute.String("environment", env),
	)
	return r
}

// InitOtel creates a tracer provider
func InitOtel(
	serviceName string, environment string,
	conf config.JaegerConfig,
) (traceProvider trace.TracerProvider, shutdown func()) {
	exporter := newJaegerExporter(conf)
	providersOpts := []sdktrace.TracerProviderOption{
		sdktrace.WithResource(newResource(serviceName, environment)),
		sdktrace.WithSampler(sdktrace.TraceIDRatioBased(conf.Ratio)),
	}
	if conf.Enabled {
		providersOpts = append(providersOpts, sdktrace.WithBatcher(exporter))
	}

	tracerProvider := sdktrace.NewTracerProvider(providersOpts...)

	shutdown = func() {
		_ = tracerProvider.Shutdown(context.Background())
	}
	return tracerProvider, shutdown
}

//func ignoreHealthCheckFilter() otelgrpc.Filter {
//	const healthLive = "/health.v1.HealthCheckService/Liveness"
//	const healthReady = "/health.v1.HealthCheckService/Readiness"
//
//	return func(info *otelgrpc.InterceptorInfo) bool {
//		if info.UnaryServerInfo.FullMethod == healthLive || info.UnaryServerInfo.FullMethod == healthReady {
//			return false
//		}
//
//		return true
//	}
//}

// UnaryServerInterceptor ...
func UnaryServerInterceptor(tp trace.TracerProvider) grpc.UnaryServerInterceptor {
	return otelgrpc.UnaryServerInterceptor(
		otelgrpc.WithTracerProvider(tp),
		otelgrpc.WithPropagators(propagation.TraceContext{}),
		//otelgrpc.WithInterceptorFilter(ignoreHealthCheckFilter()),
	)
}

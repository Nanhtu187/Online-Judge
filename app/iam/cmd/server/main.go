package main

import (
	"context"
	"fmt"
	"github.com/Nanhtu187/Online-Judge/app/iam/config"
	"github.com/Nanhtu187/Online-Judge/app/iam/pkg/errors"
	"github.com/Nanhtu187/Online-Judge/app/iam/pkg/grpclib"
	log "github.com/Nanhtu187/Online-Judge/app/iam/pkg/logger"
	"github.com/Nanhtu187/Online-Judge/app/iam/pkg/otellib"
	iam2 "github.com/Nanhtu187/Online-Judge/app/iam/service/iam"
	"github.com/Nanhtu187/Online-Judge/proto/rpc/iam/v1"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"sync"
	"syscall"
	"time"
)

func main() {
	errors.FinishNewErrors()

	rootCmd := cobra.Command{
		Use: "server",
	}
	rootCmd.AddCommand(
		startServerCommand(),
	)

	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
	}
}

func registerGRPCGateway(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) {
	_ = iam.RegisterIamServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
}

func startServer() {
	conf, err := config.Load()
	if err != nil {
		fmt.Println("ERROR:", err)
	}

	logger := config.NewLogger(conf.Log)
	db := conf.Database.MustConnect()
	rdb := conf.Redis.MustConnect()
	if err := rdb.Get(context.Background(), "test").Err(); err != nil {
		fmt.Println("Connect redis successful")
	} else {
		fmt.Println("Fail to connect redis")
	}

	tracerProvider, shutdown := otellib.InitOtel("journey-builder-service", "local", conf.Jaeger)
	defer shutdown()

	otel.SetTracerProvider(tracerProvider)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	grpcServer := grpc.NewServer(
		grpclib.ChainUnaryInterceptorIgnoreHealthCheck(
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_recovery.UnaryServerInterceptor(grpc_recovery.WithRecoveryHandler(recoveryHandlerFunc)),
			grpc_prometheus.UnaryServerInterceptor,
			otellib.UnaryServerInterceptor(tracerProvider),
			log.SetTraceInfoInterceptor(logger),
			errors.UnaryServerInterceptor,
		),
		grpc.ChainStreamInterceptor(
			grpc_recovery.StreamServerInterceptor(grpc_recovery.WithRecoveryHandler(recoveryHandlerFunc)),
			grpc_ctxtags.StreamServerInterceptor(),
			grpc_prometheus.StreamServerInterceptor,
			grpc_zap.StreamServerInterceptor(logger),
			//grpc_zap.PayloadStreamServerInterceptor(logger, loggingDecider),
		),
	)

	iamServer := iam2.InitServer(db, conf)
	iam.RegisterIamServiceServer(grpcServer, iamServer)
	grpc_prometheus.EnableHandlingTimeHistogram()
	grpc_prometheus.Register(grpcServer)

	startHTTPAndGRPCServers(*conf, grpcServer)
}

func startHTTPAndGRPCServers(conf config.Config, grpcServer *grpc.Server) {
	fmt.Println("GRPC:", conf.Server.GRPC.ListenString())
	fmt.Println("HTTP:", conf.Server.HTTP.ListenString())

	mux := runtime.NewServeMux(
		runtime.WithErrorHandler(errors.CustomerHTTPError),
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{}),
	)

	ctx := context.Background()
	grpcHost := conf.Server.GRPC.String()
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}

	registerGRPCGateway(ctx, mux, grpcHost, opts)

	httpMux := http.NewServeMux()
	httpMux.Handle("/metrics", promhttp.Handler())
	httpMux.Handle("/", mux)

	httpServer := &http.Server{
		Addr:    conf.Server.HTTP.ListenString(),
		Handler: httpMux,
	}

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()

		err := httpServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			panic(err)
		}
		fmt.Println("Shutdown HTTP server successfully")
	}()

	go func() {
		defer wg.Done()

		listener, err := net.Listen("tcp", conf.Server.GRPC.ListenString())
		if err != nil {
			panic(err)
		}

		err = grpcServer.Serve(listener)
		if err != nil {
			panic(err)
		}
		fmt.Println("Shutdown gRPC server successfully")
	}()

	//--------------------------------
	// Graceful Shutdown
	//--------------------------------
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	ctx = context.Background()
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	grpcServer.GracefulStop()
	err := httpServer.Shutdown(ctx)
	if err != nil {
		panic(err)
	}

	wg.Wait()
}

func startServerCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "start",
		Short: "start the server",
		Run: func(cmd *cobra.Command, args []string) {
			startServer()
		},
	}
}

func recoveryHandlerFunc(p interface{}) error {
	fmt.Println("stacktrace from panic:\n" + string(debug.Stack()))
	return status.Errorf(codes.Internal, "panic: %s", p)
}

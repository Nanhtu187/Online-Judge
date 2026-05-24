package cmd

import (
	"context"
	"log"
	"net"
	"net/http"

	"github.com/Nanhtu187/online-judge/src/apps/backend/iam/config"
	pbcommon "github.com/Nanhtu187/online-judge/src/packages/proto/gen/go/common"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func startServer(cmd *cobra.Command, args []string) {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	go func() {
		lis, err := net.Listen("tcp", cfg.Server.GrpcAddr())
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		s := grpc.NewServer()
		pbcommon.RegisterHealthServiceServer(s, &healthServer{})
		log.Printf("IAM gRPC server listening at %v", lis.Addr())
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve gRPC: %v", err)
		}
	}()

	ctx := context.Background()
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err = pbcommon.RegisterHealthServiceHandlerFromEndpoint(ctx, mux, cfg.Server.GrpcAddr(), opts)
	if err != nil {
		log.Fatalf("failed to register gateway: %v", err)
	}

	log.Printf("IAM REST server listening at %v", cfg.Server.HTTPAddr())
	if err := http.ListenAndServe(cfg.Server.HTTPAddr(), mux); err != nil {
		log.Fatalf("failed to serve REST: %v", err)
	}
}

func StartServerCommand() *cobra.Command {
	return &cobra.Command{
		Use: "start",
		Run: startServer,
	}
}

type healthServer struct {
	pbcommon.UnimplementedHealthServiceServer
}

func (s *healthServer) Check(ctx context.Context, in *pbcommon.HealthCheckRequest) (*pbcommon.HealthCheckResponse, error) {
	return &pbcommon.HealthCheckResponse{
		Status: pbcommon.HealthCheckResponse_SERVING,
	}, nil
}

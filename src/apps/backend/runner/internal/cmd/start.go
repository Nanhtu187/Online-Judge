package cmd

import (
	"context"
	"log"

	"github.com/Nanhtu187/online-judge/src/apps/backend/runner/config"
	"github.com/Nanhtu187/online-judge/src/apps/backend/runner/internal/adapter"
	"github.com/Nanhtu187/online-judge/src/apps/backend/runner/internal/compiler"
	"github.com/Nanhtu187/online-judge/src/apps/backend/runner/internal/docker"
	"github.com/Nanhtu187/online-judge/src/apps/backend/runner/internal/executor"
	"github.com/Nanhtu187/online-judge/src/apps/backend/runner/internal/handler"
	"github.com/Nanhtu187/online-judge/src/apps/backend/runner/internal/service"
	kfk "github.com/Nanhtu187/online-judge/src/packages/kafka"
	pb "github.com/Nanhtu187/online-judge/src/packages/proto/gen/go/server"
	"github.com/segmentio/kafka-go"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func runRunner() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	dAdapter, err := docker.NewAdapter()
	if err != nil {
		log.Fatalf("failed to init docker adapter: %v", err)
	}

	comp := compiler.NewDockerCompiler(dAdapter)
	exec := executor.NewDockerExecutor(dAdapter)

	// Init gRPC client with interceptor for internal key
	interceptor := func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		md := metadata.Pairs("x-internal-key", cfg.InternalKey)
		ctx = metadata.NewOutgoingContext(ctx, md)
		return invoker(ctx, method, req, reply, cc, opts...)
	}

	conn, err := grpc.Dial(cfg.ServerEndpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(interceptor),
	)
	if err != nil {
		log.Fatalf("failed to connect to server gRPC: %v", err)
	}
	defer conn.Close()
	serverClient := pb.NewOnlineJudgeServiceClient(conn)

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: cfg.Kafka.Brokers,
		Topic:   kfk.SubmissionRequestTopic,
		GroupID: "runner-group",
	})
	defer reader.Close()

	writer := &kafka.Writer{
		Addr:     kafka.TCP(cfg.Kafka.Brokers...),
		Topic:    kfk.SubmissionResultTopic,
		Balancer: &kafka.LeastBytes{},
	}
	defer writer.Close()

	resultWriter := adapter.NewKafkaResultWriter(writer)
	judgerSvc := service.NewJudgerService(comp, exec, serverClient, resultWriter)
	judgerHandler := handler.NewJudgerHandler(judgerSvc, reader)

	if err := judgerHandler.Start(context.Background()); err != nil {
		log.Fatalf("runner handler failed: %v", err)
	}
}

func StartJudgerCommand() *cobra.Command {
	return &cobra.Command{
		Use: "judger",
		Run: func(cmd *cobra.Command, args []string) {
			runRunner()
		},
	}
}

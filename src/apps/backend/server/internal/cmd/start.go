package cmd

import (
	"context"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/Nanhtu187/online-judge/src/apps/backend/server/config"
	"github.com/Nanhtu187/online-judge/src/packages/database"
	"github.com/Nanhtu187/online-judge/src/packages/iam"
	"github.com/Nanhtu187/online-judge/src/packages/logger"
	"github.com/Nanhtu187/online-judge/src/apps/backend/server/internal/handler"
	"github.com/Nanhtu187/online-judge/src/apps/backend/server/internal/repository"
	"github.com/Nanhtu187/online-judge/src/apps/backend/server/internal/service"
	common "github.com/Nanhtu187/online-judge/src/packages/proto/gen/go/common"
	pb "github.com/Nanhtu187/online-judge/src/packages/proto/gen/go/server"
	iampb "github.com/Nanhtu187/online-judge/src/packages/proto/gen/go/iam"
	"go.uber.org/zap"
	"os"
)

func startServer(cmd *cobra.Command, args []string) {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	l := logger.New(os.Getenv("GO_ENV"))
	defer l.Sync()

	// Database connection
	db, err := gorm.Open(mysql.Open(cfg.Database.DSN()), &gorm.Config{})
	if err != nil {
		l.Fatal("failed to connect to database", zap.Error(err))
	}
	dbProvider := database.NewProvider(db)
// Wire dependencies
problemRepo := repository.NewProblemRepository(dbProvider.(*database.Provider), l)
problemSvc := service.NewProblemService(problemRepo, dbProvider)

submissionRepo := repository.NewSubmissionRepository(dbProvider.(*database.Provider), l)
resultRepo := repository.NewSubmissionResultRepository(dbProvider.(*database.Provider), l)
submissionSvc := service.NewSubmissionService(submissionRepo, resultRepo, dbProvider, cfg.Kafka.Brokers)

// Init IAM client
iamConn, err := grpc.Dial(cfg.IAMEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
if err != nil {
	l.Fatal("failed to connect to iam service", zap.Error(err))
}
defer iamConn.Close()
iamClient := iampb.NewIamServiceClient(iamConn)
authMid := iam.NewAuthMiddleware(iamClient, cfg.InternalKey)
// Define permission map
perms := map[string]string{
	"/server.v1.OnlineJudgeService/UpsertProblem":   "PROBLEM_CREATE",
	"/server.v1.OnlineJudgeService/UpsertTestCases": "TEST_CASE_MANAGE",
	"/server.v1.OnlineJudgeService/SubmitCode":      "SUBMISSION_SUBMIT",
}

onlineJudgeHandler := handler.NewOnlineJudgeHandler(problemSvc, submissionSvc, authMid)

go func() {
	lis, err := net.Listen("tcp", cfg.Server.GrpcAddr())
	if err != nil {
		l.Fatal("failed to listen grpc", zap.Error(err))
	}

	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			authMid.UnaryInterceptor,
			authMid.AuthorizeInterceptor(perms),
		),
	)
	pb.RegisterOnlineJudgeServiceServer(s, onlineJudgeHandler)
		l.Info("Server gRPC server listening", zap.String("addr", lis.Addr().String()))
		if err := s.Serve(lis); err != nil {
			l.Fatal("failed to serve gRPC", zap.Error(err))
		}
	}()

	ctx := context.Background()
	// Set default database context for REST
	ctx = dbProvider.Readonly(ctx)

	mux := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				UseProtoNames:   true,
				EmitUnpopulated: true,
			},
			UnmarshalOptions: protojson.UnmarshalOptions{
				DiscardUnknown: true,
			},
		}),
	)
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	err = pb.RegisterOnlineJudgeServiceHandlerFromEndpoint(ctx, mux, cfg.Server.GrpcAddr(), opts)
	if err != nil {
		l.Fatal("failed to register online judge gateway", zap.Error(err))
	}

	l.Info("Server REST server listening", zap.String("addr", cfg.Server.HTTPAddr()))
	handler := allowCORS(mux)
	if err := http.ListenAndServe(cfg.Server.HTTPAddr(), handler); err != nil {
		l.Fatal("failed to serve REST", zap.Error(err))
	}
}

func allowCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization")
		if r.Method == "OPTIONS" {
			return
		}
		h.ServeHTTP(w, r)
	})
}

func StartServerCommand() *cobra.Command {
	return &cobra.Command{
		Use: "start",
		Run: startServer,
	}
}

type healthServer struct {
	common.UnimplementedHealthServiceServer
}

func (s *healthServer) Check(ctx context.Context, in *common.HealthCheckRequest) (*common.HealthCheckResponse, error) {
	return &common.HealthCheckResponse{
		Status: common.HealthCheckResponse_SERVING,
	}, nil
}

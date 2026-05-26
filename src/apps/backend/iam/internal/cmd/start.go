package cmd

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/Nanhtu187/online-judge/src/apps/backend/iam/config"
	"github.com/Nanhtu187/online-judge/src/apps/backend/iam/internal/handler"
	"github.com/Nanhtu187/online-judge/src/apps/backend/iam/internal/repository"
	"github.com/Nanhtu187/online-judge/src/apps/backend/iam/internal/service"
	"github.com/Nanhtu187/online-judge/src/packages/database"
	"github.com/Nanhtu187/online-judge/src/packages/logger"
	pb "github.com/Nanhtu187/online-judge/src/packages/proto/gen/go/iam"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func startServer(cmd *cobra.Command, args []string) {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	l := logger.New(os.Getenv("GO_ENV"))
	defer l.Sync()

	db, err := gorm.Open(mysql.Open(cfg.Database.DSN()), &gorm.Config{})
	if err != nil {
		l.Fatal("failed to connect to database", zap.Error(err))
	}
	dbProvider := database.NewProvider(db)

	// Wire dependencies
	iamRepo := repository.NewIamRepository(dbProvider.(*database.Provider), l)
	iamSvc := service.NewIamService(iamRepo, dbProvider, l, cfg.JWTSecret)
	iamHandler := handler.NewIamHandler(iamSvc, cfg.JWTSecret)

	go func() {
		lis, err := net.Listen("tcp", cfg.Server.GrpcAddr())
		if err != nil {
			l.Fatal("failed to listen grpc", zap.Error(err))
		}

		s := grpc.NewServer()
		pb.RegisterIamServiceServer(s, iamHandler)

		l.Info("IAM gRPC server listening", zap.String("addr", cfg.Server.GrpcAddr()))
		if err := s.Serve(lis); err != nil {
			l.Fatal("failed to serve grpc", zap.Error(err))
		}
	}()

	ctx := context.Background()
	mux := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				UseProtoNames:   true,
				EmitUnpopulated: true,
			},
		}),
	)
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	err = pb.RegisterIamServiceHandlerFromEndpoint(ctx, mux, cfg.Server.GrpcAddr(), opts)
	if err != nil {
		l.Fatal("failed to register iam gateway", zap.Error(err))
	}

	l.Info("IAM REST server listening", zap.String("addr", cfg.Server.HTTPAddr()))
	handler := allowCORS(mux)
	if err := http.ListenAndServe(cfg.Server.HTTPAddr(), handler); err != nil {
		l.Fatal("failed to serve rest", zap.Error(err))
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

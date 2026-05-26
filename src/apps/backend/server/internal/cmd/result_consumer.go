package cmd

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Nanhtu187/online-judge/src/apps/backend/server/config"
	"github.com/Nanhtu187/online-judge/src/apps/backend/server/internal/handler"
	"github.com/Nanhtu187/online-judge/src/apps/backend/server/internal/repository"
	"github.com/Nanhtu187/online-judge/src/apps/backend/server/internal/service"
	"github.com/Nanhtu187/online-judge/src/packages/database"
	kfk "github.com/Nanhtu187/online-judge/src/packages/kafka"
	"github.com/Nanhtu187/online-judge/src/packages/logger"
	"github.com/segmentio/kafka-go"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func startResultConsumer(cmd *cobra.Command, args []string) {
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

	submissionRepo := repository.NewSubmissionRepository(dbProvider.(*database.Provider), l)
	resultRepo := repository.NewSubmissionResultRepository(dbProvider.(*database.Provider), l)

	consumerSvc := service.NewResultConsumerService(submissionRepo, resultRepo, dbProvider)

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: cfg.Kafka.Brokers,
		Topic:   kfk.SubmissionResultTopic,
		GroupID: "result-consumer-group",
	})
	defer reader.Close()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	consumerHandler := handler.NewResultConsumerHandler(consumerSvc, reader, l)
	consumerHandler.Start(ctx)
}

func StartResultConsumerCommand() *cobra.Command {
	return &cobra.Command{
		Use: "result-consumer",
		Run: startResultConsumer,
	}
}

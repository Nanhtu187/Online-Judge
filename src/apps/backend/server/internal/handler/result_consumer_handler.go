package handler

import (
	"context"
	"encoding/json"

	"github.com/Nanhtu187/online-judge/src/apps/backend/server/internal/service"
	kfk "github.com/Nanhtu187/online-judge/src/packages/kafka"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

type ResultConsumerHandler interface {
	Start(ctx context.Context)
}

type resultConsumerHandler struct {
	svc    service.ResultConsumerService
	reader *kafka.Reader
	logger *zap.Logger
}

func NewResultConsumerHandler(svc service.ResultConsumerService, reader *kafka.Reader, logger *zap.Logger) ResultConsumerHandler {
	return &resultConsumerHandler{
		svc:    svc,
		reader: reader,
		logger: logger,
	}
}

func (h *resultConsumerHandler) Start(ctx context.Context) {
	h.logger.Info("Result consumer started, consuming results...")

	for {
		m, err := h.reader.ReadMessage(ctx)
		if err != nil {
			h.logger.Error("error reading message", zap.Error(err))
			continue
		}

		var event kfk.ResultEvent
		if err := json.Unmarshal(m.Value, &event); err != nil {
			h.logger.Error("error unmarshaling event", zap.Error(err))
			continue
		}

		if err := h.svc.ProcessResultEvent(ctx, &event); err != nil {
			h.logger.Error("failed to process result event", zap.Error(err))
		} else {
			h.logger.Info("processed result event", zap.String("submission_id", event.SubmissionID))
		}
	}
}

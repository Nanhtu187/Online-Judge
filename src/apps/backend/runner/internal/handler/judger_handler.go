package handler

import (
	"context"
	"encoding/json"
	"log"

	"github.com/Nanhtu187/online-judge/src/apps/backend/runner/internal/service"
	kfk "github.com/Nanhtu187/online-judge/src/packages/kafka"
	"github.com/segmentio/kafka-go"
)

type JudgerHandler interface {
	Start(ctx context.Context) error
}

type judgerHandler struct {
	svc    service.JudgerService
	reader *kafka.Reader
}

func NewJudgerHandler(svc service.JudgerService, reader *kafka.Reader) JudgerHandler {
	return &judgerHandler{
		svc:    svc,
		reader: reader,
	}
}

func (h *judgerHandler) Start(ctx context.Context) error {
	log.Println("Runner service started, consuming events...")

	for {
		m, err := h.reader.ReadMessage(ctx)
		if err != nil {
			return err
		}

		var event kfk.SubmissionEvent
		if err := json.Unmarshal(m.Value, &event); err != nil {
			log.Printf("error unmarshaling event: %v", err)
			continue
		}

		if err := h.svc.Judge(ctx, &event); err != nil {
			log.Printf("error judging submission %s: %v", event.SubmissionID, err)
		}
	}
}

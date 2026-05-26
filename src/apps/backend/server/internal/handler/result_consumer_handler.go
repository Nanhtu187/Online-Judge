package handler

import (
	"context"
	"encoding/json"
	"sync"

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

	msgChan := make(chan kafka.Message, 500)
	var wg sync.WaitGroup

	// 1. Fetcher Goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer close(msgChan)
		for {
			m, err := h.reader.FetchMessage(ctx)
			if err != nil {
				// context.Canceled is expected on shutdown
				if ctx.Err() == nil {
					h.logger.Error("error fetching message", zap.Error(err))
				}
				return
			}
			msgChan <- m
		}
	}()

	// 2. Worker Goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			// Block for the first message
			m, ok := <-msgChan
			if !ok {
				return
			}

			// Start a drain cycle
			batch := []kafka.Message{m}
			
			// Non-blocking drain remaining messages in buffer
			draining := true
			for len(batch) < 100 && draining {
				select {
				case next, ok := <-msgChan:
					if !ok {
						draining = false
					} else {
						batch = append(batch, next)
					}
				default:
					draining = false
				}
			}

			h.processBatch(context.Background(), batch)
		}
	}()

	wg.Wait()
	h.logger.Info("Result consumer stopped gracefully")
}

func (h *resultConsumerHandler) processBatch(ctx context.Context, messages []kafka.Message) {
	if len(messages) == 0 {
		return
	}

	var events []*kfk.ResultEvent
	for _, m := range messages {
		var event kfk.ResultEvent
		if err := json.Unmarshal(m.Value, &event); err != nil {
			h.logger.Error("error unmarshaling event", zap.Error(err))
			continue
		}
		events = append(events, &event)
	}

	if err := h.svc.ProcessResultEvents(ctx, events); err != nil {
		h.logger.Error("failed to process batch events", zap.Error(err))
		// Note: In production, we should implement a retry policy here 
		// because if we fail to process, we shouldn't commit.
		// For now, we log and continue to avoid blocking the queue.
	} else {
		// Only commit if DB sync was successful
		if err := h.reader.CommitMessages(ctx, messages...); err != nil {
			h.logger.Error("failed to commit messages", zap.Error(err))
		} else {
			h.logger.Info("processed and committed batch", zap.Int("count", len(messages)))
		}
	}
}

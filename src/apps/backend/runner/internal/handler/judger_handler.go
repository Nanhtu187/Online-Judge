package handler

import (
	"context"
	"encoding/json"
	"log"
	"sync"

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

	const numWorkers = 5
	msgChan := make(chan kafka.Message, numWorkers*2)
	var wg sync.WaitGroup

	// Start Workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for m := range msgChan {
				var event kfk.SubmissionEvent
				if err := json.Unmarshal(m.Value, &event); err != nil {
					log.Printf("[Worker %d] error unmarshaling event: %v", workerID, err)
				} else {
					if err := h.svc.Judge(context.Background(), &event); err != nil {
						log.Printf("[Worker %d] error judging submission %s: %v", workerID, event.SubmissionID, err)
					}
				}

				// Commit message after processing (success or fail)
				if err := h.reader.CommitMessages(context.Background(), m); err != nil {
					log.Printf("[Worker %d] failed to commit message for submission %s: %v", workerID, event.SubmissionID, err)
				}
			}
		}(i)
	}

	// Fetcher loop
	for {
		m, err := h.reader.FetchMessage(ctx)
		if err != nil {
			// Expected on shutdown when context is cancelled
			if ctx.Err() != nil {
				log.Println("Context cancelled, stopping fetcher...")
				break
			}
			log.Printf("error fetching message: %v", err)
			continue
		}
		msgChan <- m
	}

	close(msgChan)
	log.Println("Waiting for workers to finish...")
	wg.Wait()
	log.Println("Runner service stopped gracefully.")

	return nil
}


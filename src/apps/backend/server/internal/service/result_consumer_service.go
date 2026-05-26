package service

import (
	"context"

	"github.com/Nanhtu187/online-judge/src/apps/backend/server/internal/models"
	"github.com/Nanhtu187/online-judge/src/apps/backend/server/internal/repository"
	"github.com/Nanhtu187/online-judge/src/packages/database"
	kfk "github.com/Nanhtu187/online-judge/src/packages/kafka"
)

type ResultConsumerService interface {
	ProcessResultEvents(ctx context.Context, events []*kfk.ResultEvent) error
}

type resultConsumerService struct {
	submissionRepo repository.SubmissionRepository
	resultRepo     repository.SubmissionResultRepository
	provider       database.IProvider
}

func NewResultConsumerService(
	submissionRepo repository.SubmissionRepository,
	resultRepo repository.SubmissionResultRepository,
	provider database.IProvider,
) ResultConsumerService {
	return &resultConsumerService{
		submissionRepo: submissionRepo,
		resultRepo:     resultRepo,
		provider:       provider,
	}
}

func (s *resultConsumerService) ProcessResultEvents(ctx context.Context, events []*kfk.ResultEvent) error {
	if len(events) == 0 {
		return nil
	}

	return s.provider.Transact(ctx, func(ctx context.Context) error {
		var tcBatch []*models.TestCaseResult
		submissionStatuses := make(map[string]models.SubmissionStatus)

		for _, event := range events {
			status := models.SubmissionStatus(event.Status)

			if event.TestCaseID != "" {
				tcBatch = append(tcBatch, &models.TestCaseResult{
					SubmissionID: event.SubmissionID,
					TestCaseID:   event.TestCaseID,
					Status:       status,
					ActualOutput: event.Output,
				})

				// If we haven't seen a final status yet, mark as RUNNING
				if _, ok := submissionStatuses[event.SubmissionID]; !ok {
					submissionStatuses[event.SubmissionID] = models.StatusRunning
				}
			} else {
				// This is a final summary event, it overwrites any RUNNING status
				submissionStatuses[event.SubmissionID] = status
			}
		}

		// 1. Bulk Insert Test Case Results
		if len(tcBatch) > 0 {
			if err := s.resultRepo.CreateBatch(ctx, tcBatch); err != nil {
				return err
			}
		}

		// 2. Bulk Update Submission Statuses
		for subID, status := range submissionStatuses {
			if err := s.submissionRepo.UpdateStatus(ctx, subID, status); err != nil {
				return err
			}
		}

		return nil
	})
}

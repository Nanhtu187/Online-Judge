package service

import (
	"context"

	"github.com/Nanhtu187/online-judge/src/apps/backend/server/internal/models"
	"github.com/Nanhtu187/online-judge/src/apps/backend/server/internal/repository"
	"github.com/Nanhtu187/online-judge/src/packages/database"
	kfk "github.com/Nanhtu187/online-judge/src/packages/kafka"
)

type ResultConsumerService interface {
	ProcessResultEvent(ctx context.Context, event *kfk.ResultEvent) error
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

func (s *resultConsumerService) ProcessResultEvent(ctx context.Context, event *kfk.ResultEvent) error {
	status := models.SubmissionStatus(event.Status)

	return s.provider.Transact(ctx, func(ctx context.Context) error {
		if event.TestCaseID != "" {
			result := &models.TestCaseResult{
				SubmissionID: event.SubmissionID,
				TestCaseID:   event.TestCaseID,
				Status:       status,
				ActualOutput: event.Output,
			}
			if err := s.resultRepo.Create(ctx, result); err != nil {
				return err
			}

			// If it's a test case result, keep the submission as RUNNING
			return s.submissionRepo.UpdateStatus(ctx, event.SubmissionID, models.StatusRunning)
		}

		// If TestCaseID is empty, it's a summary event (final result or compile error)
		return s.submissionRepo.UpdateStatus(ctx, event.SubmissionID, status)
	})
}

package service

import (
	"context"

	"github.com/Nanhtu187/online-judge/src/apps/backend/server/internal/models"
	"github.com/Nanhtu187/online-judge/src/apps/backend/server/internal/repository"
	"github.com/Nanhtu187/online-judge/src/packages/database"
	pb "github.com/Nanhtu187/online-judge/src/packages/proto/gen/go/server"
)

type ProblemService interface {
	UpsertProblem(ctx context.Context, req *pb.UpsertProblemRequest) (*pb.UpsertProblemResponse, error)
	UpsertTestCases(ctx context.Context, req *pb.UpsertTestCasesRequest) (*pb.UpsertTestCasesResponse, error)
	GetProblem(ctx context.Context, req *pb.GetProblemRequest) (*pb.GetProblemResponse, error)
	ListProblems(ctx context.Context, req *pb.ListProblemsRequest) (*pb.ListProblemsResponse, error)
	ListTestCases(ctx context.Context, req *pb.ListTestCasesRequest) (*pb.ListTestCasesResponse, error)
}

func (s *problemService) ListProblems(ctx context.Context, req *pb.ListProblemsRequest) (*pb.ListProblemsResponse, error) {
	ctx = s.provider.Readonly(ctx)
	page := int(req.Page)
	if page == 0 {
		page = 1
	}
	pageSize := int(req.PageSize)
	if pageSize == 0 {
		pageSize = 10
	}

	problems, err := s.repo.ListProblems(ctx, page, pageSize)
	if err != nil {
		return nil, err
	}

	var pbProblems []*pb.ProblemSummary
	for _, p := range problems {
		var tags []string
		for _, t := range p.Tags {
			tags = append(tags, t.Name)
		}
		pbProblems = append(pbProblems, &pb.ProblemSummary{
			Id:         p.ID,
			Title:      p.Title,
			Difficulty: pb.Difficulty(pb.Difficulty_value["DIFFICULTY_"+p.Difficulty]),
			Tags:       tags,
		})
	}

	return &pb.ListProblemsResponse{Problems: pbProblems}, nil
}

func (s *problemService) ListTestCases(ctx context.Context, req *pb.ListTestCasesRequest) (*pb.ListTestCasesResponse, error) {
	ctx = s.provider.Readonly(ctx)
	testCases, err := s.repo.ListTestCases(ctx, req.ProblemId, nil)
	if err != nil {
		return nil, err
	}

	var pbTestCases []*pb.TestCase
	for _, tc := range testCases {
		pbTestCases = append(pbTestCases, &pb.TestCase{
			Id:             tc.ID,
			Input:          tc.Input,
			ExpectedOutput: tc.ExpectedOutput,
			IsSample:       tc.IsSample,
		})
	}

	return &pb.ListTestCasesResponse{TestCases: pbTestCases}, nil
}

type problemService struct {
	repo     repository.ProblemRepository
	provider database.IProvider
}

func NewProblemService(repo repository.ProblemRepository, provider database.IProvider) ProblemService {
	return &problemService{repo: repo, provider: provider}
}

func (s *problemService) UpsertProblem(ctx context.Context, req *pb.UpsertProblemRequest) (*pb.UpsertProblemResponse, error) {
	difficulty := pb.Difficulty_name[int32(req.Difficulty)]
	if len(difficulty) > len("DIFFICULTY_") {
		difficulty = difficulty[len("DIFFICULTY_"):]
	} else {
		difficulty = "EASY"
	}

	problem := &models.Problem{
		ID:           req.Id,
		Title:        req.Title,
		Content:      req.Content,
		InputFormat:  req.InputFormat,
		OutputFormat: req.OutputFormat,
		TimeLimit:    req.TimeLimit,
		MemoryLimit:  req.MemoryLimit,
		Difficulty:   difficulty,
	}

	err := s.provider.Transact(ctx, func(ctx context.Context) error {
		if problem.ID != "" {
			if err := s.repo.UpdateProblem(ctx, problem); err != nil {
				return err
			}
		} else {
			if err := s.repo.CreateProblem(ctx, problem); err != nil {
				return err
			}
		}

		// Sync Tags
		var tags []*models.Tag
		for _, tagName := range req.Tags {
			tag, err := s.repo.GetTagByName(ctx, tagName)
			if err != nil {
				// Create new tag if not exists
				tag = &models.Tag{Name: tagName}
				if err := s.repo.CreateTag(ctx, tag); err != nil {
					return err
				}
			}
			tags = append(tags, tag)
		}
		return s.repo.ReplaceTags(ctx, problem.ID, tags)
	})
	if err != nil {
		return nil, err
	}

	return &pb.UpsertProblemResponse{
		Id: problem.ID,
	}, nil
}

func (s *problemService) UpsertTestCases(ctx context.Context, req *pb.UpsertTestCasesRequest) (*pb.UpsertTestCasesResponse, error) {
	err := s.provider.Transact(ctx, func(ctx context.Context) error {
		if err := s.repo.DeleteTestCases(ctx, req.ProblemId); err != nil {
			return err
		}
		for _, tc := range req.TestCases {
			if err := s.repo.CreateTestCase(ctx, &models.TestCase{
				ProblemID:      req.ProblemId,
				Input:          tc.Input,
				ExpectedOutput: tc.ExpectedOutput,
				IsSample:       tc.IsSample,
			}); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &pb.UpsertTestCasesResponse{Success: true}, nil
}

func (s *problemService) GetProblem(ctx context.Context, req *pb.GetProblemRequest) (*pb.GetProblemResponse, error) {
	ctx = s.provider.Readonly(ctx)
	problem, err := s.repo.GetByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	isSample := true
	testCases, err := s.repo.ListTestCases(ctx, req.Id, &isSample)
	if err != nil {
		return nil, err
	}

	var pbTestCases []*pb.TestCase
	for _, tc := range testCases {
		pbTestCases = append(pbTestCases, &pb.TestCase{
			Id:             tc.ID,
			Input:          tc.Input,
			ExpectedOutput: tc.ExpectedOutput,
			IsSample:       tc.IsSample,
		})
	}

	var tags []string
	for _, t := range problem.Tags {
		tags = append(tags, t.Name)
	}

	return &pb.GetProblemResponse{
		Problem: &pb.Problem{
			Id:           problem.ID,
			Title:        problem.Title,
			Content:      problem.Content,
			InputFormat:  problem.InputFormat,
			OutputFormat: problem.OutputFormat,
			TestCases:    pbTestCases,
			TimeLimit:    problem.TimeLimit,
			MemoryLimit:  problem.MemoryLimit,
			Difficulty:   pb.Difficulty(pb.Difficulty_value["DIFFICULTY_"+problem.Difficulty]),
			Tags:         tags,
		},
	}, nil
}

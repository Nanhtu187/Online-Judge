package service

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/Nanhtu187/online-judge/src/apps/backend/runner/internal/compiler"
	"github.com/Nanhtu187/online-judge/src/apps/backend/runner/internal/executor"
	kfk "github.com/Nanhtu187/online-judge/src/packages/kafka"
	pb "github.com/Nanhtu187/online-judge/src/packages/proto/gen/go/server"
)

type JudgerService interface {
	Judge(ctx context.Context, event *kfk.SubmissionEvent) error
}

type judgerService struct {
	compiler     compiler.Compiler
	executor     executor.Executor
	serverClient pb.OnlineJudgeServiceClient
	resultWriter ResultWriter
}

type ResultWriter interface {
	WriteResult(ctx context.Context, submissionID, testCaseID, status, output string) error
}

func NewJudgerService(
	comp compiler.Compiler,
	exec executor.Executor,
	client pb.OnlineJudgeServiceClient,
	writer ResultWriter,
) JudgerService {
	return &judgerService{
		compiler:     comp,
		executor:     exec,
		serverClient: client,
		resultWriter: writer,
	}
}

func (s *judgerService) Judge(ctx context.Context, event *kfk.SubmissionEvent) error {
	log.Printf("Judging submission: %s, lang: %s", event.SubmissionID, event.Language)

	// 1. Compile
	var binary []byte
	var err error
	if s.compiler.IsCompileLanguage(event.Language) {
		binary, err = s.compiler.Compile(ctx, event.Language, event.CodeContent)
		if err != nil {
			log.Printf("Compilation error for %s: %v", event.SubmissionID, err)
			return s.resultWriter.WriteResult(ctx, event.SubmissionID, "", "COMPILE_ERROR", err.Error())
		}
	} else {
		binary = []byte(event.CodeContent)
	}

	// 2. Fetch Problem and TestCases
	problemResp, err := s.serverClient.GetProblem(ctx, &pb.GetProblemRequest{Id: event.ProblemID})
	if err != nil {
		log.Printf("Failed to fetch problem %s: %v", event.ProblemID, err)
		return s.resultWriter.WriteResult(ctx, event.SubmissionID, "", "SYSTEM_ERROR", "Failed to fetch problem details")
	}
	problem := problemResp.Problem

	resp, err := s.serverClient.ListTestCases(ctx, &pb.ListTestCasesRequest{ProblemId: event.ProblemID})
	if err != nil {
		log.Printf("Failed to fetch testcases for %s: %v", event.SubmissionID, err)
		return s.resultWriter.WriteResult(ctx, event.SubmissionID, "", "SYSTEM_ERROR", "Failed to fetch testcases")
	}

	if len(resp.TestCases) == 0 {
		return s.resultWriter.WriteResult(ctx, event.SubmissionID, "", "SYSTEM_ERROR", "No testcases found")
	}

	// 3. Filter TestCases if needed
	testCases := resp.TestCases
	if event.SubmissionType == "TEST" {
		var samples []*pb.TestCase
		for _, tc := range resp.TestCases {
			if tc.IsSample {
				samples = append(samples, tc)
			}
		}
		testCases = samples
	}

	if len(testCases) == 0 {
		return s.resultWriter.WriteResult(ctx, event.SubmissionID, "", "SYSTEM_ERROR", "No matching testcases found")
	}

	// 4. Run and Judge
	finalStatus := "ACCEPTED"
	memoryLimit := int64(problem.MemoryLimit) * 1024 * 1024
	timeLimitNano := int64(problem.TimeLimit) * 1000000

	for i, tc := range testCases {
		log.Printf("Running testcase %d/%d for %s", i+1, len(testCases), event.SubmissionID)
		output, err := s.executor.Run(ctx, event.Language, binary, tc.Input, memoryLimit, timeLimitNano)
		
		tcStatus := "ACCEPTED"
		tcOutput := output

		if err != nil {
			tcStatus = "RUNTIME_ERROR"
			if strings.Contains(err.Error(), "limit exceeded") {
				tcStatus = "LIMIT_EXCEEDED"
			}
			tcOutput = fmt.Sprintf("Testcase %d: %v", i+1, err)
			
			if finalStatus == "ACCEPTED" {
				finalStatus = tcStatus
			}
		} else if strings.TrimSpace(output) != strings.TrimSpace(tc.ExpectedOutput) {
			tcStatus = "WRONG_ANSWER"
			tcOutput = fmt.Sprintf("Testcase %d: Mismatch\nExpected: %s\nActual: %s", i+1, tc.ExpectedOutput, output)
			
			if finalStatus == "ACCEPTED" {
				finalStatus = tcStatus
			}
		}

		s.resultWriter.WriteResult(ctx, event.SubmissionID, tc.Id, tcStatus, tcOutput)
	}

	return s.resultWriter.WriteResult(ctx, event.SubmissionID, "", finalStatus, "Judging completed")
}

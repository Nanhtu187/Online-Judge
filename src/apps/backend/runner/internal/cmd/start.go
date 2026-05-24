package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/Nanhtu187/online-judge/src/apps/backend/runner/config"
	"github.com/Nanhtu187/online-judge/src/apps/backend/runner/internal/compiler"
	"github.com/Nanhtu187/online-judge/src/apps/backend/runner/internal/docker"
	"github.com/Nanhtu187/online-judge/src/apps/backend/runner/internal/executor"
	kfk "github.com/Nanhtu187/online-judge/src/packages/kafka"
	pb "github.com/Nanhtu187/online-judge/src/packages/proto/gen/go/server"
	"github.com/segmentio/kafka-go"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func runRunner() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	adapter, err := docker.NewAdapter()
	if err != nil {
		log.Fatalf("failed to init docker adapter: %v", err)
	}

	comp := compiler.NewDockerCompiler(adapter)
	exec := executor.NewDockerExecutor(adapter)

	// Init gRPC client
	conn, err := grpc.Dial(cfg.ServerEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to server gRPC: %v", err)
	}
	defer conn.Close()
	client := pb.NewOnlineJudgeServiceClient(conn)

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: cfg.Kafka.Brokers,
		Topic:   kfk.SubmissionRequestTopic,
		GroupID: "runner-group",
	})
	defer reader.Close()

	writer := &kafka.Writer{
		Addr:     kafka.TCP(cfg.Kafka.Brokers...),
		Topic:    kfk.SubmissionResultTopic,
		Balancer: &kafka.LeastBytes{},
	}
	defer writer.Close()

	log.Println("Runner service started, consuming events...")

	for {
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Printf("error reading message: %v", err)
			continue
		}

		var event kfk.SubmissionEvent
		if err := json.Unmarshal(m.Value, &event); err != nil {
			log.Printf("error unmarshaling event: %v", err)
			continue
		}

		log.Printf("Received submission: %s, lang: %s", event.SubmissionID, event.Language)

		// 1. Compile
		var binary []byte
		if comp.IsCompileLanguage(event.Language) {
			log.Printf("Compiling submission: %s", event.SubmissionID)
			binary, err = comp.Compile(context.Background(), event.Language, event.CodeContent)
			if err != nil {
				log.Printf("Compilation error for %s: %v", event.SubmissionID, err)
				pushResult(writer, event.SubmissionID, "", "COMPILE_ERROR", err.Error())
				continue
			}
			log.Printf("Compilation successful for %s, binary size: %d bytes", event.SubmissionID, len(binary))
		} else {
			binary = []byte(event.CodeContent)
		}

		// 2. Fetch Problem and TestCases
		problemResp, err := client.GetProblem(context.Background(), &pb.GetProblemRequest{Id: event.ProblemID})
		if err != nil {
			log.Printf("Failed to fetch problem %s: %v", event.ProblemID, err)
			pushResult(writer, event.SubmissionID, "", "SYSTEM_ERROR", "Failed to fetch problem details")
			continue
		}
		problem := problemResp.Problem

		resp, err := client.ListTestCases(context.Background(), &pb.ListTestCasesRequest{ProblemId: event.ProblemID})
		if err != nil {
			log.Printf("Failed to fetch testcases for %s: %v", event.SubmissionID, err)
			pushResult(writer, event.SubmissionID, "", "SYSTEM_ERROR", "Failed to fetch testcases")
			continue
		}

		if len(resp.TestCases) == 0 {
			log.Printf("No testcases found for problem %s", event.ProblemID)
			pushResult(writer, event.SubmissionID, "", "SYSTEM_ERROR", "No testcases found")
			continue
		}

		// 3. Filter TestCases if needed
		testCases := resp.TestCases
		if event.SubmissionType == "TEST" {
			log.Printf("Filtering for sample testcases only (SubmissionType: TEST)")
			var samples []*pb.TestCase
			for _, tc := range resp.TestCases {
				if tc.IsSample {
					samples = append(samples, tc)
				}
			}
			testCases = samples
		}

		if len(testCases) == 0 {
			log.Printf("No matching testcases for submission %s", event.SubmissionID)
			pushResult(writer, event.SubmissionID, "", "SYSTEM_ERROR", "No matching testcases found")
			continue
		}

		// 4. Run and Judge
		finalStatus := "ACCEPTED"

		// Limits
		memoryLimit := int64(problem.MemoryLimit) * 1024 * 1024
		timeLimitNano := int64(problem.TimeLimit) * 1000000

		for i, tc := range testCases {
			log.Printf("Running testcase %d/%d for %s (Limits: %dms, %dMB)", i+1, len(testCases), event.SubmissionID, problem.TimeLimit, problem.MemoryLimit)
			output, err := exec.Run(context.Background(), event.Language, binary, tc.Input, memoryLimit, timeLimitNano)
			
			tcStatus := "ACCEPTED"
			tcOutput := output

			if err != nil {
				log.Printf("Execution error for %s on testcase %d: %v", event.SubmissionID, i+1, err)
				tcStatus = "RUNTIME_ERROR"
				if strings.Contains(err.Error(), "limit exceeded") {
					tcStatus = "LIMIT_EXCEEDED"
				}
				tcOutput = fmt.Sprintf("Testcase %d: %v", i+1, err)
				
				if finalStatus == "ACCEPTED" {
					finalStatus = tcStatus
				}
			} else if strings.TrimSpace(output) != strings.TrimSpace(tc.ExpectedOutput) {
				log.Printf("Wrong answer for %s on testcase %d", event.SubmissionID, i+1)
				tcStatus = "WRONG_ANSWER"
				tcOutput = fmt.Sprintf("Testcase %d: Mismatch\nExpected: %s\nActual: %s", i+1, tc.ExpectedOutput, output)
				
				if finalStatus == "ACCEPTED" {
					finalStatus = tcStatus
				}
			}

			pushResult(writer, event.SubmissionID, tc.Id, tcStatus, tcOutput)
		}

		log.Printf("Judging complete for %s: %s", event.SubmissionID, finalStatus)
		pushResult(writer, event.SubmissionID, "", finalStatus, "Judging completed")
	}
}

func pushResult(writer *kafka.Writer, submissionID, testCaseID, status, output string) {
	result := kfk.ResultEvent{
		SubmissionID: submissionID,
		TestCaseID:   testCaseID,
		Status:       status,
		Output:       output,
	}
	resultJSON, _ := json.Marshal(result)
	err := writer.WriteMessages(context.Background(), kafka.Message{
		Key:   []byte(submissionID),
		Value: resultJSON,
	})
	if err != nil {
		log.Printf("error pushing result for %s: %v", submissionID, err)
	}
}

func StartJudgerCommand() *cobra.Command {
	return &cobra.Command{
		Use: "judger",
		Run: func(cmd *cobra.Command, args []string) {
			runRunner()
		},
	}
}

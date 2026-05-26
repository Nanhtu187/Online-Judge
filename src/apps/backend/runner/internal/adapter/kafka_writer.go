package adapter

import (
	"context"
	"encoding/json"
	"log"

	kfk "github.com/Nanhtu187/online-judge/src/packages/kafka"
	"github.com/segmentio/kafka-go"
)

type kafkaResultWriter struct {
	writer *kafka.Writer
}

func NewKafkaResultWriter(writer *kafka.Writer) *kafkaResultWriter {
	return &kafkaResultWriter{writer: writer}
}

func (w *kafkaResultWriter) WriteResult(ctx context.Context, submissionID, testCaseID, status, output string) error {
	result := kfk.ResultEvent{
		SubmissionID: submissionID,
		TestCaseID:   testCaseID,
		Status:       status,
		Output:       output,
	}
	resultJSON, _ := json.Marshal(result)
	err := w.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(submissionID),
		Value: resultJSON,
	})
	if err != nil {
		log.Printf("error pushing result for %s: %v", submissionID, err)
	}
	return err
}

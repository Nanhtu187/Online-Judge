package kafka

const (
	SubmissionRequestTopic = "submission-requests"
	SubmissionResultTopic  = "submission-results"
)
type SubmissionEvent struct {
	SubmissionID   string `json:"submission_id"`
	ProblemID      string `json:"problem_id"`
	CodeContent    string `json:"code_content"`
	Language       string `json:"language"`
	SubmissionType string `json:"submission_type"`
}

type ResultEvent struct {
	SubmissionID string `json:"submission_id"`
	TestCaseID   string `json:"test_case_id,omitempty"`
	Status       string `json:"status"`
	Output       string `json:"output"`
}

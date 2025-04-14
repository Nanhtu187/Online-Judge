package server

type UpsertProblemRequest struct {
	id          int32  `json:"problem_id"`
	title       string `json:"title"`
	description string `json:"description"`
}

type UpsertSubmissionRequest struct {
	problemId int32  `json:"problem_id"`
	contestId int32  `json:"contest_id"`
	userId    int32  `json:"user_id"`
	language  string `json:"language_id"`
	code      string `json:"code"`
}

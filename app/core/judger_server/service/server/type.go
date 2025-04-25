package server

type UpsertProblemRequest struct {
	ProblemId   int32  `json:"problem_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type UpsertSubmissionRequest struct {
	ProblemId int32  `json:"problem_id"`
	ContestId int32  `json:"contest_id"`
	UserId    int32  `json:"user_id"`
	Language  string `json:"language_id"`
	Code      string `json:"code"`
}

type GetProblemRequest struct {
	ProblemId int32 `json:"problem_id"`
}

type ProblemDetail struct {
	ProblemId   int32  `json:"problem_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

type GetListProblemRequest struct {
	ContestId int32  `json:"contest_id"`
	CreatedBy int32  `json:"created_by"`
	Keywords  string `json:"keywords"`
	Limit     int32  `json:"limit"`
	Offset    int32  `json:"offset"`
}

type ProblemPreview struct {
	ProblemId int32  `json:"problem_id"`
	Title     string `json:"title"`
}

type UpsertContestRequest struct {
	ContestId   int32   `json:"contest_id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	StartTime   int64   `json:"start_time"`
	EndTime     int64   `json:"end_time"`
	ProblemIds  []int32 `json:"problem_ids"`
}

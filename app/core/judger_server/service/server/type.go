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

type UpsertContestResponse struct {
	Code    int32               `json:"code"`
	Message string              `json:"message"`
	Contest ContestResponseData `json:"contest"`
}

type GetContestRequest struct {
	ContestId int32  `json:"contest_id"`
	Keyword   string `json:"keyword"`
}

type GetContestResponse struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
	Data    Contest `json:"data"`
}

type GetListContestsRequest struct {
	Page     int32  `json:"page"`
	PageSize int32  `json:"page_size"`
	Keyword  string `json:"keyword"`
}

type GetListContestsResponse struct {
	Code    int32                  `json:"code"`
	Message string                 `json:"message"`
	Data    GetListContestsResponseData `json:"data"`
}

type UpsertSubmissionResponse struct {
	Code    int32                      `json:"code"`
	Message string                     `json:"message"`
	Data    UpsertSubmissionResponseData `json:"data"`
}

type GetSubmissionRequest struct {
	SubmissionId int32 `json:"submission_id"`
}

type GetSubmissionResponse struct {
	Code    int32       `json:"code"`
	Message string      `json:"message"`
	Data    Submission  `json:"data"`
}

type GetListSubmissionsRequest struct {
	Page     int32  `json:"page"`
	PageSize int32  `json:"page_size"`
	Keyword  string `json:"keyword"`
}

type GetListSubmissionsResponse struct {
	Code    int32                      `json:"code"`
	Message string                     `json:"message"`
	Data    GetListSubmissionsResponseData `json:"data"`
}

type GetResultRequest struct {
	ResultId int32 `json:"result_id"`
}

type GetResultResponse struct {
	Code    int32   `json:"code"`
	Message string  `json:"message"`
	Data    Result  `json:"data"`
}

type GetListResultsRequest struct {
	Page     int32  `json:"page"`
	PageSize int32  `json:"page_size"`
	UserId   int32  `json:"user_id"`
}

type GetListResultsResponse struct {
	Code    int32                  `json:"code"`
	Message string                 `json:"message"`
	Data    GetListResultsResponseData `json:"data"`
}

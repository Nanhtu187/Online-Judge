package problem_repo

type GetListProblemRequest struct {
	ContestId int32  `json:"contest_id"`
	CreatedBy int32  `json:"created_by"`
	Keywords  string `json:"keywords"`
	Limit     int    `json:"limit"`
	Offset    int    `json:"offset"`
}

type ProblemPreview struct {
	ProblemId int32  `json:"problem_id"`
	Title     string `json:"title"`
}

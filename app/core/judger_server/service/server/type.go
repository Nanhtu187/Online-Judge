package server

type UpsertProblemRequest struct {
	id          int32  `json:"problem_id"`
	title       string `json:"title"`
	description string `json:"description"`
}

package server

import "github.com/Nanhtu187/Online-Judge/proto/rpc/judger_server"

func validateUpsertProblemRequest(request *judger_server.UpsertProblemRequest) (UpsertProblemRequest, error) {
	return UpsertProblemRequest{
		ProblemId:   request.ProblemId,
		Title:       request.Name,
		Description: request.Description,
	}, nil
}

func validateUpsertSubmissionRequest(req *judger_server.UpsertSubmissionRequest) (UpsertSubmissionRequest, error) {
	if req.ProblemId == 0 {
		return UpsertSubmissionRequest{}, ErrMissingProblemId
	}
	if req.UserId == 0 {
		return UpsertSubmissionRequest{}, ErrMissingUserId
	}
	return UpsertSubmissionRequest{
		ProblemId: req.ProblemId,
		ContestId: req.ContestId,
		UserId:    req.UserId,
		Language:  req.Language,
		Code:      req.Code,
	}, nil
}

func validateGetProblemRequest(req *judger_server.GetProblemRequest) (GetProblemRequest, error) {
	if req.ProblemId == 0 {
		return GetProblemRequest{}, ErrMissingProblemId
	}
	return GetProblemRequest{
		ProblemId: req.ProblemId,
	}, nil
}

func validateGetProblemsRequest(req *judger_server.GetListProblemsRequest) (GetListProblemRequest, error) {
	request := GetListProblemRequest{}
	if req.Keyword != "" {
		request.Keywords = req.Keyword
	}
	if req.Page != 0 {
		request.Offset = req.Page
	}
	if req.PageSize != 0 {
		request.Limit = req.PageSize
	}
	return request, nil
}

func validateUpsertContestRequest(req *judger_server.UpsertContestRequest) (UpsertContestRequest, error) {
	if req.ContestId == 0 {
		return UpsertContestRequest{}, ErrMissingContestId
	}
	if req.Name == "" {
		return UpsertContestRequest{}, ErrMissingContestName
	}
	return UpsertContestRequest{
		ContestId: req.ContestId,
		Name:      req.Name,
		Description: req.Description,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
		
	}, nil
}

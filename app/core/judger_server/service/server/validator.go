package server

import "github.com/Nanhtu187/Online-Judge/proto/rpc/judger_server"

func validateUpsertProblemRequest(request *judger_server.UpsertProblemRequest) (UpsertProblemRequest, error) {
	return UpsertProblemRequest{
		id:          request.ProblemId,
		title:       request.Name,
		description: request.Description,
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
		problemId: req.ProblemId,
		contestId: req.ContestId,
		userId:    req.UserId,
		language:  req.Language,
		code:      req.Code,
	}, nil
}

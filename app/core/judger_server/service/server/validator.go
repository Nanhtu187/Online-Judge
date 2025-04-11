package server

import "github.com/Nanhtu187/Online-Judge/proto/rpc/judger_server"

func validateUpsertProblemRequest(request *judger_server.UpsertProblemRequest) (UpsertProblemRequest, error) {
	return UpsertProblemRequest{
		id:          request.ProblemId,
		title:       request.Name,
		description: request.Description,
	}, nil
}

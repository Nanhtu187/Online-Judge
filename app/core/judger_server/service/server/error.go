package server

import (
	common "github.com/Nanhtu187/Online-Judge/app/common/share"
	"github.com/Nanhtu187/Online-Judge/app/core/judger_server/pkg/errors"
)

// Error missing problem id
var ErrMissingProblemId = errors.WithMessage(common.ErrParamIsRequired, "problem id is required")

// Error missing user id
var ErrMissingUserId = errors.WithMessage(common.ErrParamIsRequired, "user id is required")

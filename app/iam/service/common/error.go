package common

import "github.com/Nanhtu187/Online-Judge/app/iam/pkg/errors"

// Error missing param
var ErrParamIsRequired = errors.New(402, 400, "param is required")

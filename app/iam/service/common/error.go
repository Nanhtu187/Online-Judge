package common

import "github.com/Nanhtu187/Online-Judge/app/iam/pkg/errors"

// Error missing param
var ErrParamIsRequired = errors.New(402, 400, "param is required")

// Error record not found
var ErrRecordNotFound = errors.New(404, 404, "record not found")

// Error unauthorized
var ErrUnauthorized = errors.New(401, 401, "Unauthorized")

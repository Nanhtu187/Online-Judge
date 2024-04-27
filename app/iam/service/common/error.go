package common

import "github.com/Nanhtu187/Online-Judge/app/iam/pkg/errors"

// Error missing param
var ErrParamIsRequired = errors.New(402, 400, "param is required")

// Error record not found
var ErrRecordNotFound = errors.New(404, 404, "record not found")

// Error unauthorized
var ErrUnauthorized = errors.New(401, 401, "Unauthorized")

// Error forbidden
var ErrForbidden = errors.New(403, 403, "Forbidden")

// Error objec conflict
var ErrConflict = errors.New(409, 409, "Conflict")

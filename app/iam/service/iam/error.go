package iam

import (
	"github.com/Nanhtu187/Online-Judge/app/iam/pkg/errors"
	"github.com/Nanhtu187/Online-Judge/app/iam/service/common"
)

// Error missing name when create user
var ErrNameIsRequiredWhenCreateUser = errors.WithMessage(common.ErrParamIsRequired, "Name is required when create user")

// Error missing password when create user
var ErrPasswordIsRequiredWhenCreateUser = errors.WithMessage(common.ErrParamIsRequired, "Password is required when create user")

// ErrUserIdOrUsernameIsRequired ...
var ErrUserIdOrUsernameIsRequired = errors.WithMessage(common.ErrParamIsRequired, "UserId or Username is required")

// ErrUserNotFound ...
var ErrUserNotFound = errors.WithMessage(common.ErrRecordNotFound, "User not found")

// ErrPasswordNotMatch ...
var ErrPasswordNotMatch = errors.WithMessage(common.ErrUnauthorized, "Password not match")

// ErrInvalidToken ...
var ErrInvalidToken = errors.WithMessage(common.ErrUnauthorized, "Invalid token")

// ErrForbiddenToDeleteUser ...
var ErrForbiddenToDeleteUser = errors.WithMessage(common.ErrForbidden, "Forbidden to delete user")

// ErrUserExisted ...
var ErrUserExisted = errors.WithMessage(common.ErrConflict, "User existed")

package iam

import (
	"github.com/Nanhtu187/Online-Judge/app/iam/pkg/errors"
	"github.com/Nanhtu187/Online-Judge/app/iam/service/common"
)

// Error missing name when create user
var ErrNameIsRequiredWhenCreateUser = errors.WithMessage(common.ErrParamIsRequired, "Name is required when create user")

// Error missing password when create user
var ErrPasswordIsRequiredWhenCreateUser = errors.WithMessage(common.ErrParamIsRequired, "Password is required when create user")

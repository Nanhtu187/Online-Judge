package iam

import "github.com/Nanhtu187/Online-Judge/proto/rpc/iam/v1"

func validateCreateUserRequest(request UpsertUserRequest) error {
	if request.Name == "" {
		return ErrNameIsRequiredWhenCreateUser
	}
	if request.Password == "" {
		return ErrPasswordIsRequiredWhenCreateUser
	}
	return nil
}

func vaildateGetUserRequest(request *iam.GetUserRequest) (GetUserRequest, error) {
	if request.UserId == 0 && request.Username == "" {
		return GetUserRequest{}, ErrUserIdOrUsernameIsRequired
	}
	req := GetUserRequest{}
	if request.UserId != 0 {
		req.UserId = int(request.UserId)
	}
	if request.Username != "" {
		req.Username = request.Username
	}
	return req, nil
}

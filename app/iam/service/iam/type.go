package iam

type UpsertUserRequest struct {
	Username string `validate:"required"`
	Password string `validate:"max=255"`
	Name     string `validate:"max=255"`
	School   string `validate:"max=255"`
	Class    string `validate:"max=255"`
}

type UpsertUserResponse struct {
	UserId int
}

type GetUserRequest struct {
	UserId   int
	Username string
}

type GetUserResponse struct {
	UserId int
	Name   string
	School string
	Class  string
}

type DeleteUserRequest struct {
}

type DeleteUserResponse struct {
}

type LoginRequest struct {
}

type LoginResponse struct {
}

type RefreshTokenRequest struct {
}

type RefreshTokenResponse struct {
}

type GetListUserRequest struct {
}

type GetListUserResponse struct {
}

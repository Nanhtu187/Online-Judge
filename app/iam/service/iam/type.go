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
}

type GetUserResponse struct {
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
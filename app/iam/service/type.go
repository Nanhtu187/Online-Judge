package service

type UpsertUserRequest struct {
	Username string `validate:"required"`
	Password string `validate:"required"`
	Name     string `validate:"required"`
	School   string `validate:"max=255"`
	Class    string `validate:"max=255"`
}

type UpsertUserResponse struct {
	UserId int `json:"user_id"`
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

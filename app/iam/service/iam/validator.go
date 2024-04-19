package iam

func validateCreateUserRequest(request UpsertUserRequest) error {
	if request.Name == "" {
		return ErrNameIsRequiredWhenCreateUser
	}
	if request.Password == "" {
		return ErrPasswordIsRequiredWhenCreateUser
	}
	return nil
}

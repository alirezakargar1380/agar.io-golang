package users_types

type SignInRequest struct {
	Username string `validate:"required"`
	Password string `validate:"required"`
}

type SignUpRequerst struct {
	Username string
	Password string
}

package auth

type DummyLoginRequest struct {
	Role string `json:"role" validate:"required,oneof=moderator employee"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Role     string `json:"role" validate:"required,oneof=moderator employee"`
}

type RegisterResponse struct {
	Id    string `json:"id" validate:"required"`
	Email string `json:"email" validate:"required"`
	Role  string `json:"role" validate:"required"`
}

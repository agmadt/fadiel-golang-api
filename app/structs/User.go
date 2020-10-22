package structs

type User struct {
	ID       string `db:"id"`
	Email    string `db:"email"`
	Name     string `db:"name"`
	Password string `db:"password"`
}

type LoginRequest struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required"`
}

type LoginResponse struct {
	Token string            `json:"access_token"`
	Data  LoginResponseData `json:"data"`
}

type LoginResponseData struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

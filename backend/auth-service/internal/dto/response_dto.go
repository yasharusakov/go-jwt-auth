package dto

type UserResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

type AuthResponse struct {
	AccessToken string       `json:"access_token"`
	User        UserResponse `json:"user"`
}

type ErrorResponse struct {
	Message string `json:"error"`
	Code    int    `json:"code"`
}

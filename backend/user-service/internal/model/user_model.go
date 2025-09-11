package model

type User struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserWithoutPassword struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

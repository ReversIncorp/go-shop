package models

type User struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
	// другие поля, например, имя, дата регистрации и т.д.
}

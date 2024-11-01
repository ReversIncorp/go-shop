package entities

type User struct {
	ID       uint64 `json:"id"`
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	IsSeller bool   `json:"is_seller" validate:"required"`
}

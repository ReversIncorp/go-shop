package entities

type User struct {
	ID       uint64   `json:"id"`
	Email    string   `json:"email" validate:"required,email"`
	Password string   `json:"password" validate:"required,min=8"`
	IsOwner  bool     `json:"is_owner"`
	Stores   []uint64 `json:"stores"`
}

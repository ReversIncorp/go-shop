package entities

type User struct {
	ID       string   `json:"id"`
	Email    string   `json:"email"`
	Password string   `json:"password"`
	IsOwner  bool     `json:"is_owner"`
	Stores   []uint64 `json:"stores"`
}

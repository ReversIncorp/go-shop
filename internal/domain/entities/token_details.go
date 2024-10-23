package entities

// TokenDetails хранит данные о токенах
type TokenDetails struct {
	Token     string `json:"token"`
	UUID      string `json:"uuid"`
	AtExpires int64  `json:"expires"`
}

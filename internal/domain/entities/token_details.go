package entities

// TokenDetails хранит данные о токенах
type TokenDetails struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	AccessUUID   string `json:"access_uuid"`
	RefreshUUID  string `json:"refresh_uuid"`
	AtExpires    int64  `json:"at_expires"`
	RtExpires    int64  `json:"rt_expires"`
}

func (t *TokenDetails) ToTokens() *Tokens {
	return &Tokens{RefreshToken: t.RefreshToken, AccessToken: t.AccessToken}
}

package entities

type Tokens struct {
	RefreshToken *TokenDetails `json:"refresh_token"`
	AccessToken  *TokenDetails `json:"access_token"`
}

func (t *Tokens) CleanOutput() map[string]string {
	return map[string]string{
		"access_token":  t.AccessToken.Token,
		"refresh_token": t.RefreshToken.Token,
	}
}

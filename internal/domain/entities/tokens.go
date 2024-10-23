package entities

type Tokens struct {
	RefreshToken string `json:"refresh_token"`
	AccessToken  string `json:"access_token"`
}

func (t *Tokens) CleanOutput() map[string]string {
	return map[string]string{
		"access_token":  t.AccessToken,
		"refresh_token": t.RefreshToken,
	}
}

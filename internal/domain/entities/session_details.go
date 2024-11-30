package entities

type SessionDetails struct {
	DeviceInfo   string
	IPAddress    string
	AccessToken  string
	RefreshToken string
	ExpiresAt    int64
}

func (t *SessionDetails) CleanOutput() map[string]string {
	return map[string]string{
		"access_token":  t.AccessToken,
		"refresh_token": t.RefreshToken,
	}
}

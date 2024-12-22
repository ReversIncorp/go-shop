package entities

type SessionDetails struct {
	DeviceInfo   string `json:"device_info,omitempty"`
	IPAddress    string `json:"ip_address,omitempty"`
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	ExpiresAt    int64  `json:"expires_at,omitempty"`
}

func (t *SessionDetails) CleanOutput() map[string]string {
	return map[string]string{
		"access_token":  t.AccessToken,
		"refresh_token": t.RefreshToken,
	}
}

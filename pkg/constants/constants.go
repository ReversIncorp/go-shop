package constants

import "time"

const (
	RefreshTokenLifetime = time.Hour * 24 * 365
	AccessTokenLifetime  = 72 * time.Hour
)

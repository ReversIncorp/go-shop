package config

import (
	"os"
	"sync"
)

type Config struct {
	JWTKey []byte
}

type configHolder struct {
	once   sync.Once
	config *Config
}

// GetConfig предоставляет доступ к конфигурации.
func GetConfig() *Config {
	var holder configHolder

	holder.once.Do(func() {
		holder.config = &Config{
			JWTKey: []byte(os.Getenv("JWT_SECRET_KEY")),
		}
	})

	return holder.config
}

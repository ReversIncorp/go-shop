package repository

import (
	"marketplace/internal/domain/entities"
)

type JWTRepository interface {
	SaveSession(userID uint64, sessionID string, session *entities.SessionDetails) error
	GetSession(userID uint64, sessionID string) (*entities.SessionDetails, error)
	GetAllSessions(userID uint64) (map[string]*entities.SessionDetails, error)
	DeleteSession(userID uint64, sessionID string) error
}

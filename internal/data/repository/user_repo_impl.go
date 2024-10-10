package repository

import (
	"errors"
	"marketplace/internal/domain/entities"
	repository2 "marketplace/internal/domain/repository"
)

type userRepositoryImpl struct {
	// Можно использовать базу данных здесь, например, Gorm или другое хранилище
	users map[string]entities.User
}

func NewUserRepository() repository2.UserRepository {
	return &userRepositoryImpl{
		users: make(map[string]entities.User),
	}
}

func (r *userRepositoryImpl) Create(user entities.User) error {
	if _, exists := r.users[user.Email]; exists {
		return errors.New("user already exists")
	}
	r.users[user.Email] = user
	return nil
}

func (r *userRepositoryImpl) FindByEmail(email string) (entities.User, error) {
	user, exists := r.users[email]
	if !exists {
		return entities.User{}, errors.New("user not found")
	}
	return user, nil
}

func (r *userRepositoryImpl) FindByID(id uint64) (entities.User, error) {
	for _, user := range r.users {
		if user.ID == id {
			return user, nil
		}
	}
	return entities.User{}, errors.New("user not found")
}

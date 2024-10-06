package repository

import (
	"errors"
	"marketplace/internal/domain/models"
	repository2 "marketplace/internal/domain/repository"
)

type userRepositoryImpl struct {
	// Можно использовать базу данных здесь, например, Gorm или другое хранилище
	users map[string]models.User
}

func NewUserRepository() repository2.UserRepository {
	return &userRepositoryImpl{
		users: make(map[string]models.User),
	}
}

func (r *userRepositoryImpl) Create(user models.User) error {
	if _, exists := r.users[user.Email]; exists {
		return errors.New("user already exists")
	}
	r.users[user.Email] = user
	return nil
}

func (r *userRepositoryImpl) FindByEmail(email string) (models.User, error) {
	user, exists := r.users[email]
	if !exists {
		return models.User{}, errors.New("user not found")
	}
	return user, nil
}

func (r *userRepositoryImpl) FindByID(id string) (models.User, error) {
	for _, user := range r.users {
		if user.ID == id {
			return user, nil
		}
	}
	return models.User{}, errors.New("user not found")
}

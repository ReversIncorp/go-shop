package repository

import (
	"errors"
	"gorm.io/gorm"
	"marketplace/internal/domain/entities"
	repository2 "marketplace/internal/domain/repository"
)

type userRepositoryImpl struct {
	// Можно использовать базу данных здесь, например, Gorm или другое хранилище
	//users map[string]entities.User
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) repository2.UserRepository {
	return &userRepositoryImpl{
		db: db,
	}
}

func (r *userRepositoryImpl) Create(user entities.User) error {
	var existingUser entities.User
	if err := r.db.Where("email = ?", user.Email).First(&existingUser).Error; err == nil {
		return errors.New("user already exists")
	}

	return r.db.Create(&user).Error
}

func (r *userRepositoryImpl) FindByEmail(email string) (entities.User, error) {
	var user entities.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities.User{}, errors.New("user not found")
		}
		return entities.User{}, err
	}
	return user, nil
}

func (r *userRepositoryImpl) FindByID(id uint64) (entities.User, error) {
	var user entities.User
	if err := r.db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities.User{}, errors.New("user not found")
		}
		return entities.User{}, err
	}
	return user, nil
}

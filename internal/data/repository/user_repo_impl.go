package repository

import (
	"errors"
	"marketplace/internal/domain/entities"
	"marketplace/internal/domain/repository"
	errorHandling "marketplace/pkg/error_handling"

	"github.com/ztrue/tracerr"
	"gorm.io/gorm"
)

type userRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &userRepositoryImpl{
		db: db,
	}
}

func (r *userRepositoryImpl) Create(user *entities.User) (uint64, error) {
	// Проверка на уникальность email
	var existingUser entities.User
	err := r.db.Where("email = ?", user.Email).First(&existingUser).Error
	if err == nil {
		return 0, tracerr.Wrap(errors.New("user already exists"))
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, tracerr.Wrap(err)
	}

	if err := r.db.Create(user).Error; err != nil {
		return 0, tracerr.Wrap(err)
	}
	return user.ID, nil
}

func (r *userRepositoryImpl) FindByEmail(email string) (entities.User, error) {
	var user entities.User
	err := r.db.
		Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities.User{}, errorHandling.ErrUserNotFound
		}
		return entities.User{}, tracerr.Wrap(err)
	}
	return user, nil
}

func (r *userRepositoryImpl) FindByID(id uint64) (entities.User, error) {
	var user entities.User
	err := r.db.
		First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities.User{}, errorHandling.ErrUserNotFound
		}
		return entities.User{}, tracerr.Wrap(err)
	}
	return user, nil
}

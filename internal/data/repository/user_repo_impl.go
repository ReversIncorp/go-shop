package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"marketplace/internal/domain/entities"
	"marketplace/internal/domain/repository"
	errorResponses "marketplace/pkg/error_handling"

	"github.com/ztrue/tracerr"
)

type userRepositoryImpl struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) repository.UserRepository {
	return &userRepositoryImpl{
		db: db,
	}
}

func (r *userRepositoryImpl) Create(user entities.User) (uint64, error) {
	var existingUserID uint64
	query := `SELECT id FROM users WHERE email = $1`
	err := r.db.QueryRow(query, user.Email).Scan(&existingUserID)
	if err == nil {
		return 0, tracerr.Wrap(errors.New("user already exists"))
	} else if !errors.Is(err, sql.ErrNoRows) {
		return 0, tracerr.Wrap(fmt.Errorf("failed to check existing user: %w", err))
	}

	insertQuery := `INSERT INTO users 
    (name, 
     email, 
     password, 
     is_seller) 
	VALUES ($1, $2, $3, $4) RETURNING id`
	var newUserID uint64
	err = r.db.QueryRow(insertQuery,
		user.Name,
		user.Email,
		user.Password,
		user.IsSeller).Scan(&newUserID)

	if err != nil {
		return 0,  tracerr.Wrap(fmt.Errorf("failed to create user: %w", err))
	}

	return newUserID, nil
}

func (r *userRepositoryImpl) FindByEmail(email string) (entities.User, error) {
	var user entities.User
	query := `SELECT 
    	id,
       name, 
       email,
       password,
       is_seller
	          FROM users WHERE email = $1`
	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.IsSeller)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entities.User{}, errorResponses.ErrUserNotFound
		}
		return entities.User{}, tracerr.Wrap(err)
	}
	return user, nil
}

func (r *userRepositoryImpl) FindByID(id uint64) (entities.User, error) {
	var user entities.User
	query := `SELECT 
    	id, 
       name, 
       email, 
       password, 
       is_seller 
	          FROM users WHERE id = $1`
	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.IsSeller)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entities.User{}, errorResponses.ErrUserNotFound
		}
		return entities.User{}, err
	}
	return user, nil
}

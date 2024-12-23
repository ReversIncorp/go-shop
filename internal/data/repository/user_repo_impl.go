package repository

import (
	"database/sql"
	"errors"
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

func (r *userRepositoryImpl) Create(user entities.User) error {
	var existingUser entities.User
	query := `SELECT id, email FROM users WHERE email = $1`
	err := r.db.QueryRow(query, user.Email).Scan(&existingUser.ID, &existingUser.Email)
	if err == nil {
		return tracerr.Wrap(errors.New("user already exists"))
	}

	insertQuery := `INSERT INTO users 
    (name, 
     email, 
     password, 
     is_seller) 
	                VALUES ($1, $2, $3, $4) RETURNING id`
	err = r.db.QueryRow(insertQuery,
		user.Name,
		user.Email,
		user.Password,
		user.IsSeller).Scan(&user.ID)

	if err != nil {
		return tracerr.Wrap(err)
	}

	return nil
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
		if err == sql.ErrNoRows {
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
		if err == sql.ErrNoRows {
			return entities.User{}, errorResponses.ErrUserNotFound
		}
		return entities.User{}, tracerr.Wrap(err)
	}
	return user, nil
}

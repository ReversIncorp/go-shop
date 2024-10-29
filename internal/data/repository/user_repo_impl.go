package repository

import (
	"database/sql"
	"errors"
	"marketplace/internal/domain/entities"
	"marketplace/internal/domain/repository"

	"github.com/lib/pq"
)

type userRepositoryImpl struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) repository.UserRepository {
	return &userRepositoryImpl{
		db: db,
	}
}

func (r *userRepositoryImpl) AddOwningStore(userID, storeID uint64) error {
	_, err := r.db.Exec(`
		UPDATE users 
		SET owning_stores = array_append(array_remove(owning_stores, $1), $1) 
		WHERE id = $2
	`, storeID, userID)
	return err
}

func (r *userRepositoryImpl) IsOwnsStore(userID, storeID uint64) (bool, error) {
	var owningStores []int64 // Изменение типа для лучшей совместимости
	err := r.db.QueryRow(
		"SELECT owning_stores FROM users WHERE id = $1",
		userID,
	).Scan(pq.Array(&owningStores))

	if err != nil {
		return false, err
	}

	for _, sID := range owningStores {
		if uint64(sID) == storeID { // Приведение типа при сравнении
			return true, nil
		}
	}

	return false, errors.New("user does not own this store")
}

func (r *userRepositoryImpl) Create(user entities.User) error {
	var existingUser entities.User
	query := `SELECT id, email FROM users WHERE email = $1`
	err := r.db.QueryRow(query, user.Email).Scan(&existingUser.ID, &existingUser.Email)
	if err == nil {
		return errors.New("user already exists")
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
		return err
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
			return entities.User{}, errors.New("user not found")
		}
		return entities.User{}, err
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
			return entities.User{}, errors.New("user not found")
		}
		return entities.User{}, err
	}
	return user, nil
}

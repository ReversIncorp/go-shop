package utils

import (
	"database/sql"
	"errors"
	"github.com/lib/pq"
)

// Проверка существует ли магазин
func CheckStoreExists(db *sql.DB, storeID int64) (bool, error) {
	var existingStoreID uint64
	err := db.QueryRow("SELECT id FROM stores WHERE id = $1", storeID).Scan(&existingStoreID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, errors.New("store not found")
		}
		return false, err
	}
	return true, nil
}

// Проверка является ли юзер владельцем магазина
func CheckUserOwnsStore(db *sql.DB, userID, storeID int64) (bool, error) {
	var owningStores []int64
	err := db.QueryRow("SELECT owning_stores FROM users WHERE id = $1", userID).Scan(pq.Array(&owningStores))
	if err != nil {
		return false, err
	}

	for _, sID := range owningStores {
		if sID == storeID {
			return true, nil
		}
	}

	return false, errors.New("user does not own this store")
}

package utils

import (
	"database/sql"
	"errors"
)

// CheckCategoryBelongsToStore проверяет, принадлежит ли категория конкретному магазину
func CheckCategoryBelongsToStore(db *sql.DB, categoryID, storeID int64) (bool, error) {
	var existingStoreID int64
	err := db.QueryRow("SELECT store_id FROM categories WHERE id = $1", categoryID).Scan(&existingStoreID)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, errors.New("category not found")
		}
		return false, err
	}

	if existingStoreID == storeID {
		return true, nil
	}

	return false, errors.New("category does not belong to this store")
}

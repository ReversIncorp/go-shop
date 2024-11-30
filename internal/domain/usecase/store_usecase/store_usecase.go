package storeUsecase

import (
	"errors"
	"marketplace/internal/domain/entities"
	"marketplace/internal/domain/repository"
)

// StoreUseCase реализует интерфейс StoreUseCase
type StoreUseCase struct {
	storeRepo    repository.StoreRepository
	categoryRepo repository.CategoryRepository
}

// NewStoreUseCase создает новый экземпляр StoreUseCase
func NewStoreUseCase(
	storeRepo repository.StoreRepository,
	categoryRepo repository.CategoryRepository,
) *StoreUseCase {
	return &StoreUseCase{
		storeRepo:    storeRepo,
		categoryRepo: categoryRepo,
	}
}

// CreateStore создает новый магазин
func (s *StoreUseCase) CreateStore(store entities.Store, userID uint64) error {
	return s.storeRepo.Save(store, userID)
}

// GetStoreByID получает магазин по ID
func (s *StoreUseCase) GetStoreByID(id uint64) (entities.Store, error) {
	return s.storeRepo.FindByID(id)
}

// UpdateStore обновляет существующий магазин
func (s *StoreUseCase) UpdateStore(store entities.Store) error {
	return s.storeRepo.Update(store)
}

// DeleteStore удаляет магазин по ID
func (s *StoreUseCase) DeleteStore(id uint64) error {
	return s.storeRepo.Delete(id)
}

// GetAllStores получает все магазины
func (s *StoreUseCase) GetAllStores() ([]entities.Store, error) {
	return s.storeRepo.FindAll()
}

// IsUserStoreAdmin проверка является ли пользователь админом стора
func (s *StoreUseCase) IsUserStoreAdmin(storeID uint64, uid uint64) (bool, error) {
	storeExists, err := s.storeRepo.IsExist(storeID)
	if err != nil || !storeExists {
		return false, errors.New("store not found")
	}

	return s.storeRepo.IsUserStoreAdmin(storeID, uid)
}

// AttachCategoryToStore добавляет категорию к магазину
func (s *StoreUseCase) AttachCategoryToStore(storeID, categoryID uint64) error {
	categoryExist, err := s.categoryRepo.IsExist(categoryID)
	if err != nil || !categoryExist {
		return errors.New("category not found")
	}

	isAttached, err := s.storeRepo.IsCategoryAttached(storeID, categoryID)
	if err != nil || isAttached {
		return errors.New("category is attached to store")
	}

	return s.storeRepo.AttachCategory(storeID, categoryID)
}

// DetachCategoryFromStore открепляет категорию от магазина
func (s *StoreUseCase) DetachCategoryFromStore(storeID, categoryID uint64) error {
	categoryExist, err := s.categoryRepo.IsExist(categoryID)
	if err != nil || !categoryExist {
		return errors.New("category not found")
	}

	isAttached, err := s.storeRepo.IsCategoryAttached(storeID, categoryID)
	if err != nil || !isAttached {
		return errors.New("category is not attached to store")
	}

	return s.storeRepo.DetachCategory(storeID, categoryID)
}

// GetStoresByFilters получает все магазины по фильтрам
func (s *StoreUseCase) GetStoresByFilters(filters entities.StoreSearchParams) ([]entities.Store, *uint64, error) {
	return s.storeRepo.FindStoresByParams(filters)
}

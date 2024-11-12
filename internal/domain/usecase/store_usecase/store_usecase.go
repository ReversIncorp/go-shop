package storeUsecase

import (
	"marketplace/internal/domain/entities"
	"marketplace/internal/domain/repository"

	errorResponses "marketplace/pkg/errors"
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
	storeExists, err := s.storeRepo.IsExist(id)
	if err != nil || !storeExists {
		return errorResponses.ErrStoreNotFound
	}

	err = s.storeRepo.Delete(id)
	if err != nil {
		return errorResponses.ErrInternalServerError
	}

	return nil
}

// GetAllStores получает все магазины
func (s *StoreUseCase) GetAllStores() ([]entities.Store, error) {
	return s.storeRepo.FindAll()
}

// IsUserStoreAdmin проверка является ли пользователь админом стора
func (s *StoreUseCase) IsUserStoreAdmin(storeID uint64, uid uint64) (bool, error) {
	storeExists, err := s.storeRepo.IsExist(storeID)
	if err != nil || !storeExists {
		return false, errorResponses.ErrStoreNotFound
	}

	admin, err := s.storeRepo.IsUserStoreAdmin(storeID, uid)
	if err != nil {
		return false, errorResponses.ErrInternalServerError
	}

	return admin, nil
}

// AttachCategoryToStore добавляет категорию к магазину
func (s *StoreUseCase) AttachCategoryToStore(storeID, categoryID uint64) error {
	categoryExist, err := s.categoryRepo.IsExist(categoryID)
	if err != nil || !categoryExist {
		return errorResponses.ErrCategoryNotFound
	}

	isAttached, err := s.storeRepo.IsCategoryAttached(storeID, categoryID)
	if err != nil || isAttached {
		return errorResponses.ErrCategoryAttached
	}

	err = s.storeRepo.AttachCategory(storeID, categoryID)
	if err != nil {
		return errorResponses.ErrInternalServerError
	}

	return nil
}

// DetachCategoryFromStore открепляет категорию от магазина
func (s *StoreUseCase) DetachCategoryFromStore(storeID, categoryID uint64) error {
	categoryExist, err := s.categoryRepo.IsExist(categoryID)
	if err != nil || !categoryExist {
		return errorResponses.ErrCategoryNotFound
	}

	isAttached, err := s.storeRepo.IsCategoryAttached(storeID, categoryID)
	if err != nil || !isAttached {
		return errorResponses.ErrCategoryNotAttached
	}

	err = s.storeRepo.DetachCategory(storeID, categoryID)
	if err != nil {
		return errorResponses.ErrInternalServerError
	}

	return nil
}

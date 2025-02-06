package storeusecase

import (
	"marketplace/internal/domain/entities"
	"marketplace/internal/domain/repository"

	errorHandling "marketplace/pkg/error_handling"

	"github.com/ztrue/tracerr"
)

// StoreUseCase реализует интерфейс StoreUseCase.
type StoreUseCase struct {
	storeRepo    repository.StoreRepository
	categoryRepo repository.CategoryRepository
}

// NewStoreUseCase создает новый экземпляр StoreUseCase.
func NewStoreUseCase(
	storeRepo repository.StoreRepository,
	categoryRepo repository.CategoryRepository,
) *StoreUseCase {
	return &StoreUseCase{
		storeRepo:    storeRepo,
		categoryRepo: categoryRepo,
	}
}

// CreateStore создает новый магазин.
func (s *StoreUseCase) CreateStore(store entities.Store, userID uint64) error {
	err := s.storeRepo.Save(store, userID)
	if err != nil {
		return tracerr.Wrap(err)
	}
	return nil
}

// GetStoreByID получает магазин по ID.
func (s *StoreUseCase) GetStoreByID(id uint64) (entities.Store, error) {
	store, err := s.storeRepo.FindByID(id)
	if err != nil {
		return store, tracerr.Wrap(err)
	}
	return store, nil
}

// UpdateStore обновляет существующий магазин.
func (s *StoreUseCase) UpdateStore(store entities.Store) error {
	err := s.storeRepo.Update(store)
	if err != nil {
		return tracerr.Wrap(err)
	}
	return nil
}

// DeleteStore удаляет магазин по ID.
func (s *StoreUseCase) DeleteStore(id uint64) error {
	storeExists, err := s.storeRepo.IsExist(id)

	if err != nil {
		return tracerr.Wrap(err)
	}

	if !storeExists {
		return errorHandling.ErrStoreNotFound
	}

	errDelete := s.storeRepo.Delete(id)
	if errDelete != nil {
		return tracerr.Wrap(errDelete)
	}
	return nil
}

// GetAllStores получает все магазины.
func (s *StoreUseCase) GetAllStores() ([]entities.Store, error) {
	stores, err := s.storeRepo.FindAll()
	if err != nil {
		return nil, tracerr.Wrap(err)
	}
	return stores, nil
}

// IsUserStoreAdmin проверка является ли пользователь админом стора.
func (s *StoreUseCase) IsUserStoreAdmin(storeID uint64, uid uint64) (bool, error) {
	storeExists, err := s.storeRepo.IsExist(storeID)

	if err != nil {
		return false, tracerr.Wrap(err)
	}

	if !storeExists {
		return false, errorHandling.ErrStoreNotFound
	}

	is, errTwo := s.storeRepo.IsUserStoreAdmin(storeID, uid)
	if errTwo != nil {
		return false, tracerr.Wrap(errTwo)
	}
	return is, nil
}

// AttachCategoryToStore добавляет категорию к магазину.
func (s *StoreUseCase) AttachCategoryToStore(storeID, categoryID uint64) error {
	categoryExist, err := s.categoryRepo.IsExist(categoryID)
	if err != nil || !categoryExist {
		return errorHandling.ErrCategoryNotFound
	}

	isAttached, err := s.storeRepo.IsCategoryAttached(storeID, categoryID)
	if err != nil || isAttached {
		return errorHandling.ErrCategoryAttached
	}

	errAttach := s.storeRepo.AttachCategory(storeID, categoryID)
	if errAttach != nil {
		return tracerr.Wrap(errAttach)
	}
	return nil
}

// DetachCategoryFromStore открепляет категорию от магазина.
func (s *StoreUseCase) DetachCategoryFromStore(storeID, categoryID uint64) error {
	categoryExist, err := s.categoryRepo.IsExist(categoryID)
	if err != nil || !categoryExist {
		return errorHandling.ErrCategoryNotFound
	}

	isAttached, err := s.storeRepo.IsCategoryAttached(storeID, categoryID)
	if err != nil || !isAttached {
		return errorHandling.ErrCategoryNotAttached
	}

	errDetach := s.storeRepo.DetachCategory(storeID, categoryID)
	if errDetach != nil {
		return tracerr.Wrap(errDetach)
	}
	return nil
}

// GetStoresByFilters получает все магазины по фильтрам.
func (s *StoreUseCase) GetStoresByFilters(filters entities.StoreSearchParams) ([]entities.Store, *uint64, error) {
	stores, id, err := s.storeRepo.FindStoresByParams(filters)
	if err != nil {
		return nil, nil, tracerr.Wrap(err)
	}

	return stores, id, nil

}

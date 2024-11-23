package entities

// StoreSearchParams структура для фильтрации магазинов
type StoreSearchParams struct {
	CategoryID *uint64 `json:"category_id"`
	Name       *string `json:"name"`
	Limit      *uint64 `json:"limit" validate:"required,gt=0"`
	Cursor     *uint64 `json:"cursor" validate:"required"`
}

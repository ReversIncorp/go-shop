package entities

// ProductSearchParams структура для фильтрации продуктов
type ProductSearchParams struct {
	StoreID    *uint64  `json:"store_id"`
	CategoryID *uint64  `json:"category_id"`
	MinPrice   *float64 `json:"min_price"`
	MaxPrice   *float64 `json:"max_price"`
	Name       *string  `json:"name"`
	Limit      *uint64  `json:"limit" validate:"required,gt=0"`
	Cursor     *uint64  `json:"cursor" validate:"required,gt=0"`
}

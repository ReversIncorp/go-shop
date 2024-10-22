package entities

// ProductSearchParams структура для фильтрации продуктов
type ProductSearchParams struct {
	StoreID    *int64   `json:"store_id"`
	CategoryID *int64   `json:"category_id"`
	MinPrice   *float64 `json:"min_price"`
	MaxPrice   *float64 `json:"max_price"`
	Name       *string  `json:"name"`
}

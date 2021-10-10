package entity

// Album represents an album record.
type Product struct {
	ID    int64   `json:"id"`
	Name  string  `json:"name"`
	Stock int32   `json:"stock"`
	Price float64 `json:"price"`
}

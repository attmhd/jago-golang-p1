package model

type Product struct {
	ID           int     `json:"id"`
	CategoryID   int     `json:"category_id"`
	CategoryName string  `json:"category_name"`
	Name         string  `json:"name"`
	Price        float64 `json:"price"`
	Stock        int     `json:"stock"`
}

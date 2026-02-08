package models

type Product struct {
	ID           int     `json:"id"`
	CategoryID   int     `json:"category_id"`
	CategoryName string  `json:"category_name"`
	Name         string  `json:"name"`
	Price        float64 `json:"price"`
	Stock        int     `json:"stock"`
}

// Model untuk menampilkan produk terlaris dengan jumlah terjual
type TopSellingProduct struct {
	Name    string `json:"name"`
	QtySold int    `json:"qty_sold"`
}

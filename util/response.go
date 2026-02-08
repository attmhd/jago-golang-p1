package util

type JSONResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type Category struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type ProductResp struct {
	ID       int      `json:"id"`
	Name     string   `json:"name"`
	Price    float64  `json:"price"`
	Stock    int      `json:"stock"`
	Category Category `json:"category"`
}

type SalesSummary struct {
	TotalRevenue   int            `json:"total_revenue"`
	TotalTransaksi int            `json:"total_transaksi"`
	ProdukTerlaris ProdukTerlaris `json:"produk_terlaris"`
}

type ProdukTerlaris struct {
	Nama       string `json:"nama"`
	QtyTerjual int    `json:"qty_terjual"`
}

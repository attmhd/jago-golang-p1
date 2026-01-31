package handler

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

package handler

type JSONResponse struct {
	Message string `json:"message"`
	Data    any    `json:"data"`
}

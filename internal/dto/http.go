package dto

type SuccessResponse[T any] struct {
	Message string `json:"message"`
	Data    T      `json:"data"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

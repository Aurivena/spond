package models

type ErrorResponse struct {
	Status string `json:"status"`
	Error  ErrorDetail
}

type ErrorDetail struct {
	Title   any `json:"title"`
	Message any `json:"message"`
}

package response

type SendErrorOutput struct {
	Status string      `json:"status"`
	Error  ErrorDetail `json:"error"`
}

type SendSuccessOutput struct {
	Status string `json:"status"`
	Output any    `json:"error"`
}

type ErrorResponse struct {
	Status StatusCode  `json:"status"`
	Error  ErrorDetail `json:"error"`
}

type ErrorDetail struct {
	Title   any    `json:"title"`
	Message string `json:"message"`
}

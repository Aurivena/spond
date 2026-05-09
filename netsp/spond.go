// Package netsp lives to facilitate communication
// between server and web via JSON structures.
package netsp

import (
	"github.com/Aurivena/spond/v3/netstatus"
)

type ErrorDetail struct {
	Title    string `json:"title"`
	Message  string `json:"message"`
	Solution string `json:"solution"`
}

type Response[T any] struct {
	Code netstatus.Code
	Data T
}

func ValidateBuildError(code int, title, message, solution string) bool {
	if err := validate(title, message); err != nil {
		return false
	}
	return true
}

// BuildError forms an error structure for responding to the client.
// If the input parameters do not pass validation, it returns an error with the UnprocessableEntity code.
func BuildError[T any](code netstatus.Code, data T) *Response[T] {
	return &Response[T]{
		Code: code,
		Data: data,
	}
}

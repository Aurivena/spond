// Package netsp lives to facilitate communication
// between server and web via JSON structures.
package netsp

import (
	"fmt"
	"log"
	"net/http"
)

type ErrorDetail struct {
	Title    string
	Message  string
	Solution string
}

type AppError struct {
	Code   int
	Detail ErrorDetail
}

type writeError struct {
	Code  string   `json:"code"`
	Error errorDTO `json:"error"`
}

type errorDTO struct {
	Title    string `json:"title"`
	Message  string `json:"message"`
	Solution string `json:"solution"`
}

// SendResponseSuccess sends a successful JSON envelope.
// status is the envelope status, is the payload for the client.
// Generics Type usages for typing data
func SendResponseSuccess[T any](w http.ResponseWriter, code int, data T) {
	if !isValid(code) {
		log.Printf("[ERROR] SendResponseSuccess: status code %d don`t exists", code)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if code == NoContent {
		w.WriteHeader(int(code))
		return
	}

	write(w, data, code)
}

// SendResponseError sends the error to the client as JSON.
// err — structure with error details.
func SendResponseError(w http.ResponseWriter, err *AppError) {
	if err == nil {
		return
	}
	if !isValid(int(err.Code)) {
		log.Printf("[ERROR] SendResponseSuccess: status code %d don`t exists", err.Code)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	output := &writeError{
		Error: errorDTO{
			Message:  err.Detail.Message,
			Title:    err.Detail.Title,
			Solution: err.Detail.Solution,
		},
	}

	write(w, output, err.Code)
}

// AppendCode adds a new status code and message to the statusMessages card.
// If the code already exists, returns the error.
func AppendCode(code int, message string) error {
	if code < 100 || code > 599 {
		return fmt.Errorf("spond: invalid HTTP status code %d", code)
	}

	if err := appendCode(code, message); err != nil {
		return fmt.Errorf("spond: failed to append code %d: %w", code, err)
	}

	return nil
}

// BuildError forms an error structure for responding to the client.
// If the input parameters do not pass validation, it returns an error with the UnprocessableEntity code.
func BuildError(code int, title, message, solution string) *AppError {
	if err := validate(title, message); err != nil {
		return &AppError{
			Code: UnprocessableEntity,
			Detail: ErrorDetail{
				Title:    Invalid,
				Message:  err.Error(),
				Solution: SolutionError,
			},
		}
	}
	return &AppError{
		Code: code,
		Detail: ErrorDetail{
			Title:    title,
			Message:  message,
			Solution: solution,
		},
	}
}

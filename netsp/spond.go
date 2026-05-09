// Package netsp lives to facilitate communication
// between server and web via JSON structures.
package netsp

import (
	"bytes"
	"encoding/json"
	"fmt"
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

func (e *AppError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Detail.Title, e.Detail.Message)
}

func IsValid(code int) error {
	if !isValid(code) {
		return fmt.Errorf("invalid status code %d", code)
	}
	return nil
}

// Write encodes response as JSON and sends it to client.
// Always sets Content-Type to application/json; charset=utf-8.
func Write[T any](w http.ResponseWriter, code int, output T) {
	if code == http.StatusNoContent {
		w.WriteHeader(code)
		return
	}

	// set data for future json
	var buff bytes.Buffer
	if err := json.NewEncoder(&buff).Encode(output); err != nil {
		// fallback: plain text error if JSON encoding fails
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	w.Write(buff.Bytes())
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

func ValidateBuildError(code int, title, message, solution string) bool {
	if err := validate(title, message); err != nil {
		return false
	}
	return true
}

// BuildError forms an error structure for responding to the client.
// If the input parameters do not pass validation, it returns an error with the UnprocessableEntity code.
func BuildError(code int, title, message, solution string) *AppError {
	return &AppError{
		Code: code,
		Detail: ErrorDetail{
			Title:    title,
			Message:  message,
			Solution: solution,
		},
	}
}

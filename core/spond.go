// Package core lives to facilitate communication
// between server and web via JSON structures.
package core

import (
	"fmt"
	"net/http"

	"sync"

	"github.com/Aurivena/spond/v3/envelope"
)

type Spond struct {
	statusMessages map[int]string //storage status code and provides code append.
	mu             *sync.RWMutex
}

var defaultSpond = &Spond{
	statusMessages: make(map[int]string),
	mu:             &sync.RWMutex{},
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
func SendResponseSuccess(w http.ResponseWriter, code int, data any) {
	if !envelope.IsValid(code) {
		// It error developer
		panic(fmt.Errorf("Status code %d don`t exists", code))
	}

	if code == envelope.NoContent {
		w.WriteHeader(int(code))
		return
	}

	write(w, &data, code)
}

// SendResponseError sends the error to the client as JSON.
// err — structure with error details.
func SendResponseError(w http.ResponseWriter, err *envelope.AppError) {
	if err == nil {
		return
	}
	if !envelope.IsValid(int(err.Code)) {
		// It error developer
		panic(fmt.Errorf("status code %d don`t exists", err.Code))
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
	defaultSpond.mu.Lock()
	defer defaultSpond.mu.Unlock()

	return envelope.AppendCode(code, message)
}

// BuildError forms an error structure for responding to the client.
// If the input parameters do not pass validation, it returns an error with the UnprocessableEntity code.
func BuildError(code int, title, message, solution string) *envelope.AppError {
	if err := validate(title, message); err != nil {
		return &envelope.AppError{
			Code: envelope.UnprocessableEntity,
			Detail: envelope.ErrorDetail{
				Title:    envelope.Invalid,
				Message:  err.Error(),
				Solution: envelope.SolutionError,
			},
		}
	}
	return &envelope.AppError{
		Code: code,
		Detail: envelope.ErrorDetail{
			Title:    title,
			Message:  message,
			Solution: solution,
		},
	}
}

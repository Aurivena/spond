// Package core lives to facilitate communication
// between server and web via JSON structures.
package core

import (
	"fmt"
	"net/http"

	"sync"

	"github.com/Aurivena/spond/v2/envelope"
)

type Spond struct {
	statusMessages map[envelope.StatusCode]string //storage status code and provides code append.
	mu             sync.RWMutex
}

type writeSuccess struct {
	Data any
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

// For initialization  struct Spond
// Usage example
// spond:=NewSpond()
// spond.SendResponseSuccess(w, spond.Created, nil)
func NewSpond() *Spond {
	return &Spond{statusMessages: envelope.StatusMessages}
}

// SendResponseSuccess sends a successful JSON envelope.
// status is the envelope status, is the payload for the client.
func (s *Spond) SendResponseSuccess(w http.ResponseWriter, code envelope.StatusCode, data any) {
	if !s.codeExists(code) {
		// It error developer
		panic(fmt.Errorf("Status code %d don`t exists", code))
	}

	if code == envelope.NoContent {
		w.WriteHeader(int(code))
		return
	}

	output := writeSuccess{
		Data: &data,
	}

	write(w, output, code)
}

// SendResponseError sends the error to the client as JSON.
// err — structure with error details.
func (s *Spond) SendResponseError(w http.ResponseWriter, err *envelope.AppError) {
	if err == nil {
		return
	}
	if !s.codeExists(err.Code) {
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
func (s *Spond) AppendCode(code envelope.StatusCode, message string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exist := s.statusMessages[code]; exist {
		return envelope.ErrorAppendCode
	}
	s.statusMessages[code] = message
	return nil
}

// BuildError forms an error structure for responding to the client.
// If the input parameters do not pass validation, it returns an error with the UnprocessableEntity code.
func (s *Spond) BuildError(code envelope.StatusCode, title, message, solution string) *envelope.AppError {
	if err := validate(title, message); err != nil {
		return &envelope.AppError{
			Code: envelope.UnprocessableEntity,
			Detail: envelope.ErrorDetail{
				Title:    envelope.Invalid.Error(),
				Message:  err.Error(),
				Solution: envelope.SolutionError.Error(),
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

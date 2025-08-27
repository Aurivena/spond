// Package spond lives to facilitate communication
// between server and web via JSON structures.
package spond

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"sync"

	"github.com/Aurivena/spond/envelope"
)

type Spond struct {
	statusMessages map[envelope.StatusCode]string //storage status code and provides code append.
	mu             sync.RWMutex                   // for code append
}

type writeSuccess struct {
	Data any `json:"data,omitempty"`
}

type writeError struct {
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
// spond.SendenvelopeSuccess(c, spond.OK, nil)
func NewSpond() *Spond {
	return &Spond{statusMessages: envelope.StatusMessages}
}

// SendResponseSuccess sends a successful JSON envelope via gin.Context.
// status is the envelope status, Successenvelope is the payload for the client.
func (s *Spond) SendResponseSuccess(w http.ResponseWriter, code envelope.StatusCode, data any) {
	if !s.codeExists(code) {
		// It warn developer
		panic(fmt.Errorf("Status code %d don`t exists", code))
	}

	if code == envelope.NoContent {
		w.WriteHeader(int(code))
		return
	}

	output := writeSuccess{
		Data: &data,
	}

	var buff bytes.Buffer
	if err := json.NewEncoder(&buff).Encode(output); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(int(code))
	w.Write(buff.Bytes())
}

// SendResponseError sends the error to the client as JSON via gin.Context.
// rsp — structure with error details.
func (s *Spond) SendResponseError(w http.ResponseWriter, err envelope.AppError) {
	if !s.codeExists(err.Code) {
		// It warn developer
		panic(fmt.Errorf("Status code %d don`t exists", err.Code))
	}

	output := &writeError{
		Error: errorDTO{
			Message:  err.Detail.Message,
			Title:    err.Detail.Title,
			Solution: err.Detail.Solution,
		},
	}

	var buff bytes.Buffer
	if err := json.NewEncoder(&buff).Encode(output); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(int(err.Code))
	w.Write(buff.Bytes())
}

// AppendCode adds a new status code and message to the statusMessages card.
// If the code already exists, returns the errorAppendCode error.
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
				Title:    "invalid",
				Message:  err.Error(),
				Solution: "Recheck limits for title and message pls :)",
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

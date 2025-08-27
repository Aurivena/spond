// Package spond lives to facilitate communication
// between server and web via JSON structures.
package spond

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"sync"

	"github.com/Aurivena/spond/envelope"
)

const (
	titleInvalid     = "invalid value for title"
	messageInvalid   = "invalid value for message"
	invalid          = "invalid"
	unknownStatus    = "unknown status"
	maxTitleLength   = 256
	maxMessageLength = 1024
)

var errorAppendCode = errors.New("this code already exists")

type Spond struct {
	statusMessages map[envelope.StatusCode]string //storage status code and provides code append.
	mu             sync.RWMutex                   // for code append
}

type writeSuccess struct {
	Data  *any        `json:"data,omitempty"`
	Error *writeError `json:"error,omitempty"`
}

type writeError struct {
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
	if !s.codeExist(code) {
		// It warn developer
		panic(fmt.Errorf("status code %d don`t exists", code))
	}

	output := writeSuccess{
		Data: &data,
	}

	w.WriteHeader(int(code))
	_ = json.NewEncoder(w).Encode(output)
}

// SendResponseError sends the error to the client as JSON via gin.Context.
// rsp — structure with error details.
func (s *Spond) SendResponseError(w http.ResponseWriter, err envelope.Error) {
	if !s.codeExist(err.Code) {
		// It warn developer
		panic(fmt.Errorf("status code %d don`t exists", err.Code))
	}

	w.WriteHeader(int(err.Code))

	output := writeSuccess{
		Error: &writeError{
			Message:  err.Detail.Message,
			Title:    err.Detail.Title,
			Solution: err.Detail.Solution,
		},
	}

	_ = json.NewEncoder(w).Encode(output)
}

// AppendCode adds a new status code and message to the statusMessages card.
// If the code already exists, returns the errorAppendCode error.
func (s *Spond) AppendCode(code envelope.StatusCode, message string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exist := s.statusMessages[code]; exist {
		return errorAppendCode
	}
	s.statusMessages[code] = message
	return nil
}

// BuildError forms an error structure for responding to the client.
// If the input parameters do not pass validation, it returns an error with the UnprocessableEntity code.
func (s *Spond) BuildError(code envelope.StatusCode, title, message, solution string) *envelope.Error {
	if err := validate(title, message); err != nil {
		return &envelope.Error{
			Code: envelope.UnprocessableEntity,
			Detail: envelope.ErrorDetail{
				Title:    "invalid",
				Message:  err.Error(),
				Solution: "Recheck limits for title and message pls :)",
			},
		}
	}
	return &envelope.Error{
		Code: code,
		Detail: envelope.ErrorDetail{
			Title:    title,
			Message:  message,
			Solution: solution,
		},
	}
}

func (s *Spond) codeExist(code envelope.StatusCode) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, exist := s.statusMessages[code]
	return exist
}

// validate checks the length of the title and message.
// Returns the error text when restrictions are violated.
func validate(title, message string) error {
	if len(title) == 0 || len(title) > maxTitleLength {
		return fmt.Errorf("%s", titleInvalid)
	}
	if len(message) == 0 || len(message) > maxMessageLength {
		return fmt.Errorf("%s", messageInvalid)
	}
	return nil
}

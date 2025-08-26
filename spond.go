// Package spond lives to facilitate communication
// between server and web via JSON structures.
package spond

import (
	"errors"
	"fmt"
	"log/slog"

	"sync"

	"github.com/Aurivena/spond/envelope"
	"github.com/gin-gonic/gin"
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

type sendResponse[T any] struct {
	Data   *T           `json:"data,omitempty"`
	Status string       `json:"status"`
	Error  *errorDetail `json:"error,omitempty"`
}

type errorDetail struct {
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
func (s *Spond) SendResponseSuccess(c *gin.Context, code envelope.StatusCode, data any) {
	if c == nil {
		slog.Error("SendenvelopeSuccess: gin.Context == nil")
		return
	}

	if _, exist := s.statusMessages[code]; exist {
		panic(fmt.Errorf("status code %d don`t exists", code))
	}

	output := sendResponse[any]{
		Data:   &data,
		Status: code.String(),
	}

	c.JSON(int(code), output)
}

// SendResponseError sends the error to the client as JSON via gin.Context.
// rsp — structure with error details.
func (s *Spond) SendResponseError(c *gin.Context, code envelope.StatusCode, err sendResponse[any]) {
	if c == nil {
		slog.Error("gin.Context == nil")
		return
	}

	if _, exist := s.statusMessages[code]; exist {
		panic(fmt.Errorf("status code %d don`t exists", code))
	}

	//TODO добавить mapToHTTP
	output := sendResponse[any]{
		Data:   nil,
		Status: code.String(),
		Error:  err.Error,
	}

	c.AbortWithStatusJSON(int(code), output)
}

// AppendCode adds a new status code and message to the statusMessages card.
// If the code already exists, returns the errorAppendCode error.
func (s *Spond) AppendCode(code envelope.StatusCode, message string) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if _, exist := s.statusMessages[code]; exist {
		return errorAppendCode
	}
	s.statusMessages[code] = message
	return nil
}

// BuildError forms an error structure for responding to the client.
// If the input parameters do not pass validation, it returns an error with the UnprocessableEntity code.
func (s *Spond) BuildError(code envelope.StatusCode, title, message, solution string) errorDetail {
	if err := validate(title, message); err != nil {
		return errorDetail{
			Title:    invalid,
			Message:  err.Error(),
			Solution: "",
		}
	}

	return errorDetail{
		Title:    title,
		Message:  message,
		Solution: solution,
	}
}

// validate checks the length of the title and message.
// Returns the error text when restrictions are violated.
func validate(title, message string) error {
	if len(title) == 0 || len(title) > maxTitleLength {
		return fmt.Errorf("$w", titleInvalid)
	}
	if len(message) == 0 || len(message) > maxMessageLength {
		return fmt.Errorf("$w", messageInvalid)
	}
	return nil
}

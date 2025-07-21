// Package spond lives to facilitate communication
// between server and web via JSON structures.
package spond

import (
	"errors"
	"log/slog"
	"spond/response"
	"sync"

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
	statusMessages map[response.StatusCode]string //storage status code and provides code append.
	mu             sync.RWMutex                   // for code append
}

// For initialization  struct Spond
func NewSpond() *Spond {
	return &Spond{statusMessages: response.StatusMessages}
}

// SendResponseSuccess sends a successful JSON response via gin.Context.
// status is the response status, SuccessResponse is the payload for the client.
func (s *Spond) SendResponseSuccess(c *gin.Context, status response.StatusCode, successResponse any) {
	if c == nil {
		slog.Error("SendResponseSuccess: gin.Context == nil")
		return
	}
	output := response.SendSuccessOutput{
		Status: status.String(),
		Output: successResponse,
	}

	c.JSON(int(status.MapToHTTPStatus()), output)
}

// SendResponseError sends the error to the client as JSON via gin.Context.
// rsp — structure with error details.
func (s *Spond) SendResponseError(c *gin.Context, rsp response.ErrorResponse) {
	if c == nil {
		slog.Error("gin.Context == nil")
		return
	}

	statusMessage, ok := s.statusMessages[rsp.Status]
	if !ok {
		slog.Warn("SendResponseError: неизвестный статус код", "статус", rsp.Status)
		statusMessage = unknownStatus
	}

	output := response.SendErrorOutput{
		Status: statusMessage,
		Error:  rsp.Error,
	}

	c.AbortWithStatusJSON(int(rsp.Status.MapToHTTPStatus()), output)
}

// AppendCode adds a new status code and message to the statusMessages card.
// If the code already exists, returns the errorAppendCode error.
func (s *Spond) AppendCode(code response.StatusCode, message string) error {
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
func (s *Spond) BuildError(code response.StatusCode, title, message string) response.ErrorResponse {
	if err := validate(title, message); err != "" {
		return response.ErrorResponse{
			Status: response.UnprocessableEntity,
			Error:  response.ErrorDetail{Title: invalid, Message: err},
		}
	}

	return response.ErrorResponse{
		Status: code.MapToHTTPStatus(),
		Error:  response.ErrorDetail{Title: title, Message: message},
	}
}

// validate checks the length of the title and message.
// Returns the error text when restrictions are violated.
func validate(title, message string) string {
	if len(title) == 0 || len(title) > maxTitleLength {
		return titleInvalid
	}
	if len(message) == 0 || len(message) > maxMessageLength {
		return messageInvalid
	}
	return ""
}

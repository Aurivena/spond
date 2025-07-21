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
	titleInvalid     = "title invalid"
	messageInvalid   = "message invalid"
	invalid          = "Invalid"
	unknownStatus    = "unknown status"
	maxTitleLength   = 256
	maxMessageLength = 1024
)

var errorAppendCode = errors.New("this code already exists")

type Spond struct {
	statusMessages map[response.StatusCode]string
	mu             sync.RWMutex
}

func NewSpond() *Spond {
	return &Spond{statusMessages: response.StatusMessages}
}

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

func (s *Spond) AppendCode(code response.StatusCode, message string) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if _, exist := s.statusMessages[code]; exist {
		return errorAppendCode
	}
	s.statusMessages[code] = message
	return nil
}

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

func validate(title, message string) string {
	if len(title) == 0 || len(title) > maxTitleLength {
		return titleInvalid
	}
	if len(message) == 0 || len(message) > maxMessageLength {
		return messageInvalid
	}
	return ""
}

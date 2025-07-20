package spond

import (
	"encoding/json"
	"errors"
	"log/slog"
	"spond/response"
	"sync"

	"github.com/gin-gonic/gin"
)

var (
	errorAppendCode = errors.New("this code already exists")
	titleInvalid    = "title invalid"
	messageInvalid  = "message invalid"
	invalid         = "Invalid"
	unknownStatus   = "unknown status"
)

type Spond struct {
	statusMessages map[response.StatusCode]string
	mu             sync.Mutex
}

func NewSpond() (*Spond, error) {
	return &Spond{
		statusMessages: response.StatusMessages,
	}, nil
}

func (s *Spond) SendResponseSuccess(c *gin.Context, status response.StatusCode, successResponse any) {
	s.mu.Lock()
	defer s.mu.Unlock()

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
	s.mu.Lock()
	defer s.mu.Unlock()

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
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exist := s.statusMessages[code]; exist {
		return errorAppendCode
	}
	s.statusMessages[code] = message
	return nil
}

func (s *Spond) BuildError(code response.StatusCode, title, message string) response.ErrorResponse {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := validate(title, message); err != "" {
		return response.ErrorResponse{
			Status: response.BadRequest,
			Error:  response.ErrorDetail{Title: invalid, Message: err},
		}
	}

	return response.ErrorResponse{
		Status: code.MapToHTTPStatus(),
		Error:  response.ErrorDetail{Title: title, Message: message},
	}
}

func validate(title, message any) string {
	if _, err := json.Marshal(title); err != nil {
		return titleInvalid
	}
	if _, err := json.Marshal(message); err != nil {
		return messageInvalid
	}
	return ""
}

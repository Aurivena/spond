package spond

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log/slog"
	"net/http"
	"os"
	"spond/faults"
	"sync"
)

type Impl struct {
	StatusMessages map[StatusCode]string
	Mu             *sync.RWMutex
	Out            io.Writer
}

func NewImpl() *Impl {
	return &Impl{StatusMessages: statusMessages, Mu: &sync.RWMutex{}, Out: os.Stdout}
}

func (e *Impl) SendResponseSuccess(c *gin.Context, status StatusCode, successResponse any) {
	if c == nil {
		slog.Error("SendResponseSuccess: gin.Context == nil")
		return
	}
	output := SendSuccessOutput{
		Status: status.String(),
		Output: successResponse,
	}

	c.JSON(http.StatusOK, output)
}

func (e *Impl) SendResponseError(c *gin.Context, rsp ErrorResponse) {
	if c == nil {
		slog.Error("gin.Context == nil")
		return
	}

	statusMessage, ok := e.StatusMessages[rsp.Status]
	if !ok {
		slog.Warn("SendResponseError: неизвестный статус код", "статус", rsp.Status)
		statusMessage = faults.UnknownStatus
	}

	output := SendErrorOutput{
		Status: statusMessage,
		Error:  rsp.Error,
	}

	c.AbortWithStatusJSON(http.StatusOK, output)
}

func (e *Impl) SayHello() {
	fmt.Fprint(e.Out, "Hello it Spond!\n")
}

func (e *Impl) AppendCode(code StatusCode, message string) error {
	e.Mu.Lock()
	defer e.Mu.Unlock()

	if _, exist := e.StatusMessages[code]; exist {
		return faults.ErrorAppendCode
	}
	e.StatusMessages[code] = message
	return nil
}

func (e *Impl) BuildError(code StatusCode, title, message any) ErrorResponse {
	if err := validate(title, message); err != "" {
		return ErrorResponse{
			Status: BadRequest,
			Error:  ErrorDetail{Title: faults.Invalid, Message: err},
		}
	}

	return ErrorResponse{
		Status: code.mapToHTTPStatus(),
		Error:  ErrorDetail{Title: title, Message: message},
	}
}

func validate(title, message any) string {
	if _, err := json.Marshal(title); err != nil {
		return faults.TitleInvalid
	}
	if _, err := json.Marshal(message); err != nil {
		return faults.MessageInvalid
	}
	return ""
}

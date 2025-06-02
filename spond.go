package spond

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"os"
	"spond/models"
	"sync"
)

type Impl struct {
	StatusMessages map[StatusCode]string
	Mu             *sync.RWMutex
	Out            io.Writer
}

var (
	ErrorAppendCode = errors.New("этот код уже существует")
)

func NewImpl() *Impl {
	return &Impl{StatusMessages: statusMessages, Mu: &sync.RWMutex{}, Out: os.Stdout}
}

func (e *Impl) SayHello() {
	fmt.Fprint(e.Out, "Hello it Spond!\n")
}

func (e *Impl) AppendCode(code StatusCode, message string) error {
	e.Mu.Lock()
	defer e.Mu.Unlock()

	if _, exist := e.StatusMessages[code]; exist {
		return ErrorAppendCode
	}
	e.StatusMessages[code] = message
	return nil
}

func (e *Impl) BuildError(c *gin.Context, code StatusCode, title, message any) models.ErrorResponse {
	if c == nil {
		return models.ErrorResponse{
			Status: ContextIsNil.String(),
			Error:  models.ErrorDetail{Title: "", Message: ""},
		}
	}

	if err := validate(title, message); err != nil {
		return models.ErrorResponse{
			Status: InternalServerError.String(),
			Error:  models.ErrorDetail{Title: "Invalid", Message: err.Error()},
		}
	}

	return models.ErrorResponse{
		Status: code.String(),
		Error:  models.ErrorDetail{Title: title, Message: message},
	}
}

func validate(title, message any) error {
	if _, err := json.Marshal(title); err != nil {
		return errors.New("title invalid")
	}
	if _, err := json.Marshal(message); err != nil {
		return errors.New("message invalid")
	}
	return nil
}

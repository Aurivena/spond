package spond

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Aurivena/spond/response"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAppendCode(t *testing.T) {
	s := NewSpond()
	tests := []struct {
		name    string
		code    response.StatusCode
		message string
		wantErr error
	}{
		{"добавление нового кода", 9999, "new code", nil},
		{"повторный код", 9999, "again", errorAppendCode},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := s.AppendCode(tt.code, tt.message)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestBuildError(t *testing.T) {
	s := NewSpond()
	tests := []struct {
		name    string
		code    response.StatusCode
		title   string
		message string
		want    response.ErrorResponse
	}{
		{"валидный ответ", response.UnprocessableEntity, "", "Описание",
			response.ErrorResponse{Status: response.UnprocessableEntity, Error: response.ErrorDetail{Title: invalid, Message: titleInvalid}}},
		{"пустые значения", response.UnprocessableEntity, "title", "",
			response.ErrorResponse{Status: response.UnprocessableEntity, Error: response.ErrorDetail{Title: invalid, Message: messageInvalid}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := s.BuildError(tt.code, tt.title, tt.message)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSendResponseSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	s := NewSpond()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	s.SendResponseSuccess(c, response.Success, map[string]string{"foo": "bar"})
	assert.Equal(t, http.StatusOK, w.Code)
	var output response.SendSuccessOutput
	err := json.Unmarshal(w.Body.Bytes(), &output)
	assert.NoError(t, err)
	assert.Equal(t, response.Success.String(), output.Status)
	assert.Equal(t, map[string]interface{}{"foo": "bar"}, output.Output)
}

func TestSendResponseError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	s := NewSpond()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	errTitle := "Доступ запрещен"
	errMessage := "У вас недостаточно прав"
	errResp := response.ErrorResponse{
		Status: response.BadRequest,
		Error:  response.ErrorDetail{Title: errTitle, Message: errMessage},
	}
	s.SendResponseError(c, errResp)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	var output response.SendErrorOutput
	err := json.Unmarshal(w.Body.Bytes(), &output)
	assert.NoError(t, err)
	assert.Equal(t, s.statusMessages[response.BadRequest], output.Status)
	assert.Equal(t, errTitle, output.Error.Title)
	assert.Equal(t, errMessage, output.Error.Message)
}

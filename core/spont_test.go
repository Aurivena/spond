package spond

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Aurivena/spond/envelope"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestAppendCode(t *testing.T) {
	s := NewSpond()
	tests := []struct {
		name    string
		code    envelope.StatusCode
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
		code    envelope.StatusCode
		title   string
		message string
		want    envelope.Errorenvelope
	}{
		{"валидный ответ", envelope.UnprocessableEntity, "", "Описание",
			envelope.Errorenvelope{Status: envelope.UnprocessableEntity, Error: envelope.ErrorDetail{Title: invalid, Message: titleInvalid}}},
		{"пустые значения", envelope.UnprocessableEntity, "title", "",
			envelope.Errorenvelope{Status: envelope.UnprocessableEntity, Error: envelope.ErrorDetail{Title: invalid, Message: messageInvalid}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := s.BuildError(tt.code, tt.title, tt.message)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSendenvelopeSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)
	s := NewSpond()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	s.SendenvelopeSuccess(c, envelope.Success, map[string]string{"foo": "bar"})
	assert.Equal(t, http.StatusOK, w.Code)
	var output envelope.SendSuccessOutput
	err := json.Unmarshal(w.Body.Bytes(), &output)
	assert.NoError(t, err)
	assert.Equal(t, envelope.Success.String(), output.Status)
	assert.Equal(t, map[string]interface{}{"foo": "bar"}, output.Output)
}

func TestSendenvelopeError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	s := NewSpond()
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	errTitle := "Доступ запрещен"
	errMessage := "У вас недостаточно прав"
	errResp := envelope.Errorenvelope{
		Status: envelope.BadRequest,
		Error:  envelope.ErrorDetail{Title: errTitle, Message: errMessage},
	}
	s.SendenvelopeError(c, errResp)
	assert.Equal(t, http.StatusBadRequest, w.Code)
	var output envelope.SendErrorOutput
	err := json.Unmarshal(w.Body.Bytes(), &output)
	assert.NoError(t, err)
	assert.Equal(t, s.statusMessages[envelope.BadRequest], output.Status)
	assert.Equal(t, errTitle, output.Error.Title)
	assert.Equal(t, errMessage, output.Error.Message)
}

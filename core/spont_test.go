package core_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Aurivena/spond/v2/core"
	"github.com/Aurivena/spond/v2/envelope"
	"github.com/stretchr/testify/assert"
)

type writeSuccess struct {
	Data any `json:"data,omitempty"`
}
type errorDTO struct {
	Title    string `json:"title"`
	Message  string `json:"message"`
	Solution string `json:"solution,omitempty"`
}
type writeError struct {
	Error errorDTO `json:"error"`
}

func TestAppendCode(t *testing.T) {
	s := core.NewSpond()

	err := s.AppendCode(envelope.StatusCode(9999), "new code")
	assert.NoError(t, err)

	err = s.AppendCode(envelope.StatusCode(9999), "again")
	assert.Error(t, err)
}

func TestBuildError(t *testing.T) {
	s := core.NewSpond()

	tests := []struct {
		name     string
		code     envelope.StatusCode
		title    string
		message  string
		solution string
		want     envelope.AppError
	}{
		{
			name:     "пустой title → UnprocessableEntity + invalid",
			code:     envelope.UnprocessableEntity,
			title:    "",
			message:  "Описание",
			solution: "",
			want: envelope.AppError{
				Code: envelope.UnprocessableEntity,
				Detail: envelope.ErrorDetail{
					Title:    "invalid",
					Message:  "invalid value for title",
					Solution: "recheck limits for title and message pls :)",
				},
			},
		},
		{
			name:     "пустой message → UnprocessableEntity + invalid",
			code:     envelope.UnprocessableEntity,
			title:    "title",
			message:  "",
			solution: "",
			want: envelope.AppError{
				Code: envelope.UnprocessableEntity,
				Detail: envelope.ErrorDetail{
					Title:    "invalid",
					Message:  "invalid value for message",
					Solution: "recheck limits for title and message pls :)",
				},
			},
		},
		{
			name:     "валидный ввод → указанный код и детали",
			code:     envelope.BadRequest,
			title:    "Bad input",
			message:  "Некорректные данные",
			solution: "Проверьте поля",
			want: envelope.AppError{
				Code: envelope.BadRequest,
				Detail: envelope.ErrorDetail{
					Title:    "Bad input",
					Message:  "Некорректные данные",
					Solution: "Проверьте поля",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPtr := s.BuildError(tt.code, tt.title, tt.message, tt.solution)
			if gotPtr == nil {
				t.Fatalf("BuildError returned nil")
			}
			assert.Equal(t, tt.want, *gotPtr)
		})
	}
}

func TestSendResponseSuccess(t *testing.T) {
	s := core.NewSpond()
	w := httptest.NewRecorder()

	payload := map[string]string{"foo": "bar"}
	s.SendResponseSuccess(w, envelope.Success, payload)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))

	var out writeSuccess
	err := json.Unmarshal(w.Body.Bytes(), &out)
	assert.NoError(t, err)

	assert.Equal(t, map[string]any{"foo": "bar"}, out.Data)
}

func TestSendResponseSuccess_NoContent(t *testing.T) {
	s := core.NewSpond()
	w := httptest.NewRecorder()

	s.SendResponseSuccess(w, envelope.NoContent, nil)

	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Equal(t, "", w.Header().Get("Content-Type"))
	assert.Equal(t, 0, w.Body.Len())
}

func TestSendResponseError(t *testing.T) {
	s := core.NewSpond()
	w := httptest.NewRecorder()

	errTitle := "Доступ запрещен"
	errMessage := "У вас недостаточно прав"
	appErr := envelope.AppError{
		Code:   envelope.BadRequest,
		Detail: envelope.ErrorDetail{Title: errTitle, Message: errMessage, Solution: ""},
	}

	s.SendResponseError(w, &appErr)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))

	var out writeError
	err := json.Unmarshal(w.Body.Bytes(), &out)
	assert.NoError(t, err)

	assert.Equal(t, errTitle, out.Error.Title)
	assert.Equal(t, errMessage, out.Error.Message)
	assert.Equal(t, "", out.Error.Solution)
}

package netsp_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Aurivena/spond/v3/netsp"
	"github.com/stretchr/testify/assert"
)

func TestAppendCode(t *testing.T) {
	err := netsp.AppendCode(9999, "again")
	assert.Error(t, err)

	err = netsp.AppendCode(204, "no content")
	assert.Error(t, err)
}

func TestValidateBuildError(t *testing.T) {
	tests := []struct {
		name    string
		title   string
		message string
		want    bool
	}{
		{
			name:    "валидный ввод",
			title:   "Ошибка доступа",
			message: "У вас нет прав",
			want:    true,
		},
		{
			name:    "пустой заголовок",
			title:   "",
			message: "Описание",
			want:    false,
		},
		{
			name:    "пустое сообщение",
			title:   "Заголовок",
			message: "",
			want:    false,
		},
		{
			name:    "слишком длинный заголовок (> 256)",
			title:   string(make([]byte, 257)),
			message: "Описание",
			want:    false,
		},
		{
			name:    "слишком длинное сообщение (> 1024)",
			title:   "Заголовок",
			message: string(make([]byte, 1025)),
			want:    false,
		},
		{
			name:    "оба поля пусты",
			title:   "",
			message: "",
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := netsp.ValidateBuildError(http.StatusBadRequest, tt.title, tt.message, "")
			assert.Equal(t, tt.want, got, "Валидация должна возвращать %v для случая: %s", tt.want, tt.name)
		})
	}
}

func TestBuildError(t *testing.T) {
	tests := []struct {
		name     string
		code     int
		title    string
		message  string
		solution string
		want     netsp.AppError
	}{
		{
			name:     "стандартный случай",
			code:     http.StatusBadRequest,
			title:    "Ошибка",
			message:  "Что-то пошло не так",
			solution: "Попробуйте снова",
			want: netsp.AppError{
				Code: http.StatusBadRequest,
				Detail: netsp.ErrorDetail{
					Title:    "Ошибка",
					Message:  "Что-то пошло не так",
					Solution: "Попробуйте снова",
				},
			},
		},
		{
			name:     "пустые значения (свобода ответа)",
			code:     http.StatusInternalServerError,
			title:    "",
			message:  "",
			solution: "",
			want: netsp.AppError{
				Code: http.StatusInternalServerError,
				Detail: netsp.ErrorDetail{
					Title:    "",
					Message:  "",
					Solution: "",
				},
			},
		},
		{
			name:     "очень длинные строки (свобода ответа)",
			code:     http.StatusOK,
			title:    string(make([]byte, 1000)),
			message:  string(make([]byte, 5000)),
			solution: "Много текста",
			want: netsp.AppError{
				Code: http.StatusOK,
				Detail: netsp.ErrorDetail{
					Title:    string(make([]byte, 1000)),
					Message:  string(make([]byte, 5000)),
					Solution: "Много текста",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPtr := netsp.BuildError(tt.code, tt.title, tt.message, tt.solution)

			if gotPtr == nil {
				t.Fatalf("BuildError вернул nil, хотя должен был вернуть объект AppError")
			}

			assert.Equal(t, tt.want, *gotPtr, "Ошибка сборки объекта в случае: %s", tt.name)
		})
	}
}

func TestSendResponseSuccess(t *testing.T) {
	w := httptest.NewRecorder()

	payload := map[string]string{"foo": "bar"}
	netsp.Write(w, http.StatusOK, payload)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
	assert.JSONEq(t, `{"foo":"bar"}`, w.Body.String())
}

func TestSendResponseSuccess_NoContent(t *testing.T) {
	w := httptest.NewRecorder()

	netsp.Write[any](w, http.StatusNoContent, nil)

	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Equal(t, "", w.Header().Get("Content-Type"))
	assert.Equal(t, 0, w.Body.Len())
}

func TestSendResponseError(t *testing.T) {
	w := httptest.NewRecorder()

	errTitle := "Доступ запрещен"
	errMessage := "У вас недостаточно прав"
	appErr := netsp.AppError{
		Code:   http.StatusBadRequest,
		Detail: netsp.ErrorDetail{Title: errTitle, Message: errMessage, Solution: ""},
	}

	netsp.Write(w, appErr.Code, appErr.Detail)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))

	var out netsp.ErrorDetail
	err := json.Unmarshal(w.Body.Bytes(), &out)
	assert.NoError(t, err)

	assert.Equal(t, errTitle, out.Title)
	assert.Equal(t, errMessage, out.Message)
	assert.Equal(t, "", out.Solution)
}

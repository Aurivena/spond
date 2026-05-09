package netsp_test

import (
	"net/http"
	"testing"

	"github.com/Aurivena/spond/v3/netsp"
	"github.com/Aurivena/spond/v3/netstatus"
	"github.com/stretchr/testify/assert"
)

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
		code     netstatus.Code
		title    string
		message  string
		solution string
		want     netsp.Response[netsp.ErrorDetail]
	}{
		{
			name:     "стандартный случай",
			code:     netstatus.CodeBadRequest,
			title:    "Ошибка",
			message:  "Что-то пошло не так",
			solution: "Попробуйте снова",
			want: netsp.Response[netsp.ErrorDetail]{
				Code: netstatus.CodeBadRequest,
				Data: netsp.ErrorDetail{
					Title:    "Ошибка",
					Message:  "Что-то пошло не так",
					Solution: "Попробуйте снова",
				},
			},
		},
		{
			name:     "пустые значения (свобода ответа)",
			code:     netstatus.CodeInternalError,
			title:    "",
			message:  "",
			solution: "",
			want: netsp.Response[netsp.ErrorDetail]{
				Code: netstatus.CodeInternalError,
				Data: netsp.ErrorDetail{
					Title:    "",
					Message:  "",
					Solution: "",
				},
			},
		},
		{
			name:     "очень длинные строки (свобода ответа)",
			code:     netstatus.CodeSuccess,
			title:    string(make([]byte, 1000)),
			message:  string(make([]byte, 5000)),
			solution: "Много текста",
			want: netsp.Response[netsp.ErrorDetail]{
				Code: netstatus.CodeSuccess,
				Data: netsp.ErrorDetail{
					Title:    string(make([]byte, 1000)),
					Message:  string(make([]byte, 5000)),
					Solution: "Много текста",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := netsp.ErrorDetail{
				Title:    tt.title,
				Message:  tt.message,
				Solution: tt.solution,
			}
			gotPtr := netsp.BuildError(tt.code, data)

			if gotPtr == nil {
				t.Fatalf("BuildError вернул nil, хотя должен был вернуть объект AppError")
			}

			assert.Equal(t, tt.want, *gotPtr, "Ошибка сборки объекта в случае: %s", tt.name)
		})
	}
}

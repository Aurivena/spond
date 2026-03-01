package netsp_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Aurivena/spond/v3/netsp"
	"github.com/stretchr/testify/assert"
)

type errorDTO struct {
	Title    string `json:"title"`
	Message  string `json:"message"`
	Solution string `json:"solution,omitempty"`
}
type writeError struct {
	Error errorDTO `json:"error"`
}

func TestAppendCode(t *testing.T) {
	err := netsp.AppendCode(9999, "again")
	assert.Error(t, err)

	err = netsp.AppendCode(204, "no content")
	assert.Error(t, err)
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
			name:     "пустой title → UnprocessableEntity + invalid",
			code:     netsp.UnprocessableEntity,
			title:    "",
			message:  "Описание",
			solution: "",
			want: netsp.AppError{
				Code: netsp.UnprocessableEntity,
				Detail: netsp.ErrorDetail{
					Title:    netsp.Invalid,
					Message:  netsp.TitleInvalid,
					Solution: "Recheck limits for title and message pls :)",
				},
			},
		},
		{
			name:     "пустой message → UnprocessableEntity + invalid",
			code:     netsp.UnprocessableEntity,
			title:    "title",
			message:  "",
			solution: "",
			want: netsp.AppError{
				Code: netsp.UnprocessableEntity,
				Detail: netsp.ErrorDetail{
					Title:    netsp.Invalid,
					Message:  netsp.MessageInvalid,
					Solution: "Recheck limits for title and message pls :)",
				},
			},
		},
		{
			name:     "валидный ввод → указанный код и детали",
			code:     netsp.BadRequest,
			title:    "Bad input",
			message:  "Некорректные данные",
			solution: "Проверьте поля",
			want: netsp.AppError{
				Code: netsp.BadRequest,
				Detail: netsp.ErrorDetail{
					Title:    "Bad input",
					Message:  "Некорректные данные",
					Solution: "Проверьте поля",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPtr := netsp.BuildError(tt.code, tt.title, tt.message, tt.solution)
			if gotPtr == nil {
				t.Fatalf("BuildError returned nil")
			}
			assert.Equal(t, tt.want, *gotPtr)
		})
	}
}

func TestSendResponseSuccess(t *testing.T) {
	w := httptest.NewRecorder()

	payload := map[string]string{"foo": "bar"}
	netsp.SendResponseSuccess(w, netsp.Success, payload)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))
	assert.JSONEq(t, `{"foo":"bar"}`, w.Body.String())
}

func TestSendResponseSuccess_NoContent(t *testing.T) {
	w := httptest.NewRecorder()

	netsp.SendResponseSuccess[any](w, netsp.NoContent, nil)

	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Equal(t, "", w.Header().Get("Content-Type"))
	assert.Equal(t, 0, w.Body.Len())
}

func TestSendResponseError(t *testing.T) {
	w := httptest.NewRecorder()

	errTitle := "Доступ запрещен"
	errMessage := "У вас недостаточно прав"
	appErr := netsp.AppError{
		Code:   netsp.BadRequest,
		Detail: netsp.ErrorDetail{Title: errTitle, Message: errMessage, Solution: ""},
	}

	netsp.SendResponseError(w, &appErr)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, "application/json; charset=utf-8", w.Header().Get("Content-Type"))

	var out writeError
	err := json.Unmarshal(w.Body.Bytes(), &out)
	assert.NoError(t, err)

	assert.Equal(t, errTitle, out.Error.Title)
	assert.Equal(t, errMessage, out.Error.Message)
	assert.Equal(t, "", out.Error.Solution)
}

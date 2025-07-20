package spond_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"spond"
	"spond/faults"
	"spond/intertnal/response"
	"sync"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type testBuildErrorResponse struct {
	name     string
	title    any
	message  any
	code     response.StatusCode
	expected response.ErrorResponse
}

type testAppendCode struct {
	name    string
	code    response.StatusCode
	message string
	wantErr error
}

type testSendSuccessAndError struct {
	name     string
	c        *gin.Context
	status   response.StatusCode
	output   any
	expected response.SendSuccessOutput
}

type testSendError struct {
	Status string               `json:"status"`
	Error  response.ErrorDetail `json:"error"`
}

var invalidMessage = func() {}
var testStruct = struct {
	Name string
	Age  int
}{
	Name: "Aurivena",
	Age:  666,
}

var testsAppendCode = []testAppendCode{
	{
		name:    "TestAppendCode_RealImplementation: UnknownCode666",
		code:    666,
		message: "TestCode",
		wantErr: nil,
	},
	{
		name:    "TestAppendCode_RealImplementation: SuccessCodeExists",
		code:    response.Success,
		message: "Success",
		wantErr: faults.ErrorAppendCode,
	},
	{
		name:    "TestAppendCode_RealImplementation: BadRequestCodeExists",
		code:    response.BadRequest,
		message: "TestCode",
		wantErr: faults.ErrorAppendCode,
	},
	{
		name:    "TestAppendCode_RealImplementation: UnknownCode999",
		code:    999,
		message: "TestCode",
		wantErr: nil,
	},
	{
		name:    "TestAppendCode_RealImplementation: UnknownCode7892",
		code:    7892,
		message: "TestCode",
		wantErr: nil,
	},
}

var testsBuildError = []testBuildErrorResponse{
	{
		name:    "правильный ответ без ошибок",
		code:    response.ResourceCreated,
		title:   "пустой",
		message: "пустой",
		expected: response.ErrorResponse{
			Status: response.ResourceCreated,
			Error: response.ErrorDetail{
				Title:   "пустой",
				Message: "пустой",
			},
		},
	},
	{
		name:    "invalid title",
		code:    response.Success,
		title:   invalidMessage,
		message: "пустой",
		expected: response.ErrorResponse{
			Status: response.BadRequest,
			Error: response.ErrorDetail{
				Title:   faults.Invalid,
				Message: faults.TitleInvalid,
			},
		},
	},
	{
		name:    "invalid message",
		code:    response.Success,
		title:   "пустой",
		message: invalidMessage,
		expected: response.ErrorResponse{
			Status: response.BadRequest,
			Error: response.ErrorDetail{
				Title:   faults.Invalid,
				Message: faults.MessageInvalid,
			},
		},
	},
	{
		name:    "правильно отдает ответ с title = struct",
		code:    response.Success,
		title:   testStruct,
		message: "пустой",
		expected: response.ErrorResponse{
			Status: response.Success,
			Error: response.ErrorDetail{
				Title:   testStruct,
				Message: "пустой",
			},
		},
	},
	{
		name: "правильно отдает ответ с message = struct",

		code:    response.Success,
		title:   "пустой",
		message: testStruct,
		expected: response.ErrorResponse{
			Status: response.Success,
			Error: response.ErrorDetail{
				Title:   "пустой",
				Message: testStruct,
			},
		},
	},
	{
		name:    "правильно отдает ответ с message = struct и title = struct",
		code:    response.Success,
		title:   testStruct,
		message: testStruct,
		expected: response.ErrorResponse{
			Status: response.Success,
			Error: response.ErrorDetail{
				Title:   testStruct,
				Message: testStruct,
			},
		},
	},
}

var testsSendSuccess = []testSendSuccessAndError{
	{
		name:   "success",
		status: response.Success,
		output: float64(222),
		expected: response.SendSuccessOutput{
			Status: response.Success.String(),
			Output: float64(222),
		},
	},
	{
		name:   "testStruct",
		status: response.Success,
		output: testStruct,
		expected: response.SendSuccessOutput{
			Status: response.Success.String(),
			Output: map[string]interface{}{"Age": float64(666), "Name": "Aurivena"},
		},
	},
	{
		name:   "пустой ответ",
		status: response.Success,
		output: nil,
		expected: response.SendSuccessOutput{
			Status: response.Success.String(),
			Output: nil,
		},
	},
}

var testsSendError = []testSendSuccessAndError{
	{
		name:   "BadRequest",
		status: response.BadRequest,
		output: float64(222),
		expected: response.SendSuccessOutput{
			Status: response.BadRequest.String(),
			Output: response.ErrorDetail{
				Title:   "Ошибка",
				Message: float64(222),
			},
		},
	},
	{
		name:   "testStruct",
		status: response.InternalServerError,
		output: nil,
		expected: response.SendSuccessOutput{
			Status: response.InternalServerError.String(),
			Output: response.ErrorDetail{
				Title:   "Ошибка",
				Message: nil,
			},
		},
	},
}

func TestAppendCode_RealImplementation(t *testing.T) {
	impl := spond.NewImpl()

	for _, tt := range testsAppendCode {
		t.Run(fmt.Sprintf("%s%s", "TestAppendCode_RealImplementation: ", tt.name), func(t *testing.T) {

			err := impl.AppendCode(tt.code, tt.message)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr, "AppendCode должен вернуть error для code %d", tt.code)
			} else {
				assert.NoError(t, err, "AppendCode не должен вернуть ошибку для code %d", tt.code)
			}
		})
	}
}

func TestBuildError_RealImplementation(t *testing.T) {
	impl := spond.NewImpl()

	for _, tt := range testsBuildError {
		t.Run(fmt.Sprintf("%s%s", "TestBuildError_RealImplementation: ", tt.name), func(t *testing.T) {
			out := impl.BuildError(tt.code, tt.title, tt.message)
			assert.Equal(t, tt.expected, out, "BuildError should return expected response for %s", tt.name)
		})
	}
}

func TestSayHello_RealImplementation(t *testing.T) {
	t.Run("TestSayHello_RealImplementation: CallSayHello", func(t *testing.T) {
		var buf bytes.Buffer
		impl := spond.Impl{
			StatusMessages: make(map[response.StatusCode]string),
			Out:            &buf,
			Mu:             &sync.RWMutex{},
		}

		impl.SayHello()
		assert.Equal(t, "Hello it Spond!\n", buf.String(), `Функция должна была сказать Hello it Spond!\n`)
	})
}

func TestSendSuccessResponse_RealImplementation(t *testing.T) {
	impl := spond.NewImpl()

	for _, tt := range testsSendSuccess {
		t.Run(fmt.Sprintf("%s%s", "TestSendSuccessResponse_RealImplementation: ", tt.name), func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			impl.SendResponseSuccess(c, tt.status, tt.output)

			assert.Equal(t, http.StatusOK, w.Code)

			var response response.SendSuccessOutput
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			assert.Equal(t, tt.expected.Status, response.Status)
			assert.Equal(t, tt.expected.Output, response.Output)
		})

	}
}

func TestSendErrorResponse_RealImplementation(t *testing.T) {
	impl := spond.NewImpl()

	for _, tt := range testsSendError {
		t.Run(fmt.Sprintf("%s%s", "TestSendErrorResponse_RealImplementation: ", tt.name), func(t *testing.T) {
			gin.SetMode(gin.TestMode)
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			impl.SendResponseError(c, impl.BuildError(tt.status, "Ошибка", tt.output))

			assert.Equal(t, http.StatusOK, w.Code)

			var response testSendError
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			assert.Equal(t, tt.expected.Status, response.Status)
			assert.Equal(t, tt.expected.Output, response.Error)
		})

	}
}

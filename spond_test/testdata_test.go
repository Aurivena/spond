package spond

import (
	"github.com/gin-gonic/gin"
	"spond"
	"spond/faults"
)

type testBuildErrorResponse struct {
	name     string
	title    any
	message  any
	code     spond.StatusCode
	expected spond.ErrorResponse
}

type testAppendCode struct {
	name    string
	code    spond.StatusCode
	message string
	wantErr error
}

type testSendSuccessAndError struct {
	name     string
	c        *gin.Context
	status   spond.StatusCode
	output   any
	expected spond.SendSuccessOutput
}

type testSendError struct {
	Status string            `json:"status"`
	Error  spond.ErrorDetail `json:"error"`
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
		code:    spond.Success,
		message: "Success",
		wantErr: faults.ErrorAppendCode,
	},
	{
		name:    "TestAppendCode_RealImplementation: BadRequestCodeExists",
		code:    spond.BadRequest,
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
		code:    spond.ResourceCreated,
		title:   "пустой",
		message: "пустой",
		expected: spond.ErrorResponse{
			Status: spond.ResourceCreated,
			Error: spond.ErrorDetail{
				Title:   "пустой",
				Message: "пустой",
			},
		},
	},
	{
		name:    "invalid title",
		code:    spond.Success,
		title:   invalidMessage,
		message: "пустой",
		expected: spond.ErrorResponse{
			Status: spond.BadRequest,
			Error: spond.ErrorDetail{
				Title:   faults.Invalid,
				Message: faults.TitleInvalid,
			},
		},
	},
	{
		name:    "invalid message",
		code:    spond.Success,
		title:   "пустой",
		message: invalidMessage,
		expected: spond.ErrorResponse{
			Status: spond.BadRequest,
			Error: spond.ErrorDetail{
				Title:   faults.Invalid,
				Message: faults.MessageInvalid,
			},
		},
	},
	{
		name:    "правильно отдает ответ с title = struct",
		code:    spond.Success,
		title:   testStruct,
		message: "пустой",
		expected: spond.ErrorResponse{
			Status: spond.Success,
			Error: spond.ErrorDetail{
				Title:   testStruct,
				Message: "пустой",
			},
		},
	},
	{
		name: "правильно отдает ответ с message = struct",

		code:    spond.Success,
		title:   "пустой",
		message: testStruct,
		expected: spond.ErrorResponse{
			Status: spond.Success,
			Error: spond.ErrorDetail{
				Title:   "пустой",
				Message: testStruct,
			},
		},
	},
	{
		name:    "правильно отдает ответ с message = struct и title = struct",
		code:    spond.Success,
		title:   testStruct,
		message: testStruct,
		expected: spond.ErrorResponse{
			Status: spond.Success,
			Error: spond.ErrorDetail{
				Title:   testStruct,
				Message: testStruct,
			},
		},
	},
}

var testsSendSuccess = []testSendSuccessAndError{
	{
		name:   "success",
		status: spond.Success,
		output: float64(222),
		expected: spond.SendSuccessOutput{
			Status: spond.Success.String(),
			Output: float64(222),
		},
	},
	{
		name:   "testStruct",
		status: spond.Success,
		output: testStruct,
		expected: spond.SendSuccessOutput{
			Status: spond.Success.String(),
			Output: map[string]interface{}{"Age": float64(666), "Name": "Aurivena"},
		},
	},
	{
		name:   "пустой ответ",
		status: spond.Success,
		output: nil,
		expected: spond.SendSuccessOutput{
			Status: spond.Success.String(),
			Output: nil,
		},
	},
}

var testsSendError = []testSendSuccessAndError{
	{
		name:   "BadRequest",
		status: spond.BadRequest,
		output: float64(222),
		expected: spond.SendSuccessOutput{
			Status: spond.BadRequest.String(),
			Output: spond.ErrorDetail{
				Title:   "Ошибка",
				Message: float64(222),
			},
		},
	},
	{
		name:   "testStruct",
		status: spond.InternalServerError,
		output: nil,
		expected: spond.SendSuccessOutput{
			Status: spond.InternalServerError.String(),
			Output: spond.ErrorDetail{
				Title:   "Ошибка",
				Message: nil,
			},
		},
	},
}

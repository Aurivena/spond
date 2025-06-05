package test

import (
	"github.com/gin-gonic/gin"
	"net/http/httptest"
	"spond"
	"spond/models"
)

type testBuildErrorResponse struct {
	name     string
	c        *gin.Context
	title    any
	message  any
	code     spond.StatusCode
	expected models.ErrorResponse
}

type testAppendCode struct {
	name    string
	code    spond.StatusCode
	message string
	wantErr error
}

var c, _ = gin.CreateTestContext(httptest.NewRecorder())

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
		wantErr: spond.ErrorAppendCode,
	},
	{
		name:    "TestAppendCode_RealImplementation: BadRequestCodeExists",
		code:    spond.BadRequest,
		message: "TestCode",
		wantErr: spond.ErrorAppendCode,
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
		name:    "c == nil",
		c:       nil,
		code:    spond.Success,
		title:   "пустой",
		message: "пустой",
		expected: models.ErrorResponse{
			Status: spond.ContextIsNil.String(),
			Error: models.ErrorDetail{
				Title:   "",
				Message: "",
			},
		},
	},
	{
		name:    "правильный ответ без ошибок",
		c:       c,
		code:    spond.ResourceCreated,
		title:   "пустой",
		message: "пустой",
		expected: models.ErrorResponse{
			Status: spond.ResourceCreated.String(),
			Error: models.ErrorDetail{
				Title:   "пустой",
				Message: "пустой",
			},
		},
	},
	{
		name:    "invalid title",
		c:       c,
		code:    spond.Success,
		title:   invalidMessage,
		message: "пустой",
		expected: models.ErrorResponse{
			Status: spond.BadRequest.String(),
			Error: models.ErrorDetail{
				Title:   "Invalid",
				Message: "title invalid",
			},
		},
	},
	{
		name:    "invalid message",
		c:       c,
		code:    spond.Success,
		title:   "пустой",
		message: invalidMessage,
		expected: models.ErrorResponse{
			Status: spond.BadRequest.String(),
			Error: models.ErrorDetail{
				Title:   "Invalid",
				Message: "message invalid",
			},
		},
	},
	{
		name:    "правильно отдает ответ с title = struct",
		c:       c,
		code:    spond.Success,
		title:   testStruct,
		message: "пустой",
		expected: models.ErrorResponse{
			Status: spond.Success.String(),
			Error: models.ErrorDetail{
				Title:   testStruct,
				Message: "пустой",
			},
		},
	},
	{
		name:    "правильно отдает ответ с message = struct",
		c:       c,
		code:    spond.Success,
		title:   "пустой",
		message: testStruct,
		expected: models.ErrorResponse{
			Status: spond.Success.String(),
			Error: models.ErrorDetail{
				Title:   "пустой",
				Message: testStruct,
			},
		},
	},
	{
		name:    "правильно отдает ответ с message = struct и title = struct",
		c:       c,
		code:    spond.Success,
		title:   testStruct,
		message: testStruct,
		expected: models.ErrorResponse{
			Status: spond.Success.String(),
			Error: models.ErrorDetail{
				Title:   testStruct,
				Message: testStruct,
			},
		},
	},
}

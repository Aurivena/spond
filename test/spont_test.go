package test

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"spond"
	"spond/models"
	"testing"
)

func TestAppendCode_RealImplementation(t *testing.T) {
	impl := spond.NewImpl()

	tests := []struct {
		name    string
		code    spond.StatusCode
		message string
		wantErr error
	}{
		{
			name:    "UnknownCode666",
			code:    666,
			message: "TestCode",
			wantErr: nil,
		},
		{
			name:    "SuccessCodeExists",
			code:    spond.Success,
			message: "Success",
			wantErr: spond.ErrorAppendCode,
		},
		{
			name:    "BadRequestCodeExists",
			code:    spond.BadRequest,
			message: "TestCode",
			wantErr: spond.ErrorAppendCode,
		},
		{
			name:    "UnknownCode999",
			code:    999,
			message: "TestCode",
			wantErr: nil,
		},
		{
			name:    "UnknownCode7892",
			code:    7892,
			message: "TestCode",
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

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
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	invalidMessage := func() {}
	testStruct := struct {
		Name string
		Age  int
	}{
		Name: "Aurivena",
		Age:  666,
	}

	tests := []testBuildErrorResponse{
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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := impl.BuildError(tt.c, tt.code, tt.title, tt.message)
			assert.Equal(t, tt.expected, out, "BuildError should return expected response for %s", tt.name)
		})
	}
}

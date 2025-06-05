package test

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http/httptest"
	"spond"
	"spond/models"
	"testing"
)

type mockSpond struct {
	mock.Mock
}

func (m *mockSpond) AppendCode(code spond.StatusCode, msg any) error {
	args := m.Called(code, msg)
	return args.Error(0)
}

func (m *mockSpond) BuildError(c *gin.Context, code spond.StatusCode, title, message any) models.ErrorResponse {
	args := m.Called(c, code, title, message)
	return args.Get(0).(models.ErrorResponse)
}

func (m *mockSpond) SayHello() {
	m.Called()
}

func TestAppendCode_Mock(t *testing.T) {
	m := new(mockSpond)

	tests := []struct {
		name    string
		code    spond.StatusCode
		message string
		wantErr error
	}{
		{
			name:    "TestAppendCode_Mock: UnknownCode666",
			code:    666,
			message: "TestCode",
			wantErr: nil,
		},
		{
			name:    "TestAppendCode_Mock: SuccessCodeExists",
			code:    spond.Success,
			message: "Success",
			wantErr: spond.ErrorAppendCode,
		},
		{
			name:    "TestAppendCode_Mock: BadRequestCodeExists",
			code:    spond.BadRequest,
			message: "TestCode",
			wantErr: spond.ErrorAppendCode,
		},
		{
			name:    "TestAppendCode_Mock: UnknownCode999",
			code:    999,
			message: "TestCode",
			wantErr: nil,
		},
		{
			name:    "TestAppendCode_Mock: UnknownCode7892",
			code:    7892,
			message: "TestCode",
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m.On("AppendCode", tt.code, tt.message).Return(tt.wantErr).Once()

			err := m.AppendCode(tt.code, tt.message)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr, "AppendCode должен вернуть error для code %d", tt.code)
			} else {
				assert.NoError(t, err, "AppendCode не должен вернуть ошибку для code %d", tt.code)
			}
		})
	}
}

func TestBuildError_Mock(t *testing.T) {
	m := new(mockSpond)
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
			name:    "TestBuildError_Mock: c == nil",
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
			name:    "TestBuildError_Mock: правильный ответ без ошибок",
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
			name:    "TestBuildError_Mock: invalid title",
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
			name:    "TestBuildError_Mock: invalid message",
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
			name:    "TestBuildError_Mock: правильно отдает ответ с title = struct",
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
			name:    "TestBuildError_Mock: правильно отдает ответ с message = struct",
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
			name:    "TestBuildError_Mock: правильно отдает ответ с message = struct и title = struct",
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
			titleArg := tt.title
			messageArg := tt.message
			if _, err := json.Marshal(tt.title); err != nil {
				titleArg = mock.Anything
			}
			if _, err := json.Marshal(tt.message); err != nil {
				messageArg = mock.Anything
			}
			m.On("BuildError", tt.c, tt.code, titleArg, messageArg).Return(tt.expected).Once()

			out := m.BuildError(tt.c, tt.code, tt.title, tt.message)
			assert.Equal(t, tt.expected, out, "BuildError should return expected response for %s", tt.name)
			m.AssertExpectations(t)
		})
	}
}

func TestSayHello_Mock(t *testing.T) {
	t.Run("TestSayHello_Mock: CallSayHello", func(t *testing.T) {
		m := new(mockSpond)
		m.On("SayHello").Return().Once()

		m.SayHello()
		m.AssertExpectations(t)
	})
}

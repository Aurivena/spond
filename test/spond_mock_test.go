package test

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

	for _, tt := range testsAppendCode {
		t.Run(fmt.Sprintf("%s%s", "TestAppendCode_Mock: ", tt.name), func(t *testing.T) {
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

	for _, tt := range testsBuildError {
		t.Run(fmt.Sprintf("%s%s", "TestBuildError_Mock: ", tt.name), func(t *testing.T) {
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

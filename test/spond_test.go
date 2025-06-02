package test

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"spond"
	"spond/models"
	"sync"
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

func TestAppendCode(t *testing.T) {
	m := new(mockSpond)

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
			m.On("AppendCode", tt.code, tt.message).Return(tt.wantErr).Once()

			impl := spond.NewImpl()
			err := impl.AppendCode(tt.code, tt.message)
			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr, "AppendCode should return error for code %d", tt.code)
			} else {
				assert.NoError(t, err, "AppendCode should not return error for code %d", tt.code)
			}
		})
	}
}

func TestBuildError(t *testing.T) {

}

func TestSayHello(t *testing.T) {
	t.Run("CallSayHello", func(t *testing.T) {
		var buf bytes.Buffer
		impl := &spond.Impl{
			StatusMessages: make(map[spond.StatusCode]string),
			Out:            &buf,
			Mu:             &sync.RWMutex{},
		}

		impl.SayHello()
		assert.Equal(t, "Hello it Spond!\n", buf.String(), "Функция должна была написать `Hello it Spond!\n`")
	})
}

package test

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"spond"
	"sync"
	"testing"
)

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
			out := impl.BuildError(tt.c, tt.code, tt.title, tt.message)
			assert.Equal(t, tt.expected, out, "BuildError should return expected response for %s", tt.name)
		})
	}
}

func TestSayHello_RealImplementation(t *testing.T) {
	t.Run("TestSayHello_RealImplementation: CallSayHello", func(t *testing.T) {
		var buf bytes.Buffer
		impl := spond.Impl{
			StatusMessages: make(map[spond.StatusCode]string),
			Out:            &buf,
			Mu:             &sync.RWMutex{},
		}

		impl.SayHello()
		assert.Equal(t, "Hello it Spond!\n", buf.String(), `Функция должна была сказать Hello it Spond!\n`)
	})
}

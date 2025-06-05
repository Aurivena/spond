package test

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"spond"
	"testing"
)

func FuzzAppendCode(f *testing.F) {
	impl := spond.NewImpl()

	for _, tt := range testsAppendCode {
		f.Add(int(tt.code), tt.message)
	}

	f.Fuzz(func(t *testing.T, code int, message string) {
		err := impl.AppendCode(spond.StatusCode(code), message)

		if err != nil {
			assert.ErrorIs(t, err, spond.ErrorAppendCode,
				"AppendCode returned an unexpected error type for code %d, message '%s'", code, message)
		} else {
			err = impl.AppendCode(spond.StatusCode(code), message)
			assert.ErrorIs(t, err, spond.ErrorAppendCode,
				"AppendCode did not return ErrorAppendCode on second attempt for code %d, message '%s'", code, message)
		}
	})
}

func FuzzBuildError(f *testing.F) {
	impl := spond.NewImpl()

	f.Fuzz(func(t *testing.T, code int, title, message string) {
		var ctx, _ = gin.CreateTestContext(httptest.NewRecorder())
		out := impl.BuildError(ctx, spond.StatusCode(code), title, message)

		_ = out
	})
}

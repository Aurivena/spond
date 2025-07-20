package spond_test

import (
	"spond"
	"spond/faults"
	"spond/intertnal/response"
	"testing"

	"github.com/stretchr/testify/assert"
)

func FuzzAppendCode(f *testing.F) {
	impl := spond.NewImpl()

	for _, tt := range testsAppendCode {
		f.Add(int(tt.code), tt.message)
	}

	f.Fuzz(func(t *testing.T, code int, message string) {
		err := impl.AppendCode(response.StatusCode(code), message)

		if err != nil {
			assert.ErrorIs(t, err, faults.ErrorAppendCode,
				"AppendCode returned an unexpected error type for code %d, message '%s'", code, message)
		} else {
			err = impl.AppendCode(response.StatusCode(code), message)
			assert.ErrorIs(t, err, faults.ErrorAppendCode,
				"AppendCode did not return ErrorAppendCode on second attempt for code %d, message '%s'", code, message)
		}
	})
}

func FuzzBuildError(f *testing.F) {
	impl := response.NewImpl()

	f.Fuzz(func(t *testing.T, code int, title, message string) {
		out := impl.BuildError(response.StatusCode(code), title, message)

		_ = out
	})
}

package netstatus

import (
	"net/http"

	"google.golang.org/grpc/codes"
)

type Code int

const (
	CodeSuccess Code = iota
	CodeBadRequest
	CodeUnauthorized
	CodeNotFound
	CodeInternalError
)

func (c Code) HTTP() int {
	switch c {
	case CodeSuccess:
		return http.StatusOK
	case CodeBadRequest:
		return http.StatusBadRequest
	case CodeUnauthorized:
		return http.StatusUnauthorized
	case CodeNotFound:
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}

func (c Code) GRPC() codes.Code {
	switch c {
	case CodeSuccess:
		return codes.OK
	case CodeBadRequest:
		return codes.InvalidArgument
	case CodeUnauthorized:
		return codes.Unauthenticated
	case CodeNotFound:
		return codes.NotFound
	default:
		return codes.Internal
	}
}

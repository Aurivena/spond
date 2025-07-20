package response

import "fmt"

type StatusCode int

func (e StatusCode) MapToHTTPStatus() StatusCode {
	switch {
	case e >= 100 && e < 527:
		return e
	case e >= 4000 && e < 5000:
		return 400
	case e >= 5000 && e < 6000:
		return 500
	default:
		return 500
	}
}

func (e StatusCode) String() string {
	if str, ok := StatusMessages[e]; ok {
		return str
	}
	return fmt.Sprintf("unknown  StatusCode (%d)", e)
}

const (
	Success               StatusCode = 200
	ResourceCreated       StatusCode = 201
	NoContent             StatusCode = 204
	BadRequest            StatusCode = 400
	Unauthorized          StatusCode = 401
	Forbidden             StatusCode = 403
	NotFound              StatusCode = 404
	NotAcceptable         StatusCode = 406
	ConfirmationTimeout   StatusCode = 408
	ResourceAlreadyExists StatusCode = 409
	ResourceInTrash       StatusCode = 410
	BadHeader             StatusCode = 412
	UnsupportedMediaType  StatusCode = 415
	InternalServerError   StatusCode = 500
)

var StatusMessages = map[StatusCode]string{
	Success:               "Success",
	ResourceCreated:       "ResourceCreated",
	NoContent:             "NoContent",
	BadRequest:            "BadRequest",
	Unauthorized:          "Unauthorized",
	Forbidden:             "Forbidden",
	NotFound:              "NotFound",
	NotAcceptable:         "NotAcceptable",
	ConfirmationTimeout:   "ConfirmationTimeout",
	ResourceAlreadyExists: "ResourceAlreadyExists",
	ResourceInTrash:       "ResourceInTrash",
	BadHeader:             "BadHeader",
	UnsupportedMediaType:  "UnsupportedMediaType",
	InternalServerError:   "InternalServerError",
}

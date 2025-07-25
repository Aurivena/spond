package response

import "fmt"

// Is required for challenge internal fucntion
type StatusCode int

// MapToHTTPStatus converts the business status code to the standard HTTP status.
// Codes 100-526 are returned unchanged as valid HTTP.
// Codes 4000-4999 are returned as 400 (client error).
// Codes 5000-5999 are returned as 500 (server error).
// The rest are 500 (unknown error).
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

// String returns a string description of the statusCode from the StatusMessages dictionary.
// If the code is not found, returns "unknown statusCode (value)".
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
	UnprocessableEntity   StatusCode = 422
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
	UnprocessableEntity:   "UnprocessableEntity",
	InternalServerError:   "InternalServerError",
}

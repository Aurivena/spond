package envelope

import "fmt"

// Is required for challenge internal function
type StatusCode int

// String returns a string description of the statusCode from the StatusMessages dictionary.
// If the code is not found, returns "unknown statusCode (value)".
func (e StatusCode) String() string {
	if str, ok := StatusMessages[e]; ok {
		return str
	}
	return fmt.Sprintf("Unknown  StatusCode (%d)", e)
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

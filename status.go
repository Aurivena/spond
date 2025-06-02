package spond

import "fmt"

type StatusCode int

const (
	Success               StatusCode = 200
	ResourceCreated       StatusCode = 201
	NoContent             StatusCode = 204
	BadRequest            StatusCode = 400
	InvalidAccountStatus  StatusCode = 4001
	AccountIsActive       StatusCode = 4002
	Unauthorized          StatusCode = 401
	Forbidden             StatusCode = 403
	AccountIsInactive     StatusCode = 4031
	NotFound              StatusCode = 404
	NotAcceptable         StatusCode = 406
	ConfirmationTimeout   StatusCode = 408
	ResourceAlreadyExists StatusCode = 409
	ResourceInTrash       StatusCode = 410
	BadHeader             StatusCode = 412
	UnsupportedMediaType  StatusCode = 415
	InternalServerError   StatusCode = 500
	ContextIsNil          StatusCode = 5001
)

var statusMessages = map[StatusCode]string{
	Success:               "Success",
	ResourceCreated:       "ResourceCreated",
	NoContent:             "NoContent",
	BadRequest:            "BadRequest",
	InvalidAccountStatus:  "InvalidAccountStatus",
	AccountIsActive:       "AccountIsActive",
	Unauthorized:          "Unauthorized",
	Forbidden:             "Forbidden",
	AccountIsInactive:     "AccountIsInactive",
	NotFound:              "NotFound",
	NotAcceptable:         "NotAcceptable",
	ConfirmationTimeout:   "ConfirmationTimeout",
	ResourceAlreadyExists: "ResourceAlreadyExists",
	ResourceInTrash:       "ResourceInTrash",
	BadHeader:             "BadHeader",
	UnsupportedMediaType:  "UnsupportedMediaType",
	InternalServerError:   "InternalServerError",
	ContextIsNil:          "ContextIsNil",
}

func (e StatusCode) String() string {
	if str, ok := statusMessages[e]; ok {
		return str
	}
	return fmt.Sprintf("Неизвестный StatusCode (%d)", e)
}

package spond

import "fmt"

type SendErrorOutput struct {
	Status string      `json:"status"`
	Error  ErrorDetail `json:"error"`
}

type SendSuccessOutput struct {
	Status string `json:"status"`
	Output any    `json:"error"`
}

type ErrorResponse struct {
	Status StatusCode  `json:"status"`
	Error  ErrorDetail `json:"error"`
}

type ErrorDetail struct {
	Title   any `json:"title"`
	Message any `json:"message"`
}

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

func (e StatusCode) mapToHTTPStatus() StatusCode {
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
	if str, ok := statusMessages[e]; ok {
		return str
	}
	return fmt.Sprintf("Неизвестный StatusCode (%d)", e)
}

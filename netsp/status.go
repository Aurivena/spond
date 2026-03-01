package netsp

import (
	"fmt"
	"sync"
)

func appendCode(code int, message string) error {
	mu := sync.RWMutex{}

	mu.Lock()
	defer mu.Unlock()

	if _, exist := statusMessages[code]; exist {
		return fmt.Errorf("Status code %d already exists", code)
	}
	statusMessages[code] = message
	return nil
}

func isValid(code int) bool {
	_, ok := statusMessages[code]
	return ok
}

const (
	Success               int = 200
	ResourceCreated       int = 201
	NoContent             int = 204
	BadRequest            int = 400
	Unauthorized          int = 401
	Forbidden             int = 403
	NotFound              int = 404
	NotAcceptable         int = 406
	ConfirmationTimeout   int = 408
	ResourceAlreadyExists int = 409
	ResourceInTrash       int = 410
	BadHeader             int = 412
	UnsupportedMediaType  int = 415
	UnprocessableEntity   int = 422
	InternalServerError   int = 500
)

var statusMessages = map[int]string{
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

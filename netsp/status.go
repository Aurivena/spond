package netsp

import (
	"fmt"
	"net/http"
	"sync"
)

var mu sync.RWMutex

func appendCode(code int, message string) error {
	mu.Lock()
	defer mu.Unlock()

	if _, exist := statusMessages[code]; exist {
		return fmt.Errorf("status code %d already exists", code)
	}
	statusMessages[code] = message
	return nil
}

func isValid(code int) bool {
	mu.RLock()
	defer mu.RUnlock()
	_, ok := statusMessages[code]
	return ok
}

var statusMessages = map[int]string{
	http.StatusOK:                   "Success",
	http.StatusCreated:              "ResourceCreated",
	http.StatusNoContent:            "NoContent",
	http.StatusBadRequest:           "BadRequest",
	http.StatusUnauthorized:         "Unauthorized",
	http.StatusForbidden:            "Forbidden",
	http.StatusNotFound:             "NotFound",
	http.StatusNotAcceptable:        "NotAcceptable",
	http.StatusUnsupportedMediaType: "UnsupportedMediaType",
	http.StatusUnprocessableEntity:  "UnprocessableEntity",
	http.StatusInternalServerError:  "InternalServerError",
}

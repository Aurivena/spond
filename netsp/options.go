package netsp

import (
	"errors"
)

const (
	MaxTitleLength   = 256
	MaxMessageLength = 1024
)

var (
	ErrorAppendCode = errors.New("this code already exists")
	TitleInvalid    = errors.New("invalid value for title")
	MessageInvalid  = errors.New("invalid value for message")
	Invalid         = errors.New("invalid")
	UnknownStatus   = errors.New("unknown status")
	SolutionError   = errors.New("recheck limits for title and message pls :)")
)

// validate checks the length of the title and message.
// Returns the error when restrictions are violated.
func validate(title, message string) error {
	if len(title) == 0 || len(title) > MaxTitleLength {
		return TitleInvalid
	}
	if len(message) == 0 || len(message) > MaxMessageLength {
		return MessageInvalid
	}
	return nil
}

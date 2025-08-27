package envelope

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

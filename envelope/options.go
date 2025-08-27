package envelope

import (
	"errors"
)

const (
	MaxTitleLength   = 256
	MaxMessageLength = 1024
)

var (
	ErrorAppendCode = errors.New("This code already exists")
	TitleInvalid    = errors.New("Invalid value for title")
	MessageInvalid  = errors.New("Invalid value for message")
	Invalid         = errors.New("Invalid")
	UnknownStatus   = errors.New("Unknown status")
)

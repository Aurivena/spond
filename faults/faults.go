package faults

import "errors"

const (
	UnknownStatus = "UnknownStatus"
)

var (
	ErrorAppendCode = errors.New("этот код уже существует")
	TitleInvalid    = "title invalid"
	MessageInvalid  = "message invalid"
	Invalid         = "Invalid"
)

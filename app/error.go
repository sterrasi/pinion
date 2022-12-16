package app

import "strings"

type ErrorCode uint8

// Error interface describes an application error
type Error interface {
	error
	Code() ErrorCode
	CodeValue() string
	Cause() error
	GetContext() string
	SetContext(ctx string)
}

// ErrorImpl implements an Error
type ErrorImpl struct {
	code      ErrorCode
	codeValue string
	context   string
	cause     error
	metadata  map[string]string
	message   string
}

// Error returns a descriptor string that encapsulates all the error's metadata. This
// satisfies the 'error' interface
func (e *ErrorImpl) Error() string {

	var sb strings.Builder

	// add the error code
	sb.WriteString("[")
	sb.WriteString(e.codeValue)
	sb.WriteString("] ")

	// add the optional context
	if e.context != "" {
		sb.WriteString(e.context)
		sb.WriteString(": ")
	}

	// add the message
	sb.WriteString(e.message)

	// add the optional cause
	if e.cause != nil {
		sb.WriteString("; Cause=")
		sb.WriteString(e.cause.Error())
	}
	return sb.String()
}

// Code returns the associated error code
func (e *ErrorImpl) Code() ErrorCode {
	return e.code
}

// CodeValue returns the associated error code
func (e *ErrorImpl) CodeValue() string {
	return e.codeValue
}

// Cause returns the optional underlying error
func (e *ErrorImpl) Cause() error {
	return e.cause
}

// GetContext returns the error's optional context
func (e *ErrorImpl) GetContext() string {
	return e.context
}

// SetContext identifies the context of the error
func (e *ErrorImpl) SetContext(ctx string) {
	e.context = ctx
}

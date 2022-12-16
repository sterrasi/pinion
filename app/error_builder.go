package app

import "fmt"

// ErrorBuilder builds an Error
type ErrorBuilder struct {
	code      ErrorCode
	codeValue string
	fn        func(impl *ErrorImpl) Error
	context   string
	cause     error
	metadata  map[string]string
	message   string
}

// NewErrorBuilderWithFactory creates a new ErrorBuilder
func NewErrorBuilderWithFactory(code ErrorCode, codeValue string,
	fn func(impl *ErrorImpl) Error) *ErrorBuilder {

	return &ErrorBuilder{
		code:      code,
		codeValue: codeValue,
		fn:        fn}
}

// NewErrorBuilder creates a new ErrorBuilder
func NewErrorBuilder(code ErrorCode, codeValue string) *ErrorBuilder {

	return &ErrorBuilder{
		code:      code,
		codeValue: codeValue}
}

// Cause sets the optional error that caused this error
func (b *ErrorBuilder) Cause(err error) *ErrorBuilder {
	b.cause = err
	return b
}

// Context sets a contextual string that defines the operation in-process when the error occurred
func (b *ErrorBuilder) Context(context string) *ErrorBuilder {
	b.context = context
	return b
}

// Str sets a string key and value to the metadata associated with this error
func (b *ErrorBuilder) Str(key string, val string) *ErrorBuilder {
	if len(key) != 0 {
		if b.metadata == nil {
			b.metadata = make(map[string]string)
		}
		b.metadata[key] = val
	}
	return b
}

// Msg sets the message for the Error and creates it
func (b *ErrorBuilder) Msg(msg string) Error {
	b.message = msg
	if b.fn != nil {
		return b.fn(&ErrorImpl{
			code:      b.code,
			codeValue: b.codeValue,
			context:   b.context,
			cause:     b.cause,
			metadata:  b.metadata,
			message:   b.message,
		})
	}
	return &ErrorImpl{
		code:      b.code,
		codeValue: b.codeValue,
		context:   b.context,
		cause:     b.cause,
		metadata:  b.metadata,
		message:   b.message,
	}
}

// Msgf sets the message for the Error using message formatting and creates it
func (b *ErrorBuilder) Msgf(format string, args ...any) Error {

	if len(args) > 0 {
		return b.Msg(fmt.Sprintf(format, args...))
	}
	return b.Msg(format)
}

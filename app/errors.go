package app

const UnknownErrorCode ErrorCode = 0

// InternalErrorCode relates to a general internal server error that should be avoided if
// a more specific one can be chosen
const InternalErrorCode ErrorCode = 1

func BuildInternalError() *ErrorBuilder {
	return NewErrorBuilder(InternalErrorCode, "internal")
}
func NewInternalError(format string, args ...any) Error {
	return BuildInternalError().Msgf(format, args...)
}

// SystemConfigurationErrorCode signifies a server error that keeps the server from starting.  It is
// related to an issue that can be fixed in the software's configuration
const SystemConfigurationErrorCode ErrorCode = 2

func BuildSysConfigError() *ErrorBuilder {
	return NewErrorBuilder(SystemConfigurationErrorCode, "system-configuration")
}
func NewSysConfigError(format string, args ...any) Error {
	return BuildSysConfigError().Msgf(format, args...)
}

// ServiceUnavailableErrorCode signifies that either the server or one of its dependencies is not able to service the
// request.
const ServiceUnavailableErrorCode ErrorCode = 3

func BuildSvcUnavailableError() *ErrorBuilder {
	return NewErrorBuilder(ServiceUnavailableErrorCode, "service-unavailable")
}
func NewSvcUnavailableError(format string, args ...any) Error {
	return BuildSvcUnavailableError().Msgf(format, args...)
}

// IllegalArgumentError relates to an internal server error that means an internal argument check failed.  This type
// of error usually signifies a bug in the software
const IllegalArgumentError ErrorCode = 4

func BuildIllegalArgumentError() *ErrorBuilder {
	return NewErrorBuilder(IllegalArgumentError, "illegal-argument")
}
func NewIllegalArgumentError(format string, args ...any) Error {
	return BuildIllegalArgumentError().Msgf(format, args...)
}

// ValidationErrorCode signifies a client level error that means the data provided by the client to the server
// is invalid
const ValidationErrorCode ErrorCode = 5

func BuildValidationError() *ErrorBuilder {
	return NewErrorBuilder(ValidationErrorCode, "validation")
}
func NewValidationError(format string, args ...any) Error {
	return BuildValidationError().Msgf(format, args...)
}

// IllegalStateErrorCode relates to a client level error that signifies the operation asked of the server
// cannot be performed because it is not in the proper state
const IllegalStateErrorCode ErrorCode = 6

func BuildIllegalStateError() *ErrorBuilder {
	return NewErrorBuilder(IllegalStateErrorCode, "illegal-state")
}
func NewIllegalStateError(format string, args ...any) Error {
	return BuildIllegalStateError().Msgf(format, args...)
}

// NotFoundErrorCode relates to a client level error where an entity is referenced by the client that does not exist
const NotFoundErrorCode ErrorCode = 7

func BuildNotFoundError() *ErrorBuilder {
	return NewErrorBuilder(NotFoundErrorCode, "not-found")
}
func NewNotFoundError(format string, args ...any) Error {
	return BuildNotFoundError().Msgf(format, args...)
}

// AlreadyExistsErrorCode relates to a client level error where the result of an operation is to produce a new entity
// but the entity already exists in the system
const AlreadyExistsErrorCode ErrorCode = 8

func BuildAlreadyExistsError() *ErrorBuilder {
	return NewErrorBuilder(AlreadyExistsErrorCode, "already-exists")
}
func NewAlreadyExistsError(format string, args ...any) Error {
	return BuildAlreadyExistsError().Msgf(format, args...)
}

// IOErrorCode relates to an error while trying to access a resource like a file
const IOErrorCode ErrorCode = 9

func BuildIOError() *ErrorBuilder {
	return NewErrorBuilder(IOErrorCode, "io")
}
func NewIOError(format string, args ...any) Error {
	return BuildIOError().Msgf(format, args...)
}

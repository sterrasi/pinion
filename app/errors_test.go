package app

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestErrorCodeMappings(t *testing.T) {
	assert.Equal(t, InternalErrorCode, NewInternalError("error").Code())
	assert.Equal(t, InternalErrorCode, BuildInternalError().Msg("error").Code())

	assert.Equal(t, SystemConfigurationErrorCode, NewSysConfigError("error").Code())
	assert.Equal(t, SystemConfigurationErrorCode, BuildSysConfigError().Msg("error").Code())

	assert.Equal(t, ServiceUnavailableErrorCode, NewSvcUnavailableError("error").Code())
	assert.Equal(t, ServiceUnavailableErrorCode, BuildSvcUnavailableError().Msg("error").Code())

	assert.Equal(t, IllegalArgumentError, NewIllegalArgumentError("error").Code())
	assert.Equal(t, IllegalArgumentError, BuildIllegalArgumentError().Msg("error").Code())

	assert.Equal(t, ValidationErrorCode, NewValidationError("error").Code())
	assert.Equal(t, ValidationErrorCode, BuildValidationError().Msg("error").Code())

	assert.Equal(t, IllegalStateErrorCode, NewIllegalStateError("error").Code())
	assert.Equal(t, IllegalStateErrorCode, BuildIllegalStateError().Msg("error").Code())

	assert.Equal(t, NotFoundErrorCode, NewNotFoundError("error").Code())
	assert.Equal(t, NotFoundErrorCode, BuildNotFoundError().Msg("error").Code())

	assert.Equal(t, AlreadyExistsErrorCode, NewAlreadyExistsError("error").Code())
	assert.Equal(t, AlreadyExistsErrorCode, BuildAlreadyExistsError().Msg("error").Code())
}

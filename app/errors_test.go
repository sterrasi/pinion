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

	assert.Equal(t, ValidationErrorCode, NewValidationErrorCode("error").Code())
	assert.Equal(t, ValidationErrorCode, BuildValidationErrorCode().Msg("error").Code())

	assert.Equal(t, IllegalStateErrorCode, NewIllegalStateErrorCode("error").Code())
	assert.Equal(t, IllegalStateErrorCode, BuildIllegalStateErrorCode().Msg("error").Code())

	assert.Equal(t, NotFoundErrorCode, NewNotFoundErrorCode("error").Code())
	assert.Equal(t, NotFoundErrorCode, BuildNotFoundErrorCode().Msg("error").Code())

	assert.Equal(t, AlreadyExistsErrorCode, NewAlreadyExistsErrorCode("error").Code())
	assert.Equal(t, AlreadyExistsErrorCode, BuildAlreadyExistsErrorCode().Msg("error").Code())
}

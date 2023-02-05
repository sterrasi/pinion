package app

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestAppErrorMessageFormat tests that the Error() if a fully hydrated app.Error is well-formed
func TestAppErrorMessageFormat(t *testing.T) {

	err := BuildInternalError().Context("some context").
		Cause(errors.New("the cause")).
		Msgf("there was a %s", "problem")
	assert.NotNil(t, err)
	assert.Equal(t, InternalErrorCode, err.Code())
	assert.Equal(t, "internal", err.CodeValue())
	assert.Equal(t, "the cause", err.Cause().Error())
	assert.Equal(t, "[internal] some context: there was a problem; Cause=the cause", err.Error())
}

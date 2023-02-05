package http

import (
	"github.com/sterrasi/pinion/app"
	nethttp "net/http"
)

// GetHttpStatusCode returns the http status code for the given apperr.Error. If the error code is unknown
// then false is returned
func GetHttpStatusCode(err app.Error) (bool, int) {

	switch err.Code() {
	case app.InternalErrorCode:
		fallthrough
	case app.IllegalArgumentError:
		fallthrough
	case app.SystemConfigurationErrorCode:
		return true, nethttp.StatusInternalServerError

	case app.ServiceUnavailableErrorCode:
		return true, nethttp.StatusServiceUnavailable

	case app.ValidationErrorCode:
		return true, nethttp.StatusBadRequest

	case app.AlreadyExistsErrorCode:
		fallthrough
	case app.IllegalStateErrorCode:
		return true, nethttp.StatusConflict

	case app.NotFoundErrorCode:
		return true, nethttp.StatusNotFound

	default:
		return false, 0
	}
}

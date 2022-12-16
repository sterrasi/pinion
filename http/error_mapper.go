package http

import (
	"github.com/sterrasi/pinion/app"
	gohttp "net/http"
)

// GetHttpStatusCode returns the http status code for the given app.Error. If the error code is unknown
// then false is returned
func GetHttpStatusCode(err app.Error) (bool, int) {

	switch err.Code() {
	case app.InternalErrorCode:
		fallthrough
	case app.IllegalArgumentError:
		fallthrough
	case app.SystemConfigurationErrorCode:
		return true, gohttp.StatusInternalServerError

	case app.ServiceUnavailableErrorCode:
		return true, gohttp.StatusServiceUnavailable

	case app.ValidationErrorCode:
		return true, gohttp.StatusBadRequest

	case app.AlreadyExistsErrorCode:
		fallthrough
	case app.IllegalStateErrorCode:
		return true, gohttp.StatusConflict

	case app.NotFoundErrorCode:
		return true, gohttp.StatusNotFound

	default:
		return false, 0
	}
}

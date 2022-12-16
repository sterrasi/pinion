package grpc

import (
	"github.com/sterrasi/pinion/app"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ToStatus will return the corresponding status.Error for the given app.Error. If the app.Error cannot
// be identified then false will be returned
func ToStatus(err app.Error) (bool, error) {
	switch err.Code() {
	case app.InternalErrorCode:
		fallthrough
	case app.IllegalArgumentError:
		fallthrough
	case app.SystemConfigurationErrorCode:
		return true, status.Error(codes.Internal, err.Error())

	case app.ServiceUnavailableErrorCode:
		return true, status.Error(codes.Unavailable, err.Error())

	case app.ValidationErrorCode:
		return true, status.Error(codes.InvalidArgument, err.Error())

	case app.IllegalStateErrorCode:
		return true, status.Error(codes.FailedPrecondition, err.Error())

	case app.NotFoundErrorCode:
		return true, status.Error(codes.NotFound, err.Error())

	case app.AlreadyExistsErrorCode:
		return true, status.Error(codes.AlreadyExists, err.Error())

	default:
		return false, nil
	}
}

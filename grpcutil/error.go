package grpcutil

import (
	"errors"
	"github.com/alexandria-oss/core/exception"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ResponseErr writes the required error's gRPC code and message using error
// from the gRCP action
func ResponseErr(err error) error {
	return status.Error(ErrorToCode(err), exception.GetErrorDescription(err))
}

// ErrorToCode returns gRPC code depending on the error given
func ErrorToCode(err error) codes.Code {
	switch {
	case errors.Is(err, exception.EntityNotFound) || errors.Is(err, exception.EntitiesNotFound):
		return codes.NotFound
	case errors.Is(err, exception.RequiredField) || errors.Is(err, exception.InvalidFieldFormat) ||
		errors.Is(err, exception.InvalidID) || errors.Is(err, exception.EmptyBody):
		return codes.InvalidArgument
	case errors.Is(err, exception.InvalidFieldRange):
		return codes.OutOfRange
	case errors.Is(err, exception.EntityExists):
		return codes.AlreadyExists
	default:
		return codes.Internal
	}
}
package app_error

import (
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrUserNotFound      = errors.New("User Not Found")
	ErrEmailAlreadyExist = errors.New("Email is already exist")
	ErrInvalidInput      = errors.New("invalid input")
)

func ToGRPCError(err error) error {
	switch {
	case errors.Is(err, ErrUserNotFound):
		return status.Errorf(codes.NotFound, "User Is Not Exist Or Wrong Credentials")
	case errors.Is(err, ErrEmailAlreadyExist):
		return status.Errorf(codes.AlreadyExists, "Email Already Exists")
	case errors.Is(err, ErrInvalidInput):
		return status.Errorf(codes.InvalidArgument, "%s", err.Error())
	default:
		return status.Errorf(codes.Internal, "an internal error occurred")
	}
}

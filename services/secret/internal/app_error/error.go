package app_error

import (
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrSecretAlreadyExists = errors.New("secret already exists")
	ErrSecretNotFound      = errors.New("secret not found")
	ErrOrgNotFound         = errors.New("organization not found")
	ErrInvalidInput        = errors.New("invalid input")
	ErrForeignKeyViolation = errors.New("foreign key violation")
)

func ToGRPCError(err error) error {
	switch {
	case errors.Is(err, ErrSecretAlreadyExists):
		return status.Errorf(codes.AlreadyExists, "a secret with this name already exists in the organization")
	case errors.Is(err, ErrSecretNotFound):
		return status.Errorf(codes.NotFound, "secret not found")
	case errors.Is(err, ErrOrgNotFound):
		return status.Errorf(codes.NotFound, "organization not found")
	case errors.Is(err, ErrInvalidInput):
		return status.Errorf(codes.InvalidArgument, "%s", err.Error())
	default:
		// Never leak internal details to the caller.
		return status.Errorf(codes.Internal, "an internal error occurred")
	}
}

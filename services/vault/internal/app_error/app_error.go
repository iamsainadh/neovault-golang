package app_error

import (
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrInvalidInput      = errors.New("invalid input")
	ErrOrgNotFound       = errors.New("organization not found")
	ErrDuplicateVaultKey = errors.New("vault key version already exists")
)

func ToGRPCError(err error) error {

	switch {
	case errors.Is(err, ErrOrgNotFound):
		return status.Errorf(codes.NotFound, "organization not found")
	case errors.Is(err, ErrInvalidInput):
		return status.Errorf(codes.InvalidArgument, "%s", err.Error())
	case errors.Is(err, ErrDuplicateVaultKey):
		return status.Errorf(
			codes.AlreadyExists,
			"vault key version already exists for this organization",
		)
	default:
		return status.Errorf(codes.Internal, "an internal error occurred")
	}
}

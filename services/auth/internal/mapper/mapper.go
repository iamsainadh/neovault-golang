package mapper

import (
	"github.com/iamsainadh/neovault/auth/internal/model"
	pb "github.com/iamsainadh/neovault/auth/proto"
)

// Mapper the Grpc to Go model

func MapUserGrpcToModel(req *pb.CreateUserRequest) model.CreateUserInput {
	return model.CreateUserInput{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}
}

func MapOrganizationGrpcToModel(req *pb.CreateOrganizationRequest) model.CreateOrganizationInput {
	return model.CreateOrganizationInput{
		Name:   req.Name,
		UserId: req.UserId,
	}
}

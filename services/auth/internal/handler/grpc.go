package handler

import (
	"context"

	"github.com/iamsainadh/neovault/auth/internal/app_error"
	"github.com/iamsainadh/neovault/auth/internal/mapper"
	"github.com/iamsainadh/neovault/auth/internal/model"
	"github.com/iamsainadh/neovault/auth/internal/service"

	pb "github.com/iamsainadh/neovault/auth/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthHandler struct {
	pb.UnimplementedUserServiceServer
	userService         service.UserService
	organizationService service.OrganizationService
}

func NewAuthHandler(svc service.UserService, orgSvc service.OrganizationService) *AuthHandler {
	return &AuthHandler{
		userService:         svc,
		organizationService: orgSvc,
	}
}

func (h AuthHandler) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}

	if req.Email == "" || req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "email and password are required")
	}

	//! Map the input grpc request to user input model
	user := mapper.MapUserGrpcToModel(req)

	createdUser, err := h.userService.CreateUser(ctx, user)

	if err != nil {
		return nil, app_error.ToGRPCError(err)
	}

	return &pb.CreateUserResponse{
		UserId: createdUser.UserId,
		Token:  createdUser.Token,
	}, nil
}

func (h AuthHandler) CreateOrganization(ctx context.Context, req *pb.CreateOrganizationRequest) (*pb.CreateOrganizationResponse, error) {
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "Name is required")
	}

	organization := mapper.MapOrganizationGrpcToModel(req)

	createdOrg, err := h.organizationService.CreateOrganization(ctx, organization)

	if err != nil {
		return nil, app_error.ToGRPCError(err)
	}

	return &pb.CreateOrganizationResponse{
		Id:   createdOrg.Id,
		Name: createdOrg.Name,
	}, nil
}

func (h AuthHandler) GetUserLogin(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	if req.Email == "" || req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "Credentials are incorrect")
	}

	user := &model.LoginUserInput{
		Email:    req.Email,
		Password: req.Password,
	}

	existingUser, err := h.userService.LoginUser(ctx, user)

	if err != nil {
		return nil, app_error.ToGRPCError(err)
	}

	return &pb.GetUserResponse{
		UserId: existingUser.UserId,
		Token:  existingUser.Token,
	}, nil
}

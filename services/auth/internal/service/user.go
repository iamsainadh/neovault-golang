package service

import (
	"context"
	"os"

	"github.com/iamsainadh/neovault/auth/internal/app_error"
	"github.com/iamsainadh/neovault/auth/internal/model"
	"github.com/iamsainadh/neovault/auth/internal/repository"
	"github.com/iamsainadh/neovault/utils"
)

type UserService interface {
	CreateUser(ctx context.Context, input model.CreateUserInput) (model.CreateUserResponse, error)
	//! Log The User
	LoginUser(ctx context.Context, input *model.LoginUserInput) (*model.LoginUserResponse, error)
	//! Token Generate
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (us *userService) CreateUser(ctx context.Context, input model.CreateUserInput) (model.CreateUserResponse, error) {
	hashedPassword, err := utils.GenerateHash(input.Password)
	if err != nil {
		return model.CreateUserResponse{}, err
	}
	// User Creation
	user := model.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: hashedPassword,
	}

	createdUser, err := us.repo.CreateUser(ctx, user)

	if err != nil {
		return model.CreateUserResponse{}, err
	}

	//TODO: If User Have Organization Account return
	token, err := utils.GenerateJwt(os.Getenv("JWT_SECRET"), createdUser.ID, "user", "")

	//! Generate Jwt Token and return
	return model.CreateUserResponse{
		UserId: createdUser.ID,
		Token:  token,
	}, nil
}

func (us *userService) LoginUser(ctx context.Context, input *model.LoginUserInput) (*model.LoginUserResponse, error) {
	existingUser, err := us.repo.LoginUser(ctx, input)

	if err != nil {
		return &model.LoginUserResponse{}, err
	}

	isPasswordCorrect, err := utils.VerifyHash(input.Password, existingUser.Password)

	if err != nil {
		return &model.LoginUserResponse{}, err
	}

	var loginRes model.LoginUserResponse

	if isPasswordCorrect {
		loginRes.UserId = existingUser.UserId
		//! Generate Token
		// TODO:  Check User is exist in Organization Member Table if exist return the token with that role if not return the token as user
		token, err := utils.GenerateJwt(os.Getenv("JWT_SECRET"), existingUser.UserId, "user", "")

		if err != nil {
			return &model.LoginUserResponse{}, err
		}
		loginRes.Token = token
	} else {
		return &model.LoginUserResponse{}, app_error.ErrUserNotFound
	}

	return &loginRes, nil
}

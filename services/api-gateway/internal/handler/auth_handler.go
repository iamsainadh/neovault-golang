package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/iamsainadh/neovault/api-gateway/internal/middleware"
	"github.com/iamsainadh/neovault/api-gateway/internal/model"
	auth "github.com/iamsainadh/neovault/auth/proto"
)

type AuthHandler struct {
	client auth.UserServiceClient
}

func NewAuthHandler(auth auth.UserServiceClient) *AuthHandler {
	return &AuthHandler{
		client: auth,
	}
}

// Create User
func (h *AuthHandler) CreateUser(ctx *gin.Context) {
	var user model.CreateUserRequest

	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"status":  false,
			"message": getValidationError(err),
		})
		return
	}

	createdUser, err := h.client.CreateUser(ctx.Request.Context(), &auth.CreateUserRequest{
		Name:     user.Name,
		Email:    user.Email,
		Password: user.Password,
	})

	if err != nil {
		middleware.HandleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"code":     http.StatusCreated,
		"message":  "User Created Successfully",
		"response": createdUser,
		"status":   true,
	})
}

// Login User
func (h *AuthHandler) LoginUser(ctx *gin.Context) {
	var user model.LoginUserRequest

	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"status":  false,
			"message": getValidationError(err),
		})
		return
	}

	loginUser, err := h.client.GetUserLogin(ctx.Request.Context(), &auth.GetUserRequest{
		Email:    user.Email,
		Password: user.Password,
	})

	if err != nil {
		middleware.HandleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":     http.StatusOK,
		"message":  "Login Successfully",
		"response": loginUser,
		"status":   true,
	})

}

// Create Organization
func (h *AuthHandler) CreateOrganization(ctx *gin.Context) {
	var organizationReq model.CreateOrganizationRequest

	if err := ctx.ShouldBindJSON(&organizationReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"status":  false,
			"message": getValidationError(err),
		})
		return
	}

	userId := ctx.GetString("user_id")
	createdOrg, err := h.client.CreateOrganization(ctx.Request.Context(), &auth.CreateOrganizationRequest{
		Name:   organizationReq.OrgName,
		UserId: userId,
	})
	if err != nil {
		fmt.Println(err)
		middleware.HandleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"code":     http.StatusOK,
		"message":  "Organization Created Successfully",
		"response": createdOrg,
		"status":   true,
	})

}

// User Login

// User Token Request

// Validator
func getValidationError(err error) string {
	var ve validator.ValidationErrors

	if errors.As(err, &ve) {
		field := ve[0].Field()
		tag := ve[0].Tag()

		switch field {
		case "Name":
			if tag == "required" {
				return "name is required"
			}
		case "Email":
			if tag == "required" {
				return "email is required"
			}
			if tag == "email" {
				return "please enter a valid email"
			}
		case "Password":
			if tag == "required" {
				return "password is required"
			}
			if tag == "min" {
				return "password must be at least 8 characters"
			}
		case "OrgName":
			if tag == "required" {
				return "Organization Name is required"
			}

		case "UserId":
			if tag == "required" {
				return "User Id is required"
			}
		}
	}

	return "invalid request"
}

package model

type CreateUserRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type LoginUserRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type CreateOrganizationRequest struct {
	OrgName string `json:"org_name" binding:"required"`
	// UserId  string `json:"user_id" binding:"required"`
}

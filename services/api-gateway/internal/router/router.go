package router

import (
	"github.com/gin-gonic/gin"
	"github.com/iamsainadh/neovault/api-gateway/internal/config"
	"github.com/iamsainadh/neovault/api-gateway/internal/middleware"
	"github.com/iamsainadh/neovault/api-gateway/internal/repository"
)

func NewRouter(
	r *gin.Engine,
	jwtSecret string,
	cfg *config.Config,
	repo repository.Repository,
) {
	// Public Routes
	api := r.Group("/api/v1")

	//------------------------------------------------------------------------------ Auth Unprotected Routes ------------------------------------------------------------------------------
	authHandler := AuthRouter(cfg, api)

	protected := api.Group("/")

	protected.Use(middleware.JWTAuthMiddleware(jwtSecret))

	protected.POST("/create-org", middleware.RBACMiddleware("user"), authHandler.CreateOrganization)

	// Set Role Here (Not using from Role From JWT Because Of Security Reason if a user is removed he still have token to delete instead we are adding the role fetch from db based on user_id)
	protected.Use(middleware.OrganizationRoleMiddleware(repo))

	//------------------------------------------------------------------------------ Organization Routes ------------------------------------------------------------------------------
	org := protected.Group("/org")

	AdminRoutes(cfg, org)
}

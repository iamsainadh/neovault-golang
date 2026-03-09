package router

import (
	"github.com/gin-gonic/gin"
	"github.com/iamsainadh/neovault/api-gateway/internal/config"
	"github.com/iamsainadh/neovault/api-gateway/internal/middleware"
)

func AdminRoutes(
	cfg *config.Config,
	router *gin.RouterGroup,
) {
	
	router.POST("/:orgId/create-secret", middleware.RBACMiddleware("admin"))
}

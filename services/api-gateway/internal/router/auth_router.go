package router

import (
	"github.com/bytedance/gopkg/util/logger"
	"github.com/gin-gonic/gin"
	"github.com/iamsainadh/neovault/api-gateway/internal/config"
	"github.com/iamsainadh/neovault/api-gateway/internal/handler"
	client "github.com/iamsainadh/neovault/shared/grpcclients/auth"
)

func AuthRouter(
	cfg *config.Config,
	router *gin.RouterGroup,
) *handler.AuthHandler {
	authClient, err := client.NewAuthClient(cfg.AUTH_PORT)
	if err != nil {
		logger.Info(err)
	}
	authHandler := handler.NewAuthHandler(authClient)
	auth := router.Group("/auth")
	auth.POST("/create-user", authHandler.CreateUser)
	auth.POST("/user-login", authHandler.LoginUser)

	return authHandler
}

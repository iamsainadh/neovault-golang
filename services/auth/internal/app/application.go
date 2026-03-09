package app

import (
	"log/slog"

	"github.com/iamsainadh/neovault/auth/internal/config"
	"github.com/iamsainadh/neovault/auth/internal/handler"
	"github.com/iamsainadh/neovault/auth/internal/repository"
	"github.com/iamsainadh/neovault/auth/internal/service"
	client "github.com/iamsainadh/neovault/shared/grpcclients/vault"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Application struct {
	Handler *handler.AuthHandler
	DB      *pgxpool.Pool
	Cfg     *config.Config
}

func NewApplication(pool *pgxpool.Pool, cfg *config.Config) *Application {
	userRepo := repository.NewUserRepository(pool)

	organizationRepo := repository.NewOrganizationRepository(pool)
	memberRepo := repository.NewMemberRepository(pool)

	userSvc := service.NewUserService(userRepo)

	// Vault Client
	vaultClient, err := client.NewVaultClient(":" + cfg.VAULT_PORT)

	if err != nil {
		// TODO : Delete The Organization if vault client gives err (SAGA PATTERN)
		slog.Warn("vault unavailable", slog.Any("error", err))
	}

	orgSvc := service.NewOrganizationService(pool, organizationRepo, memberRepo, vaultClient)

	handler := handler.NewAuthHandler(userSvc, orgSvc)

	return &Application{
		Handler: handler,
		DB:      pool,
	}
}

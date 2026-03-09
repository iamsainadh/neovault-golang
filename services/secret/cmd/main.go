package main

import (
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	pb "github.com/iamsainadh/neovault/secret/proto"
	client "github.com/iamsainadh/neovault/shared/grpcclients/vault"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"

	"github.com/iamsainadh/neovault/secret/internal/config"
	"github.com/iamsainadh/neovault/secret/internal/database"
	"github.com/iamsainadh/neovault/secret/internal/handler"
	"github.com/iamsainadh/neovault/secret/internal/repository"
	"github.com/iamsainadh/neovault/secret/internal/service"
	"google.golang.org/grpc"
)

func main() {
	//--
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	slog.SetDefault(logger)
	//--

	cfg := config.LoadConfig()

	//! DB
	pool, err := database.PostgresConnectionPool(cfg)

	if err != nil {
		logger.Error("failed to connect to database", slog.Any("error", err))
		os.Exit(1)
	}
	defer pool.Close()
	logger.Info("database connected")

	vaultClient, err := client.NewVaultClient(cfg.VAULT_PORT)

	secretRepo := repository.NewSecretRepository(pool)

	secretVersionRepo := repository.NewSecretVersionRepository(pool)

	secretVersionService := service.NewSecretVersionService(secretVersionRepo, vaultClient)

	secretService := service.NewSecretService(secretRepo, secretVersionService)

	h := handler.NewSecretHandler(secretService, logger)

	versionHandler := handler.NewSecretVersionHandler(secretVersionService, logger)

	// -- Create GRPC Server
	grpcServer := grpc.NewServer(
	// TODO : Create the interceptor
	)

	healthServer := health.NewServer()

	healthpb.RegisterHealthServer(grpcServer, healthServer)

	healthServer.SetServingStatus(
		"",
		healthpb.HealthCheckResponse_SERVING,
	)

	pb.RegisterSecretServiceServer(grpcServer, h)
	pb.RegisterSecretVersionServer(grpcServer, versionHandler)

	addr := ":" + cfg.GRPCPort

	listener, err := net.Listen("tcp", addr)

	if err != nil {
		logger.Error("failed to listen", slog.String("addr", addr), slog.Any("error", err))
		os.Exit(1)
	}

	//-- Graceful Shutdown

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		logger.Info("secret gRPC server starting", slog.String("addr", addr))
		if err := grpcServer.Serve(listener); err != nil {
			logger.Error("gRPC server error", slog.Any("error", err))
			os.Exit(1)
		}
	}()

	<-quit

	logger.Info("Shutdown signal received")

	healthServer.SetServingStatus(
		"",
		healthpb.HealthCheckResponse_NOT_SERVING,
	)
	grpcServer.GracefulStop()
	logger.Info("server stopped cleanly")
}

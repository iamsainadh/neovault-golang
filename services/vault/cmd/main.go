package main

import (
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/iamsainadh/neovault/vault/internal/config"
	"github.com/iamsainadh/neovault/vault/internal/database"
	"github.com/iamsainadh/neovault/vault/internal/handler"
	"github.com/iamsainadh/neovault/vault/internal/repository"
	"github.com/iamsainadh/neovault/vault/internal/service"
	pb "github.com/iamsainadh/neovault/vault/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	slog.SetDefault(logger)

	cfg := config.LoadConfig()

	pool, err := database.NewConnection(cfg)

	if err != nil {
		logger.Error("failed to connect to database", slog.Any("error", err))
		os.Exit(1)
	}

	defer pool.Close()

	logger.Info("database connected")

	//! Repo

	repo := repository.NewVaultRepository(pool)

	svc := service.NewVaultService(repo)

	h := handler.NewVaultHandler(svc)

	grpcServer := grpc.NewServer()

	healthServer := health.NewServer()

	healthpb.RegisterHealthServer(grpcServer, healthServer)

	healthServer.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)

	pb.RegisterVaultServiceServer(grpcServer, h)

	addr := ":" + cfg.GRPCPort

	listener, err := net.Listen("tcp", addr)

	if err != nil {
		logger.Error("failed to listen", slog.String("addr", addr), slog.Any("error", err))
		os.Exit(1)
	}

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			logger.Error("gRPC server error", slog.Any("error", err))
			os.Exit(1)
		}
	}()

	<-quit

	logger.Info("shutdown signal received — draining in-flight RPCs")
	healthServer.SetServingStatus("", healthpb.HealthCheckResponse_NOT_SERVING)
	grpcServer.GracefulStop()
	logger.Info("server stopped cleanly")
}

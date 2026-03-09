package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/iamsainadh/neovault/auth/internal/app"
	"github.com/iamsainadh/neovault/auth/internal/config"
	"github.com/iamsainadh/neovault/auth/internal/database"
	"github.com/iamsainadh/neovault/auth/internal/server"
)

func main() {

	cfg := config.LoadConfig()

	pool, err := database.NewConnectionPool(cfg)

	if err != nil {
		fmt.Printf("failed to connect to database %v", err)
		os.Exit(1)
	}

	defer pool.Close()

	app := app.NewApplication(pool, cfg)

	grpcServer := server.NewGRPCServerConn(app.Handler)

	go func() {
		addr := ":" + cfg.PORT

		log.Printf("gRPC running on %s", addr)

		if err := grpcServer.Start(addr); err != nil {
			log.Fatalf("grpc server failed: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	log.Println("shutting down server...")
	grpcServer.Stop()
}

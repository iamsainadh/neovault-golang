package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/iamsainadh/neovault/scheduler/internal/config"
	"github.com/iamsainadh/neovault/scheduler/internal/database"
	"github.com/iamsainadh/neovault/scheduler/internal/jobs"
	"github.com/iamsainadh/neovault/scheduler/internal/runner"
	sharedServiceClient "github.com/iamsainadh/neovault/shared/grpcclients/secrets"
	sharedVaultClient "github.com/iamsainadh/neovault/shared/grpcclients/vault"
)

func main() {

	cfg, err := config.Load()

	if err != nil {
		fmt.Printf("Error At :%s", err)
	}

	pool, err := database.NewPoolConnection(cfg)

	if err != nil {
		log.Fatalf("db connection failed: %v", err)
	}

	secretClient, err := sharedServiceClient.NewSecretsClient(cfg.SecretPort)

	if err != nil {
		log.Fatalf("db connection failed: %v", err)
	}

	secretVersion, err := sharedServiceClient.NewSecretVersionClient(cfg.SecretPort)

	if err != nil {
		log.Fatalf("db connection failed: %v", err)
	}

	vaultClient, err := sharedVaultClient.NewVaultClient(cfg.VaultPort)

	if err != nil {
		log.Fatalf("vault client init failed: %v", err)
	}

	job := jobs.NewRotationClient(secretClient, vaultClient, secretVersion, pool)

	runner := runner.NewRunner(job)

	runner.Register("*/5 * * * *")

	runner.Start()

	log.Println("Scheduler started...")

	sig := make(chan os.Signal, 1)

	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	<-sig

	log.Println("Shutting down...")

}

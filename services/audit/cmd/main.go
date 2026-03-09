package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/iamsainadh/neovault/audit/internal/config"
	"github.com/iamsainadh/neovault/audit/internal/database"
	"github.com/iamsainadh/neovault/audit/internal/repository"
	"github.com/iamsainadh/neovault/audit/internal/service"
	"github.com/twmb/franz-go/pkg/kgo"
)

func main() {
	conf := config.LoadConfig()

	ctx := context.Background()

	pool, err := database.NewPostgresConnection(ctx, conf)

	client, err := kgo.NewClient(
		kgo.SeedBrokers(
			conf.KafkaBrokerUri,
		),
		kgo.ConsumeTopics(
			"audit.logs",
		),
		kgo.ConsumerGroup(
			"audit-service",
		),
		kgo.ConsumeStartOffset(
			kgo.NewOffset().AtEnd(),
		),
	)

	if err != nil {
		log.Fatal(err)
	}

	repo := repository.NewAuditRepository(pool)

	service.NewConsumer(client, repo)
	r := gin.Default()

	r.Run(conf.Port)
}

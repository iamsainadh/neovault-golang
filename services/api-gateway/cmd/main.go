package main

import (
	"log"

	"github.com/bytedance/gopkg/util/logger"
	"github.com/gin-gonic/gin"
	"github.com/iamsainadh/neovault/api-gateway/internal/config"
	"github.com/iamsainadh/neovault/api-gateway/internal/database"
	"github.com/iamsainadh/neovault/api-gateway/internal/middleware"
	"github.com/iamsainadh/neovault/api-gateway/internal/repository"
	"github.com/iamsainadh/neovault/api-gateway/internal/router"
	"github.com/iamsainadh/neovault/shared/kafka"
	"github.com/iamsainadh/neovault/shared/metrics"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/twmb/franz-go/pkg/kgo"
)

func main() {
	metrics.RegisterCollectors()
	cfg, err := config.LoadConfig()

	pool, err := database.NewConnectionPool(cfg)

	if err != nil {
		logger.Info(err)
		return
	}

	repo := repository.NewRepository(pool)

	client, err := kgo.NewClient(
		kgo.SeedBrokers(cfg.KAFKA_URI),
	)

	if err != nil {
		log.Fatal(err)
	}

	kafka.NewProducer(client)

	defer client.Close()

	ginRouter := gin.Default()

	ginRouter.GET("/metrics", gin.WrapH(promhttp.Handler()))

	ginRouter.Use(middleware.MetricMiddleware())

	router.NewRouter(ginRouter, cfg.JwtSecret, cfg, repo)

	if err := ginRouter.Run(cfg.PORT); err != nil {
		log.Fatal(err)
	}

}

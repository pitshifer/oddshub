package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/pitshifer/oddshub/internal/application"
	"github.com/pitshifer/oddshub/internal/cache"
	"github.com/pitshifer/oddshub/internal/collector/theoddsapi"
	"github.com/pitshifer/oddshub/internal/config"
	"github.com/pitshifer/oddshub/internal/storage/postgres"
	"github.com/pitshifer/oddshub/internal/transport/handler"
)

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatalln(err)
	}

	ctx := context.Background()

	var logLevel slog.Level
	if err := logLevel.UnmarshalText([]byte(config.LogLevel)); err != nil {
		logLevel = slog.LevelInfo
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	}))
	slog.SetDefault(logger)

	storage, err := postgres.New(ctx, config.DatabaseURL, logger, config.OddsWay)
	if err != nil {
		logger.Error("failed to initialize storage", "error", err)
		os.Exit(1)
	}
	defer storage.Close()

	cache := cache.New(storage)
	if err := cache.Warm(ctx, []string{"soccer_epl"}); err != nil {
		logger.Error("failed to warm cache", "error", err)
		os.Exit(1)
	}

	client := theoddsapi.NewClient(config.TheOddsApiKey)

	oddsService := application.NewOddsService(cache, client, logger)

	httpHandler := handler.New(cache, client, oddsService, logger)
	router := handler.NewRouter(httpHandler)
	if err = http.ListenAndServe(":8080", router); err != nil {
		logger.Error("failed to start server", "error", err)
		os.Exit(1)
	}
}

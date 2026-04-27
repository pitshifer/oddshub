package main

import (
	"context"
	"log"
	"net/http"

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

	storage, err := postgres.New(ctx, config.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer storage.Close()

	cache := cache.New(storage)
	if err := cache.Warm(ctx, []string{"soccer_epl"}); err != nil {
		log.Fatal(err)
	}

	client := theoddsapi.NewClient(config.TheOddsApiKey)

	httpHandler := handler.New(cache, client)
	router := handler.NewRouter(httpHandler)
	if err = http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}

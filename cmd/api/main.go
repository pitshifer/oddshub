package main

import (
	"context"
	"log"
	"net/http"

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

	// client := theoddsapi.NewClient(config.TheOddsApiKey)
	// odds, err := client.GetOdds(ctx, "soccer_epl")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// err = storage.SaveOdds(ctx, "theoddsapi", odds)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	httpHandler := handler.New(storage)
	router := handler.NewRouter(httpHandler)
	if err = http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}

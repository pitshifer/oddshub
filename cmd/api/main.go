package main

import (
	"context"
	"fmt"
	"log"

	"github.com/pitshifer/oddshub/internal/collector/theoddsapi"
	"github.com/pitshifer/oddshub/internal/config"
)

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatalln(err)
	}

	client := theoddsapi.NewClient(config.TheOddsApiKey)
	odds, err := client.GetOdds(context.Background(), "soccer_epl")
	if err != nil {
		log.Fatal(err)
	}

	for _, o := range odds {
		fmt.Println(o.HomeTeam, "vs", o.AwayTeam)
	}
}

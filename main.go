package main

import (
	"context"
	"log"
	"os"

	"github.com/raylas/query-bot/pkg/bot"
	"github.com/raylas/query-bot/pkg/config"
)

func main() {
	log := log.New(os.Stdout, "query-bot: ", log.Lshortfile|log.LstdFlags)
	ctx := context.Background()

	// Load configuration
	config, err := config.Load()
	if err != nil {
		log.Fatalf("[ERROR] %s \n", err)
	}

	// Initialize Slack bot
	s, err := bot.New()
	if err != nil {
		log.Fatalf("[ERROR] %s \n", err)
	}
	s.Logger = log

	// Start Slack bot
	if err := s.Listen(ctx, config); err != nil {
		s.Logger.Fatalf("[ERROR] %s \n", err)
	}
}

package main

import (
	"log"
	_ "github.com/joho/godotenv/autoload"
	"os"
)

// Config Stores all the secrets necessary for the bot to work
type Config struct {
	Token          string
}

func loadConfig() *Config {
	// Receive the required token
	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		log.Fatalf("BOT_TOKEN env variable not set")
	}
	return &Config{
		Token: token,
	}
}
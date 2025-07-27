package main

import (
	"fmt"
	_ "github.com/joho/godotenv/autoload"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Token          string
}

func loadConfig() (*Config, error) {
	// Receive the required token
	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("BOT_TOKEN env variable not set")
	}

	return &Config{
		Token: token,
	}, nil
}
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
	WhiteListUsers []int32
	AiApiKey       string
}

func loadConfig() (*Config, error) {
	// Receive the required token
	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		return nil, fmt.Errorf("BOT_TOKEN env variable not set")
	}

	// Parsing the list of users
	usersStr := os.Getenv("USERS")
	users, err := parseNumbers(usersStr)
	if err != nil {
		return nil, fmt.Errorf("Parsing error USERS: %v", err)
	}

	return &Config{
		Token:          token,
		WhiteListUsers: users,
	}, nil
}

func parseNumbers(input string) ([]int32, error) {
	if input == "" {
		return []int32{}, nil
	}

	var nums []int32
	parts := strings.SplitSeq(input, ",")

	for part := range parts {
		numStr := strings.TrimSpace(part)
		if numStr == "" {
			continue
		}

		num, err := strconv.ParseInt(numStr, 10, 32)
		if err != nil {
			return nil, fmt.Errorf("некорректное число '%s': %v", numStr, err)
		}
		nums = append(nums, int32(num))
	}

	return nums, nil
}

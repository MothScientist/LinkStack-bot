package main

import (
	"encoding/json"
	"log"
	"os"
)

var jsonHelpMsg map[string]string

func loadLocaleJson() {
	data, err := os.ReadFile("help_msg.json")
	if err != nil {
		log.Fatalf("Error reading .json file: %v;", err)
	}
	if err = json.Unmarshal(data, &jsonHelpMsg); err != nil {
		log.Fatalf("Error loading .json data into memory: %v;", err)
	}
}

func getLocaleHelpMsg(lang string) string {
	switch lang {
	case "en", "ru", "es":
	default:
		lang = "en"
	}
	return jsonHelpMsg[lang]
}
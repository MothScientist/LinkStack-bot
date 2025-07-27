package main

import (
	"context"
	"log"
	"sync"

	"github.com/go-telegram/bot"
)

// The following variables are required to store the calculated value within a single request (between Match and Handler), avoiding data races
var urlCacheLink sync.Map   // Saves link
var urlCacheLinkId sync.Map // Saves link id

// Function to launch the bot
func botProcess(token string) {
	opts := []bot.Option{
		bot.WithDefaultHandler(baseHandler),
	}

	b, err := bot.New(token, opts...)

	if err != nil {
		log.Fatalf("bot not started: %v", err)
	}

	b.RegisterHandlerMatchFunc(addMatch, addHandler)
	b.RegisterHandlerMatchFunc(getMatch, getHandler)
	b.RegisterHandlerMatchFunc(delMatch, delHandler)

	b.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, startHandler)
	//b.RegisterHandler(bot.HandlerTypeMessageText, "/help", bot.MatchTypeExact, helpHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/list", bot.MatchTypeExact, listHandler)
	b.RegisterHandler(bot.HandlerTypeMessageText, "/rdm", bot.MatchTypeExact, rdmHandler)

	b.Start(context.TODO())
}

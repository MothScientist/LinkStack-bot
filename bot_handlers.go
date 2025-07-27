package main

import (
	"context"
	"fmt"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

// Handler for adding a new link to the repository
func addHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	urlText, ok := urlCacheLink.Load(getCompositeSyncMapKey(update))
	defer urlCacheLink.Delete(getCompositeSyncMapKey(update)) // Remove link from global cache

	if !ok { return } // TODO - logs

	dbData := DbData{
		TelegramId: update.Message.From.ID,
		Url:        urlText.(string),
	}

	urlNumber, status, err := recordIsExists(&dbData)

	if !status && err == nil {
		// If the record does not exist yet, we get the title and write it down
		dbData.Title = getFirstH1Text(dbData.Url)
		urlNumber, err = addToStorage(&dbData)
	}

	var outputText string
	if err != nil {
		fmt.Println(err)
		outputText = "Error writing link to storage"
	} else if status {
		outputText = fmt.Sprintf("This link already exists in the repository\nNumber: %v", urlNumber)
	} else {
		outputText = fmt.Sprintf("Added to the list\nNumber: %v", urlNumber)
	}
	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      outputText,
		ParseMode: models.ParseModeMarkdown,
	})

	if err != nil { return } // TODO - logs
}

// Handler for getting a link by its number from the storage
func getHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	linkId, ok := urlCacheLinkId.Load(getCompositeSyncMapKey(update))
	defer urlCacheLinkId.Delete(getCompositeSyncMapKey(update)) // Remove link id from global cache
	if !ok { return } // TODO - logs

	dbData := DbData{
		TelegramId: update.Message.From.ID,
		LinkId:     linkId.(int32),
	}
	urlText, title, status, err := getFromStorage(&dbData)

	var outputText string
	if err != nil {
		outputText = "Error getting record from storage"
	} else if urlText == "" {
		outputText = "There is no link with this number in the repository"
	} else if !status && urlText != "" {
		outputText = "The entry with this number has been deleted"
	} else {
		outputText = fmt.Sprintf("<a href=\"%s\">%s</a>", urlText, title)
	}

	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      outputText,
		ParseMode: models.ParseModeHTML,
	})
	if err != nil { return } // TODO - logs
}

// Handler for deleting a link by its number from storage
func delHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	linkId, ok := urlCacheLinkId.Load(getCompositeSyncMapKey(update))
	defer urlCacheLinkId.Delete(getCompositeSyncMapKey(update)) // Remove link from global cache
	if !ok { return } // TODO - logs

	dbData := DbData{
		TelegramId: update.Message.From.ID,
		LinkId:     linkId.(int32),
	}
	status, err := delFromStorage(&dbData)

	var outputText string
	if err != nil {
		fmt.Println(err)
		outputText = "Error deleting record from storage"
	} else if !status {
		outputText = "There is no link with this number in the repository"
	} else {
		outputText = "Links removed from storage"
	}

	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      outputText,
		ParseMode: models.ParseModeMarkdown,
	})
	if err != nil { return } // TODO - logs
}

// Handler for getting a list of links from storage
func listHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	dbData := DbData{
		TelegramId: update.Message.From.ID,
	}

	urls, err := getListFromStorage(&dbData)

	var outputText string
	if err != nil {
		outputText = "Error retrieving records from storage"
		fmt.Print(err)
	} else if urls == nil {
		outputText = "There are no active records in the repository."
	} else {
		outputText = getListMsg(urls)
	}

	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      outputText,
		ParseMode: models.ParseModeHTML,
	})
	if err != nil {
		return // TODO - logs
	}
}

// Handler for getting a random link from storage (/rdm command)
func rdmHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	dbData := DbData{
		TelegramId: update.Message.From.ID,
	}
	linkId, url, title, err := getRandomFromStorage(&dbData)

	var outputText string
	if err != nil {
		fmt.Println(err)
		outputText = "Error getting random record from storage"
	} else if linkId == 0 {
		outputText = "There are no active records in the repository"
	} else {
		outputText = fmt.Sprintf("%d: <a href=\"%s\">%s</a>", linkId, url, title)
	}

	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      outputText,
		ParseMode: models.ParseModeHTML,
	})
	if err != nil { return } // TODO - logs
}

// Handler for the /start command
func startHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      "Caught */start*",
		ParseMode: models.ParseModeMarkdown,
	})
	if err != nil {
		return // TODO - logs
	}
}

// Handler for commands not handled by other functions
func baseHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	//msgText := botMessageURL{
	//	Text:               update.Message.Text,
	//	ForwardMessageText: update.Message.CaptionEntities,
	//}
	fmt.Printf("%+v\n", update.Message.Entities)
	//fmt.Printf("%+v\n", update.Message.CaptionEntities)

	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      "Caught msg",
		ParseMode: models.ParseModeMarkdown,
	})
	if err != nil {
		return // TODO - logs
	}
}

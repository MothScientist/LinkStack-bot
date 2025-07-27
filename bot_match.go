package main

import (
	"fmt"
	"github.com/go-telegram/bot/models"
	"regexp"
	"strconv"
)

func addMatch(update *models.Update) bool {
	urlText := getFirstUrl(update.Message.Text, update.Message.Entities, update.Message.CaptionEntities)
	fmt.Println(urlText)
	if urlText != "" {
		urlCacheLink.Store(getCompositeSyncMapKey(update), urlText)
		return true
	}
	return false

}

func getMatch(update *models.Update) bool {
    return matchCommand(update, `^(?i)get\s+(\d+)$`)
}

func delMatch(update *models.Update) bool {
    return matchCommand(update, `^(?i)del\s+(\d+)$`)
}

// General function to check for commands like "get X" or "del X"
func matchCommand(update *models.Update, pattern string) bool {
    re := regexp.MustCompile(pattern)
    matches := re.FindStringSubmatch(update.Message.Text)

    if len(matches) != 2 { return false }

    num, err := strconv.Atoi(matches[1])
    if err != nil { return false }

    urlCacheLinkId.Store(getCompositeSyncMapKey(update), int32(num))
    return true
}
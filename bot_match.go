package main

import (
	"regexp"
	"strconv"

	"github.com/go-telegram/bot/models"
)

func addMatch(update *models.Update) bool {
	urlText := getFirstUrl(update.Message.Text, update.Message.Entities, update.Message.CaptionEntities)
	if urlText != "" {
		bridgeLink.Store(getCompositeSyncMapKey(update), urlText)
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

// matchCommand General function to check for commands like "get X" or "del X"
func matchCommand(update *models.Update, pattern string) bool {
    re := regexp.MustCompile(pattern)
    matches := re.FindStringSubmatch(update.Message.Text)

    if len(matches) != 2 { return false }

    num, err := strconv.Atoi(matches[1])
    if err != nil { return false }

	bridgeLinkId.Store(getCompositeSyncMapKey(update), int32(num))
    return true
}

// getFirstUrl Gets the first link from the chain: message -> formatted message -> forwarded messages
func getFirstUrl(urlMsgText string, urlEntitiesText ...[]models.MessageEntity) string {
	if urlText := getUrlFromMessage(urlMsgText); urlText != "" {
		return urlText
	}
	for _, urlEntText := range urlEntitiesText {
		if urlText := getUrlFromEntityMsg(urlEntText); urlText != "" {
			return urlText
		}
	}
	return ""
}

// getUrlFromMessage Extracts a reference from a string
func getUrlFromMessage(messageText string) string {
	match := regexpUrl(messageText, false)
	if match != "" && isUrl(match) {
		return match
	}
	return ""
}

// getUrlFromEntityMsg Finds and returns a link from forwarded messages or rich text
func getUrlFromEntityMsg(entityMsg []models.MessageEntity) string {
	for _, msg := range entityMsg {
		if msg.URL != "" {
			return msg.URL
		}
	}
	return ""
}
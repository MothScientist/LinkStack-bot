package main

import (
	"github.com/go-telegram/bot/models"
)

// DbData Data passed between bot handlers and storage
type DbData struct {
	TelegramId int64
	LinkId     int32
	Status     bool
	Url        string
	Title      string
}

// SqlQueries When launched, it loads .sql queries into RAM
type SqlQueries struct {
	addRecord        string
	delRecord        string
	getRecord        string
	getListRecords   string
	getRandomRecord  string
	recordIsExists   string
}

// Link Structure for storing a link with a title
type Link struct {
	URL   string
	Title string
}

// CompositeSyncMapKey Composite key structure for sync.Map, since message id is unique within one dialog
type CompositeSyncMapKey struct {
	TelegramId int64
	MsgId      int
}

// getCompositeSyncMapKey Getting a composite key for sync.Map
func getCompositeSyncMapKey(update *models.Update) CompositeSyncMapKey {
	return CompositeSyncMapKey{
		TelegramId: update.Message.From.ID,
		MsgId:      update.Message.ID,
	}
}
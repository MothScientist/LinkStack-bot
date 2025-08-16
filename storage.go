package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func getSql(filename string) string {
	var res string
	switch filename {
	case "add_record":
		res = queries.addRecord
	case "del_record":
		res = queries.delRecord
	case "get_record":
		res = queries.getRecord
	case "get_list_records":
		res = queries.getListRecords
	case "get_random_record":
		res = queries.getRandomRecord
	case "record_is_exists":
		res = queries.recordIsExists
	}
	return res
}

// addToStorage Function of recording a link in the database
func addToStorage(dbData DbData) (linkId int32, err error) {
	db, err := sql.Open("sqlite3", "main.sqlite3")
	if err != nil {
		return 0, fmt.Errorf("Error opening connection: %w", err)
	}
	defer db.Close()

	err = db.QueryRow(
		getSql("add_record"),
		dbData.TelegramId,
		dbData.TelegramId,
		dbData.Url,
		dbData.Title,
	).Scan(&linkId)
	dbData.LinkId = linkId

	if err != nil {
		return 0, fmt.Errorf("Error retrieving data: %w", err)
	}

	cacheData := Link{
		URL:   dbData.Url,
		Title: dbData.Title,
	}
	getUserCache.Add(getCacheCompositeKeyByDbData(dbData), cacheData)
	log.Print("[DB] addToStorage")
	return linkId, nil
}

// getFromStorage Function for getting a link by id from the database
func getFromStorage(dbData DbData) (url string, title string, status bool, err error) {
	cacheLink := getUserCache.Get(getCacheCompositeKeyByDbData(dbData))
	if cacheLink != nil {
		return cacheLink.URL, cacheLink.Title, true, nil
	}
	db, err := sql.Open("sqlite3", "main.sqlite3")
	if err != nil {
		return "", "", false, fmt.Errorf("Error opening connection: %w", err)
	}
	defer db.Close()

	err = db.QueryRow(
		getSql("get_record"),
		dbData.TelegramId,
		dbData.LinkId,
	).Scan(&url, &title, &status)

	if err != nil {
		if err == sql.ErrNoRows {
			return "", "", false, nil
		}
		return "", "", false, fmt.Errorf("Error retrieving data: %w", err)
	}

	cacheData := Link{
		URL:   url,
		Title: title,
	}
	if status == true && url != "" {
		getUserCache.Add(getCacheCompositeKeyByDbData(dbData), cacheData)
	}
	log.Print("[DB] getFromStorage")
	return url, title, status, nil
}

// getRandomFromStorage Function to get a random active user record from the database
func getRandomFromStorage(dbData DbData) (linkId int32, url string, title string, err error) {
	db, err := sql.Open("sqlite3", "main.sqlite3")
	if err != nil {
		return 0, "", "", fmt.Errorf("Error opening connection: %w", err)
	}
	defer db.Close()

	err = db.QueryRow(
		getSql("get_random_record"),
		dbData.TelegramId,
	).Scan(&linkId, &url, &title)

	if err != nil {
		if err == sql.ErrNoRows {
			return 0, "", "", nil
		}
		return 0, "", "", fmt.Errorf("Error retrieving data: %w", err)
	}

	cacheData := Link{
		URL:   url,
		Title: title,
	}
	if url != "" {
		getUserCache.Add(getCacheCompositeKeyByDbData(dbData), cacheData)
	}
	log.Print("[DB] getRandomFromStorage")
	return linkId, url, title, nil
}

// delFromStorage Function for deleting a link by id from the database
func delFromStorage(dbData DbData) (bool, error) {
	db, err := sql.Open("sqlite3", "main.sqlite3")
	if err != nil {
		return false, fmt.Errorf("Error opening connection: %w", err)
	}
	defer db.Close()

	res, err := db.Exec(
		getSql("del_record"),
		dbData.TelegramId,
		dbData.LinkId,
	)
	if err != nil {
		return false, fmt.Errorf("Error executing request: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("Error getting number of deleted rows: %w", err)
	} else if rowsAffected == 0 {
		return false, nil // The record was previously deleted
	}

	getUserCache.Del(getCacheCompositeKeyByDbData(dbData))
	log.Print("[DB] delFromStorage")
	return true, nil
}

// getListFromStorage Function to get a list of valid links from the repository {id: [title, url]}
func getListFromStorage(dbData DbData) (urls map[int32]Link, err error) {
	db, err := sql.Open("sqlite3", "main.sqlite3")
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query(
		getSql("get_list_records"),
		dbData.TelegramId,
	)

	urls = make(map[int32]Link)

	if err != nil {
		if err == sql.ErrNoRows {
			return urls, nil
		}
		return nil, fmt.Errorf("Error retrieving data: %w", err)
	}

	for rows.Next() {
		var (
			id    int32
			url   string
			title string
		)

		if err := rows.Scan(&id, &url, &title); err != nil {
			return nil, fmt.Errorf("Error reading data: %w", err)
		}

		urls[id] = Link{
			URL:   url,
			Title: title,
		}
	}

	// Checking that the iteration has completed correctly
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Error processing results: %w", err)
	}

	if err := rows.Close(); err != nil {
		return nil, fmt.Errorf("Error closing connection for rows: %w", err)
	}

	log.Print("[DB] getListFromStorage")
	return urls, nil
}

func recordIsExists(dbData *DbData) (linkId int32, status bool, err error) {
	db, err := sql.Open("sqlite3", "main.sqlite3")
	if err != nil {
		return 0, false, err
	}
	defer db.Close()

	err = db.QueryRow(
		getSql("record_is_exists"),
		dbData.TelegramId,
		dbData.Url,
	).Scan(&linkId)

	if err != nil {
		if err == sql.ErrNoRows {
			return 0, false, nil
		}
		return 0, false, fmt.Errorf("Error retrieving data: %w", err)
	}

	return linkId, linkId != 0, nil
}

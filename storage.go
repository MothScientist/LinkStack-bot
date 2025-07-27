package main

import (
	"database/sql"
	"fmt"

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

// Function of recording a link in the database
func addToStorage(dbData *DbData) (linkId int32, err error) {
	db, err := sql.Open("sqlite3", "main.sqlite3")
	if err != nil {
		return 0, fmt.Errorf("Ошибка открытия соединения: %w", err)
	}
	defer db.Close()

	err = db.QueryRow(
		getSql("add_record"),
		dbData.TelegramId,
		dbData.TelegramId,
		dbData.Url,
		dbData.Title,
	).Scan(&linkId)

	if err != nil {
		return 0, fmt.Errorf("Ошибка получения данных: %w", err)
	}

	return linkId, nil
}

// Function for getting a link by id from the database
func getFromStorage(dbData *DbData) (url string, title string, status bool, err error) {
	db, err := sql.Open("sqlite3", "main.sqlite3")
	if err != nil {
		return "", "", false, fmt.Errorf("Ошибка открытия соединения: %w", err)
	}
	defer db.Close()

	err = db.QueryRow(
		getSql("get_record"),
		dbData.TelegramId,
		dbData.LinkId,
	).Scan(&url, &title, &status)

	if err != nil {
		return "", "", false, fmt.Errorf("Ошибка получения данных: %w", err)
	}

	return url, title, status, nil
}

// Function to get a random active user record from the database
func getRandomFromStorage(dbData *DbData) (linkId int32, url string, title string, err error) {
	db, err := sql.Open("sqlite3", "main.sqlite3")
	if err != nil {
		return 0, "", "", fmt.Errorf("Ошибка открытия соединения: %w", err)
	}
	defer db.Close()

	err = db.QueryRow(
		getSql("get_random_record"),
		dbData.TelegramId,
	).Scan(&linkId, &url, &title)

	if err != nil {
		return 0, "", "", fmt.Errorf("Ошибка получения данных: %w", err)
	}

	return linkId, url, title, nil
}

// Function for deleting a link by id from the database
func delFromStorage(dbData *DbData) (bool, error) {
	db, err := sql.Open("sqlite3", "main.sqlite3")
	if err != nil {
		return false, fmt.Errorf("Ошибка открытия соединения: %w", err)
	}
	defer db.Close()

	res, err := db.Exec(
		getSql("del_record"),
		dbData.TelegramId,
		dbData.LinkId,
	)
	if err != nil {
		return false, fmt.Errorf("Ошибка выполнения запрос: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("Ошибка при получении количества удаленных строк: %w", err)
	} else if rowsAffected == 0 {
		return false, fmt.Errorf("Запись не удалена: %w", err)
	}

	return true, nil
}

// Function to get a list of valid links from the repository {id: [title, url]}
func getListFromStorage(dbData *DbData) (urls map[int32]Link, err error) {
	db, err := sql.Open("sqlite3", "main.sqlite3")
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query(
		getSql("get_list_records"),
		dbData.TelegramId,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("Ошибка получения данных: %w", err)
	}

	urls = make(map[int32]Link)
	for rows.Next() {
		var (
			id    int32
			url   string
			title string
		)

		if err := rows.Scan(&id, &url, &title); err != nil {
			return nil, fmt.Errorf("Ошибка чтения данных: %w", err)
		}

		urls[id] = Link{
			URL:   url,
			Title: title,
		}
	}

	// Checking that the iteration has completed correctly
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Ошибка при обработке результатов: %w", err)
	}

	err = rows.Close()
	if err != nil {
		return nil, fmt.Errorf("Ошибка закрытия соединения для rows: %w", err)
	}

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
		return 0, false, fmt.Errorf("Ошибка получения данных: %w", err)
	}

	return linkId, linkId != 0, nil
}

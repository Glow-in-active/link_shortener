package storage

import (
	"database/sql"
	"fmt"
	"log"
)

func AddData(db *sql.DB, shortUrl, longUrl string) {
	_, err := db.Exec("INSERT INTO data (short_url, long_url) VALUES ($1, $2)", shortUrl, longUrl)
	if err != nil {
		log.Fatal("Ошибка при добавлении данных: ", err)
	}
	fmt.Println("Данные успешно добавлены!")
}

func GetData(db *sql.DB) ([]string, error) {
	rows, err := db.Query("SELECT id, short_url, long_url FROM data")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []string
	for rows.Next() {
		var id string
		var shortUrl, longUrl string
		if err := rows.Scan(&id, &shortUrl, &longUrl); err != nil {
			return nil, err
		}
		result = append(result, fmt.Sprintf("ID: %s, Short URL: %s, Long URL: %s", id, shortUrl, longUrl))
	}
	return result, nil
}

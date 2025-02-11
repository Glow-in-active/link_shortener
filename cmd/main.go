package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	// Подключаемся к базе данных
	connStr := "postgres://admin:secret@localhost:5432/mydb?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Ошибка подключения к базе данных: ", err)
	}
	defer db.Close()

	addData(db, "short.ly/abc123", "https://example.com/some-long-url")

	data, err := getData(db)
	if err != nil {
		log.Fatal("Ошибка при извлечении данных: ", err)
	}
	fmt.Println("Данные из таблицы:", data)
}

func addData(db *sql.DB, shortUrl, longUrl string) {
	_, err := db.Exec("INSERT INTO data (short_url, long_url) VALUES ($1, $2)", shortUrl, longUrl)
	if err != nil {
		log.Fatal("Ошибка при добавлении данных: ", err)
	}
	fmt.Println("Данные успешно добавлены!")
}

func getData(db *sql.DB) ([]string, error) {
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

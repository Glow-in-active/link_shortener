package storage

import (
	"database/sql"
	"fmt"
)

func AddData(db *sql.DB, shortUrl, longUrl string) error {
	_, err := db.Exec("INSERT INTO data (short_url, long_url) VALUES ($1, $2)", shortUrl, longUrl)
	if err != nil {
		return fmt.Errorf("ошибка при добавлении данных: %w", err)
	}
	return nil
}

func GetData(db *sql.DB, shortUrl string) ([]string, error) {
	rows, err := db.Query("SELECT long_url FROM data WHERE short_url = $1", shortUrl)
	if err != nil {
		return nil, fmt.Errorf("ошибка при запросе данных: %w", err)
	}
	defer rows.Close()

	var result []string
	for rows.Next() {
		var longUrl string
		if err := rows.Scan(&longUrl); err != nil {
			return nil, fmt.Errorf("ошибка при чтении данных: %w", err)
		}
		result = append(result, longUrl)
	}
	return result, nil
}

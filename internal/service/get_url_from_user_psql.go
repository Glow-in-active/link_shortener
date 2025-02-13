package service

import (
	"database/sql"
	"fmt"
	"link_shortener/internal/generator"
	"link_shortener/internal/storage"
)

func SaveURL(db *sql.DB, longURL string) (string, error) {
	shortURL := generator.GenerateShortURL(longURL)
	storage.AddData(db, shortURL, longURL)

	return shortURL, nil
}

func GetLongURL(db *sql.DB, shortURL string) (string, error) {
	results, err := storage.GetData(db, shortURL)
	if err != nil {
		return "", fmt.Errorf("ошибка получения данных из БД: %w", err)
	}
	if len(results) == 0 {
		return "", fmt.Errorf("URL не найден")
	}
	return results[0], nil
}

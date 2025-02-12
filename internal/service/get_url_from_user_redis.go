package service

import (
	"fmt"
	"link_shortener/internal/generator"
	"link_shortener/internal/storage"

	"github.com/redis/go-redis/v9"
)

func SaveURLRedis(rdb *redis.Client, longURL string) (string, error) {
	shortURL := generator.GenerateShortURL(longURL)
	err := storage.RedisAddData(rdb, shortURL, longURL)
	if err != nil {
		return "", fmt.Errorf("ошибка сохранения в Redis: %w", err)
	}
	return shortURL, nil
}

func GetLongURLRedis(rdb *redis.Client, shortURL string) (string, error) {
	longURL, err := storage.RedisGetData(rdb, shortURL)
	if err != nil {
		return "", fmt.Errorf("ошибка получения данных из Redis: %w", err)
	}
	return longURL, nil
}

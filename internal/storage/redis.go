package storage

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func RedisAddData(rdb *redis.Client, shortUrl, longUrl string) error {
	err := rdb.Set(context.Background(), shortUrl, longUrl, 0).Err()
	if err != nil {
		return fmt.Errorf("ошибка при добавлении данных в Redis: %w", err)
	}
	return nil
}

func RedisGetData(rdb *redis.Client, shortUrl string) (string, error) {
	longUrl, err := rdb.Get(context.Background(), shortUrl).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("данные не найдены")
	} else if err != nil {
		return "", fmt.Errorf("ошибка при запросе данных из Redis: %w", err)
	}
	return longUrl, nil
}

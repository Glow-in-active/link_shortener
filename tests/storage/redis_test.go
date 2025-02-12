package storage

import (
	"github.com/alicebob/miniredis/v2"
	"link_shortener/internal/storage"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

var rdb *redis.Client
var mockRedis *miniredis.Miniredis

func setupMockRedis(t *testing.T) {
	var err error
	mockRedis, err = miniredis.Run()
	assert.NoError(t, err, "Ошибка при запуске мок-Redis")

	rdb = redis.NewClient(&redis.Options{
		Addr: mockRedis.Addr(),
	})
}

func TestAddAndGetData(t *testing.T) {
	setupMockRedis(t)

	shortKey := "test123"
	longURL := "https://example.com/test"

	err := storage.RedisAddData(rdb, shortKey, longURL)
	assert.NoError(t, err, "Ошибка при добавлении данных")

	result, err := storage.RedisGetData(rdb, shortKey)
	assert.NoError(t, err, "Ошибка при получении данных")
	assert.Equal(t, longURL, result, "Полученный URL не совпадает с ожидаемым")
}

func TestGetNonExistingData(t *testing.T) {
	setupMockRedis(t)

	_, err := storage.RedisGetData(rdb, "not_exist")
	assert.Error(t, err, "Ожидалась ошибка для несуществующего ключа")
}

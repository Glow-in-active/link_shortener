package service

import (
	"github.com/redis/go-redis/v9"
	"link_shortener/internal/service"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
)

func setupRedisTest() (*redis.Client, *miniredis.Miniredis) {
	mr, err := miniredis.Run()
	if err != nil {
		panic(err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})
	return rdb, mr
}

func TestSaveURLRedis(t *testing.T) {
	rdb, mr := setupRedisTest()
	defer mr.Close()

	longURL := "https://example.com"
	shortURL, err := service.SaveURLRedis(rdb, longURL)

	assert.NoError(t, err, "ошибка при сохранении в Redis")
	assert.NotEmpty(t, shortURL, "короткий URL не должен быть пустым")

	storedURL, err := mr.Get(shortURL)
	assert.NoError(t, err)
	assert.Equal(t, longURL, storedURL, "сохранённый URL должен совпадать")
}

func TestGetLongURLRedis(t *testing.T) {
	rdb, mr := setupRedisTest()
	defer mr.Close()

	shortURL := "abc123"
	longURL := "https://example.com"

	mr.Set(shortURL, longURL)

	retrievedURL, err := service.GetLongURLRedis(rdb, shortURL)
	assert.NoError(t, err, "ошибка при получении из Redis")
	assert.Equal(t, longURL, retrievedURL, "полученный URL должен совпадать с сохранённым")
}

func TestGetLongURLRedis_NotFound(t *testing.T) {
	rdb, mr := setupRedisTest()
	defer mr.Close()

	shortURL := "notfound"

	retrievedURL, err := service.GetLongURLRedis(rdb, shortURL)
	assert.Error(t, err, "ожидалась ошибка для несуществующего ключа")
	assert.Empty(t, retrievedURL, "возвращённый URL должен быть пустым")
}

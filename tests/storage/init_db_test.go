package storage

import (
	"context"
	"link_shortener/internal/storage"
	"testing"

	"github.com/alicebob/miniredis/v2"
	_ "github.com/lib/pq" // Нужен для sql.Open
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func setuppMockRedis(t *testing.T) (*redis.Client, *miniredis.Miniredis) {
	mockRedis, err := miniredis.Run()
	assert.NoError(t, err, "Ошибка при запуске mock-Redis")

	rdb := redis.NewClient(&redis.Options{
		Addr: mockRedis.Addr(),
	})

	return rdb, mockRedis
}

func TestInitDatabasePostgres(t *testing.T) {
	db, err := storage.InitDatabase("postgres")

	assert.NoError(t, err, "Ошибка при инициализации Postgres")
	assert.NotNil(t, db, "БД не должна быть nil")
	assert.NotNil(t, db.Postgres, "Postgres клиент должен быть инициализирован")
	assert.Nil(t, db.Redis, "Redis должен быть nil для Postgres")

	db.Close()
}

func TestInitDatabaseRedis(t *testing.T) {
	rdb, mockRedis := setuppMockRedis(t)
	defer mockRedis.Close()

	assert.NotNil(t, rdb, "Redis клиент должен быть инициализирован")

	ctx := context.Background()

	err := rdb.Set(ctx, "test_key", "test_value", 0).Err()
	assert.NoError(t, err, "Ошибка при записи в Redis")

	val, err := rdb.Get(ctx, "test_key").Result()
	assert.NoError(t, err, "Ошибка при чтении из Redis")
	assert.Equal(t, "test_value", val, "Значение в Redis не совпадает")
}

func TestInitDatabaseInvalid(t *testing.T) {
	db, err := storage.InitDatabase("unknown")

	assert.Error(t, err, "Должна быть ошибка при неизвестном типе БД")
	assert.Nil(t, db, "БД должна быть nil при ошибке")
}

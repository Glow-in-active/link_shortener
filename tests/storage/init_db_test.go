package storage

import (
	"context"
	"link_shortener/internal/storage"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

var testConfig = &storage.Config{
	DBType: "postgres",
	Postgres: storage.PostgresConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "admin",
		Password: "secret",
		DBName:   "test_db",
		SSLMode:  "disable",
	},
	Redis: storage.RedisConfig{
		Host:     "localhost",
		Port:     6379,
		Password: "",
		DB:       0,
	},
}

func TestInitDatabasePostgres(t *testing.T) {
	mockDB, mock, err := sqlmock.New()
	assert.NoError(t, err, "Ошибка создания mock PostgreSQL")
	defer mockDB.Close()

	db := &storage.Database{
		Postgres: mockDB,
		Redis:    nil,
		Type:     "postgres",
	}

	assert.NotNil(t, db.Postgres, "Postgres клиент должен быть инициализирован")
	assert.Nil(t, db.Redis, "Redis должен быть nil для Postgres")

	mock.ExpectClose()
	err = db.Close()
	assert.NoError(t, err, "Ошибка при закрытии PostgreSQL")
}

func setuppMockRedis(t *testing.T) (*redis.Client, *miniredis.Miniredis) {
	mockRedis, err := miniredis.Run()
	assert.NoError(t, err, "Ошибка запуска mock-Redis")

	rdb := redis.NewClient(&redis.Options{
		Addr: mockRedis.Addr(),
	})

	return rdb, mockRedis
}

func TestInitDatabaseRedis(t *testing.T) {
	_, mockRedis := setuppMockRedis(t)
	defer mockRedis.Close()

	testConfig.DBType = "redis"
	db, err := storage.InitDatabase("redis", testConfig)
	assert.NoError(t, err, "Ошибка при инициализации Redis")
	assert.Nil(t, db.Postgres, "Postgres должен быть nil для Redis")
	assert.NotNil(t, db.Redis, "Redis клиент должен быть инициализирован")

	ctx := context.Background()
	err = db.Redis.Set(ctx, "test_key", "test_value", 0).Err()
	assert.NoError(t, err, "Ошибка при записи в Redis")

	val, err := db.Redis.Get(ctx, "test_key").Result()
	assert.NoError(t, err, "Ошибка при чтении из Redis")
	assert.Equal(t, "test_value", val, "Значение в Redis не совпадает")

	db.Close()
}

func TestInitDatabaseInvalid(t *testing.T) {
	db, err := storage.InitDatabase("unknown", testConfig)

	assert.Error(t, err, "Должна быть ошибка при неизвестном типе БД")
	assert.Nil(t, db, "БД должна быть nil при ошибке")
}

func TestDatabaseClose(t *testing.T) {
	mockDB, _, err := sqlmock.New()
	assert.NoError(t, err, "Ошибка создания mock PostgreSQL")
	defer mockDB.Close()

	db := &storage.Database{
		Postgres: mockDB,
		Redis:    nil,
		Type:     "postgres",
	}

	err = db.Close()
	assert.NoError(t, err, "Ошибка при закрытии соединений")
}

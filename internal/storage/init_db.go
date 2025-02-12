package storage

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

type Database struct {
	Postgres *sql.DB
	Redis    *redis.Client
	Type     string
}

func InitDatabase(dbType string) (*Database, error) {
	db := &Database{Type: dbType}

	switch dbType {
	case "postgres":
		connStr := "postgres://admin:secret@localhost:5432/mydb?sslmode=disable"
		pgDB, err := sql.Open("postgres", connStr)
		if err != nil {
			return nil, fmt.Errorf("ошибка подключения к Postgres: %w", err)
		}
		db.Postgres = pgDB

	case "redis":
		db.Redis = redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
		})

	default:
		return nil, fmt.Errorf("некорректный тип БД: %s", dbType)
	}

	log.Println("Подключено к", dbType)
	return db, nil
}

func (db *Database) Close() {
	if db.Postgres != nil {
		db.Postgres.Close()
	}
}

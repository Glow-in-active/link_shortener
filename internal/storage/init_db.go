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

func InitDatabase(dbType string, config interface{}) (*Database, error) {
	if dbType == "" {
		return nil, fmt.Errorf("db_type не указан в конфиге")
	}

	cfg, ok := config.(*Config)
	if !ok {
		return nil, fmt.Errorf("неверный тип конфигурации")
	}

	db := &Database{Type: dbType}

	switch dbType {
	case "postgres":
		// Используем cfg.Postgres для доступа к конфигурации PostgreSQL
		connStr := fmt.Sprintf(
			"postgres://%s:%s@%s:%d/%s?sslmode=%s",
			cfg.Postgres.User, cfg.Postgres.Password, cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.DBName, cfg.Postgres.SSLMode,
		)

		pgDB, err := sql.Open("postgres", connStr)
		if err != nil {
			return nil, fmt.Errorf("ошибка подключения к Postgres: %w", err)
		}
		db.Postgres = pgDB

	case "redis":
		// Используем cfg.Redis для доступа к конфигурации Redis
		db.Redis = redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
			Password: cfg.Redis.Password,
			DB:       cfg.Redis.DB,
		})

	default:
		return nil, fmt.Errorf("некорректный тип БД: %s", dbType)
	}

	log.Println("Подключено к", dbType)
	return db, nil
}

func (db *Database) Close() error {
	var err error

	if db.Postgres != nil {
		if closeErr := db.Postgres.Close(); closeErr != nil {
			err = fmt.Errorf("ошибка закрытия Postgres: %w", closeErr)
		}
		log.Println("Соединение с Postgres закрыто")
	}

	if db.Redis != nil {
		if closeErr := db.Redis.Close(); closeErr != nil {
			if err != nil {
				err = fmt.Errorf("%v; ошибка закрытия Redis: %w", err, closeErr)
			} else {
				err = fmt.Errorf("ошибка закрытия Redis: %w", closeErr)
			}
		}
		log.Println("Соединение с Redis закрыто")
	}

	return err
}

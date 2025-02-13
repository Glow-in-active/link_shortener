package storage

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	DBType   string
	Postgres PostgresConfig
	Redis    RedisConfig
}

type PostgresConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

func LoadConfig(filename string) (*Config, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("ошибка открытия конфигурационного файла: %w", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	config := &Config{}
	if err := decoder.Decode(config); err != nil {
		return nil, fmt.Errorf("ошибка декодирования конфигурационного файла: %w", err)
	}

	return config, nil
}

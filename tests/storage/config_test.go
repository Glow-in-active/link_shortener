package storage_test

import (
	"link_shortener/internal/storage"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createTempConfig(t *testing.T, content string) string {
	tmpFile, err := os.CreateTemp("", "config_*.json")
	assert.NoError(t, err, "Ошибка при создании временного файла")

	_, err = tmpFile.Write([]byte(content))
	assert.NoError(t, err, "Ошибка при записи во временный файл")

	err = tmpFile.Close()
	assert.NoError(t, err, "Ошибка при закрытии временного файла")

	return tmpFile.Name()
}

func TestLoadConfig_Success(t *testing.T) {
	configContent := `{
  "DBType": "postgres",
  "Postgres": {
   "Host": "localhost",
   "Port": 5432,
   "User": "admin",
   "Password": "secret",
   "DBName": "test_db",
   "SSLMode": "disable"
  },
  "Redis": {
   "Host": "localhost",
   "Port": 6379,
   "Password": "",
   "DB": 0
  }
 }`

	tmpFile := createTempConfig(t, configContent)
	defer os.Remove(tmpFile)

	config, err := storage.LoadConfig(tmpFile)
	assert.NoError(t, err, "Ошибка при загрузке конфигурации")
	assert.NotNil(t, config, "Конфигурация не должна быть nil")

	assert.Equal(t, "postgres", config.DBType, "Некорректный тип БД")
	assert.Equal(t, "localhost", config.Postgres.Host, "Некорректный хост Postgres")
	assert.Equal(t, 5432, config.Postgres.Port, "Некорректный порт Postgres")
}

func TestLoadConfig_FileNotFound(t *testing.T) {
	_, err := storage.LoadConfig("non_existent.json")
	assert.Error(t, err, "Должна быть ошибка при отсутствии файла")
}

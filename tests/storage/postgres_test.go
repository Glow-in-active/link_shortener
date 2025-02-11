package storage

import (
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"link_shortener/internal/storage"
)

func TestAddData(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Ошибка при создании мок базы данных: %s", err)
	}
	defer db.Close()

	mock.ExpectExec("INSERT INTO data").
		WithArgs("short123", "https://example.com").
		WillReturnResult(sqlmock.NewResult(1, 1))

	storage.AddData(db, "short123", "https://example.com")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Не все ожидания выполнены: %s", err)
	}
}

func TestGetData(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Ошибка при создании мок базы данных: %s", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "short_url", "long_url"}).
		AddRow(1, "short123", "https://example.com").
		AddRow(2, "short456", "https://google.com")

	mock.ExpectQuery("SELECT id, short_url, long_url FROM data").
		WillReturnRows(rows)

	data, err := storage.GetData(db)
	if err != nil {
		t.Fatalf("Ошибка выполнения GetData: %s", err)
	}

	expected := []string{
		"ID: 1, Short URL: short123, Long URL: https://example.com",
		"ID: 2, Short URL: short456, Long URL: https://google.com",
	}

	if len(data) != len(expected) {
		t.Fatalf("Ожидалось %d записей, но получено %d", len(expected), len(data))
	}

	for i, v := range expected {
		if data[i] != v {
			t.Errorf("Ожидалось %s, но получено %s", v, data[i])
		}
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Не все ожидания выполнены: %s", err)
	}
}

func TestGetData_Error(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Ошибка при создании мок базы данных: %s", err)
	}
	defer db.Close()

	mock.ExpectQuery("SELECT id, short_url, long_url FROM data").
		WillReturnError(fmt.Errorf("ошибка базы данных"))

	_, err = storage.GetData(db)
	if err == nil {
		t.Errorf("Ожидалась ошибка базы данных, но её не произошло")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Не все ожидания выполнены: %s", err)
	}
}

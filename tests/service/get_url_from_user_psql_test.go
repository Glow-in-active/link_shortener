package service

import (
	"link_shortener/internal/service"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"link_shortener/internal/generator"
)

func TestSaveURL(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open mock database: %v", err)
	}
	defer db.Close()

	longURL := "https://example.com"
	shortURL := generator.GenerateShortURL(longURL)

	mock.ExpectExec("INSERT INTO data").WithArgs(shortURL, longURL).
		WillReturnResult(sqlmock.NewResult(1, 1))

	result, err := service.SaveURL(db, longURL)
	assert.NoError(t, err)
	assert.Equal(t, shortURL, result)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unmet expectations: %v", err)
	}
}

func TestGetLongURL(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open mock database: %v", err)
	}
	defer db.Close()

	shortURL := "short1"
	longURL := "https://example.com"

	rows := sqlmock.NewRows([]string{"long_url"}).AddRow(longURL)
	mock.ExpectQuery("SELECT long_url FROM data WHERE short_url = ?").WithArgs(shortURL).
		WillReturnRows(rows)

	result, err := service.GetLongURL(db, shortURL)
	assert.NoError(t, err)
	assert.Equal(t, longURL, result)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unmet expectations: %v", err)
	}
}

func TestGetLongURL_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open mock database: %v", err)
	}
	defer db.Close()

	shortURL := "nonexistent"

	mock.ExpectQuery("SELECT long_url FROM data WHERE short_url = ?").WithArgs(shortURL).
		WillReturnRows(sqlmock.NewRows([]string{"long_url"}))

	result, err := service.GetLongURL(db, shortURL)
	assert.Error(t, err)
	assert.Equal(t, "URL не найден", err.Error())
	assert.Empty(t, result)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unmet expectations: %v", err)
	}
}

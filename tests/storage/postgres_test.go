package storage

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"link_shortener/internal/storage"
)

func TestAddData(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	shortUrl := "abcd1234"
	longUrl := "https://example.com"

	mock.ExpectExec("INSERT INTO data").
		WithArgs(shortUrl, longUrl).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = storage.AddData(db, shortUrl, longUrl)
	assert.NoError(t, err)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetData_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	shortUrl := "abcd1234"
	longUrl := "https://example.com"

	rows := sqlmock.NewRows([]string{"long_url"}).
		AddRow(longUrl)

	mock.ExpectQuery("SELECT long_url FROM data WHERE short_url = ?").
		WithArgs(shortUrl).
		WillReturnRows(rows)

	result, err := storage.GetData(db, shortUrl)
	assert.NoError(t, err)
	assert.Equal(t, []string{longUrl}, result)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetData_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	shortUrl := "abcd1234"

	mock.ExpectQuery("SELECT long_url FROM data WHERE short_url = ?").
		WithArgs(shortUrl).
		WillReturnRows(sqlmock.NewRows([]string{"long_url"}))

	result, err := storage.GetData(db, shortUrl)
	assert.NoError(t, err)
	assert.Empty(t, result)

	assert.NoError(t, mock.ExpectationsWereMet())
}

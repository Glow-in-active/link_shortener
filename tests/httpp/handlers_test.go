package httpp

import (
	"bytes"
	"encoding/json"
	"github.com/alicebob/miniredis/v2"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"link_shortener/internal/httpp"
	"link_shortener/internal/storage"
	"net/http"
	"net/http/httptest"
	"testing"
)

var mockRedis *miniredis.Miniredis
var rdb *redis.Client

func setupMockRedis(t *testing.T) {
	var err error
	mockRedis, err = miniredis.Run()
	assert.NoError(t, err, "Ошибка при запуске мок-Redis")

	rdb = redis.NewClient(&redis.Options{
		Addr: mockRedis.Addr(),
	})
}

func TestShortenURL(t *testing.T) {
	setupMockRedis(t)

	gin.SetMode(gin.TestMode)
	router := gin.Default()

	database := &storage.Database{
		Type:  "redis",
		Redis: rdb,
	}

	handler := httpp.NewHandler(database)
	router.POST("/shorten", handler.ShortenURL)

	requestBody, _ := json.Marshal(map[string]string{
		"long_url": "https://example.com",
	})

	req, _ := http.NewRequest("POST", "/shorten", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code, "Ожидался код 200 OK")

	var response map[string]string
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err, "Ошибка парсинга JSON-ответа")
	assert.Contains(t, response, "short_url", "Ответ должен содержать short_url")
}

func TestResolveURL(t *testing.T) {
	setupMockRedis(t)

	gin.SetMode(gin.TestMode)
	router := gin.Default()

	database := &storage.Database{
		Type:  "redis",
		Redis: rdb,
	}

	handler := httpp.NewHandler(database)
	router.GET("/:shortURL", handler.ResolveURL)

	shortURL := "abc123"
	expectedLongURL := "https://example.com"
	err := rdb.Set(context.Background(), shortURL, expectedLongURL, 0).Err()
	assert.NoError(t, err, "Ошибка при добавлении тестовых данных в Redis")

	req, _ := http.NewRequest("GET", "/"+shortURL, nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code, "Ожидался код 200 OK")

	var response map[string]string
	err = json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err, "Ошибка парсинга JSON-ответа")
	assert.Equal(t, expectedLongURL, response["long_url"], "Должен вернуться корректный long_url")
}

func TestResolveNonExistingURL(t *testing.T) {
	setupMockRedis(t)

	gin.SetMode(gin.TestMode)
	router := gin.Default()

	database := &storage.Database{
		Type:  "redis",
		Redis: rdb,
	}

	handler := httpp.NewHandler(database)
	router.GET("/:shortURL", handler.ResolveURL)

	req, _ := http.NewRequest("GET", "/notfound", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNotFound, resp.Code, "Ожидался код 404 Not Found")
}

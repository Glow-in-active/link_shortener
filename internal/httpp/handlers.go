package httpp

import (
	"link_shortener/internal/service"
	"link_shortener/internal/storage"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	database *storage.Database
}

func NewHandler(database *storage.Database) *Handler {
	return &Handler{database: database}
}

func SetupRoutes(router *gin.Engine, database *storage.Database) {
	handler := NewHandler(database)

	router.POST("/shorten", handler.ShortenURL)
	router.GET("/:shortURL", handler.ResolveURL)
}

func (h *Handler) ShortenURL(c *gin.Context) {
	var request map[string]string
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "невалидный JSON"})
		return
	}

	longURL, ok := request["long_url"]
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "отсутствует параметр long_url"})
		return
	}

	var shortURL string
	var err error

	switch h.database.Type {
	case "postgres":
		shortURL, err = service.SaveURL(h.database.Postgres, longURL)
	case "redis":
		shortURL, err = service.SaveURLRedis(h.database.Redis, longURL)
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "неизвестный тип БД"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "не удалось создать ссылку"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"short_url": shortURL})
}

func (h *Handler) ResolveURL(c *gin.Context) {
	shortURL := c.Param("shortURL")

	var longURL string
	var err error

	switch h.database.Type {
	case "postgres":
		longURL, err = service.GetLongURL(h.database.Postgres, shortURL)
	case "redis":
		longURL, err = service.GetLongURLRedis(h.database.Redis, shortURL)
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "неизвестный тип БД"})
		return
	}

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ссылка не найдена"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"long_url": longURL})
}

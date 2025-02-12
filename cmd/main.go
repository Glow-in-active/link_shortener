package main

import (
	"github.com/gin-gonic/gin"
	"link_shortener/internal/httpp"
	"link_shortener/internal/storage"
	"log"
	"os"
)

func main() {
	configPath := "config.json"
	if _, err := os.Stat("/app/config.json"); err == nil {
		configPath = "/app/config.json"
	}
	config, err := storage.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	database, err := storage.InitDatabase(config.DBType, config)
	if err != nil {
		log.Fatalf("Ошибка инициализации базы данных: %v", err)
	}
	defer database.Close()

	router := gin.Default()
	httpp.SetupRoutes(router, database)

	log.Println("Сервер запущен на :8080, использует", config.DBType)
	router.Run(":8080")
}

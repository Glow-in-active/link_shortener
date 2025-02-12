package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"link_shortener/internal/httpp"
	"link_shortener/internal/storage"
	"log"
	"os"
)

type Config struct {
	Database string
}

func main() {
	configFile, err := os.Open("/app/config.json")
	if err != nil {
		log.Fatal("Ошибка открытия конфигурационного файла:", err)
	}
	defer configFile.Close()

	var config Config
	decoder := json.NewDecoder(configFile)
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatal("Ошибка при чтении конфигурационного файла:", err)
	}

	if config.Database == "" {
		config.Database = "redis"
	}

	database, err := storage.InitDatabase(config.Database)
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	router := gin.Default()
	httpp.SetupRoutes(router, database)

	log.Println("Сервер запущен на :8080, использует", config.Database)
	router.Run(":8080")
}

package main

import (
	"bufio"
	"fmt"
	"link_shortener/internal/httpp"
	"link_shortener/internal/storage"
	"log"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Выберите базу данных (postgres/redis, по умолчанию redis):")
	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)
	if choice == "" {
		choice = "redis"
	}

	database, err := storage.InitDatabase(choice)
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	router := gin.Default()
	httpp.SetupRoutes(router, database)

	log.Println("Сервер запущен на :8080, использует", choice)
	router.Run(":8080")
}

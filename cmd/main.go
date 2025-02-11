package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"link_shortener/internal/generator"
	"link_shortener/internal/storage"
	"log"
)

func main() {
	connStr := "postgres://admin:secret@localhost:5432/mydb?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Ошибка подключения к базе данных: ", err)
	}
	var url string
	fmt.Scan(&url)
	shortURL := generator.GenerateShortURL(url)

	storage.AddData(db, shortURL, url)
	data, err := storage.GetData(db)
	if err != nil {
		log.Fatal("Ошибка при извлечении данных: ", err)
	}
	fmt.Println("Данные из таблицы:", data)

}

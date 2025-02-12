package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"

	"link_shortener/internal/service"

	"log"
)

func main() {
	connStr := "postgres://admin:secret@localhost:5432/mydb?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Ошибка подключения к базе данных: ", err)
	}
	var url string
	var short string
	var long_url string
	url = "github.io"
	short, err = service.SaveURL(db, url)
	long_url, err = service.GetLongURL(db, short)
	fmt.Println(short)
	fmt.Println(long_url)
}

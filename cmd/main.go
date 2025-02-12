package main

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"link_shortener/internal/service"
)

func main() {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	var (
		shortUrl string
		longurl  string
		sns      string
	)

	fmt.Scan(&longurl)
	shortUrl, _ = service.SaveURLRedis(rdb, longurl)

	sns, _ = service.GetLongURLRedis(rdb, shortUrl)
	fmt.Println(shortUrl, sns)
}

package main

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"link_shortener/internal/storage"
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
	)

	fmt.Scan(&shortUrl)
	longurl, _ = storage.RedisGetData(rdb, shortUrl)
	fmt.Println(longurl)

}

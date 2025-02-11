package main

import (
	"fmt"
	"link_shortener/internal/generator"
)

func main() {
	var urll string
	urll = "github.com/fdsxfsdf/dsfsdfdsf"
	shortURL := generator.GenerateShortURL(urll)

	fmt.Println(shortURL)
}

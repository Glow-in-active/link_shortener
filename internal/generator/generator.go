package generator

import (
	"crypto/sha256"
	"encoding/hex"
)

func GenerateShortURL(url string) string {
	hash := sha256.Sum256([]byte(url))
	hashString := hex.EncodeToString(hash[:])
	shortURL := hashString[:10]
	return shortURL
}

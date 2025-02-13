package generator

import (
	"link_shortener/internal/generator"
	"testing"
)

func TestGenerateShortURL(t *testing.T) {
	testCases := []struct {
		inputURL       string
		expectedLength int
	}{
		{"https://example.com", 10},
		{"https://google.com", 10},
		{"https://someverylongurl.com/path/to/resource", 10},
	}

	for _, tc := range testCases {
		shortURL := generator.GenerateShortURL(tc.inputURL)

		if len(shortURL) != tc.expectedLength {
			t.Errorf("Для URL %s ожидалось 10 символов, но получено: %d (%s)", tc.inputURL, len(shortURL), shortURL)
		}
	}
}

func TestGenerateShortURLConsistency(t *testing.T) {
	url := "https://example.com"

	hash1 := generator.GenerateShortURL(url)
	hash2 := generator.GenerateShortURL(url)

	if hash1 != hash2 {
		t.Errorf("Одинаковый URL %s дал разные хэши: %s и %s", url, hash1, hash2)
	}
}

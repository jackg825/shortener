package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"time"
)

// GenerateShortURL takes a URL string and returns a shortened representation of it with characters limited to 'a-zA-Z0-9' and a length of 6 characters.
func GenerateShortURL(url string) string {
	hash := sha256.Sum256([]byte(url))
	encoded := hex.EncodeToString(hash[:])
	var shortURL string
	for _, char := range encoded {
		if len(shortURL) == 7 { // Extended length to 7 characters
			break
		}
		if ('a' <= char && char <= 'z') || ('A' <= char && char <= 'Z') || ('0' <= char && char <= '9') {
			shortURL += string(char)
		}
	}
	// Adding a character based on the current timestamp to avoid collision
	timestamp := time.Now().UnixNano()
	timestampLastChar := string((timestamp % 10) + '0') // Convert last digit of timestamp to string
	shortURL += timestampLastChar

	return shortURL
}

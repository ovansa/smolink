package utils

import (
	"crypto/rand"
	"encoding/base64"
	"strings"
	"time"

	mathrand "math/rand"
)

var charset = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func GenerateShortCode(length int) string {
	r := mathrand.New(mathrand.NewSource(time.Now().UnixNano()))

	b := make([]rune, length)
	for i := range b {
		b[i] = charset[r.Intn(len(charset))]
	}

	return string(b)
}

func GenerateShortCodeSecure(length int) (string, error) {
	// Calculate how many bytes we need
	byteLength := (length * 6) / 8 // 6 bits per character in base64
	if (length*6)%8 != 0 {
		byteLength++
	}

	// Generate random bytes
	randomBytes := make([]byte, byteLength)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	// Encode to base64 and truncate to desired length
	encoded := base64.URLEncoding.EncodeToString(randomBytes)
	encoded = strings.ReplaceAll(encoded, "+", "")
	encoded = strings.ReplaceAll(encoded, "/", "")
	encoded = strings.ReplaceAll(encoded, "=", "")

	if len(encoded) > length {
		encoded = encoded[:length]
	}

	return encoded, nil
}
